package agent

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPresidioDetector_SuccessfulDetection(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/analyze" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		results := []presidioResult{
			{EntityType: "EMAIL_ADDRESS", Start: 10, End: 25, Score: 0.95},
			{EntityType: "PHONE_NUMBER", Start: 30, End: 42, Score: 0.8},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}))
	defer srv.Close()

	det := NewPresidioDetector(srv.URL, WithPresidioThreshold(0.5))
	matches := det.Detect("My email user@example.com, phone 123-456-7890")

	if len(matches) != 2 {
		t.Fatalf("want 2 matches, got %d", len(matches))
	}

	if matches[0].PatternName != "presidio:EMAIL_ADDRESS" {
		t.Errorf("match[0] name: want %q, got %q", "presidio:EMAIL_ADDRESS", matches[0].PatternName)
	}
	if matches[0].Category != PIICategoryContact {
		t.Errorf("match[0] category: want %q, got %q", PIICategoryContact, matches[0].Category)
	}
	if matches[0].Score != 0.95 {
		t.Errorf("match[0] score: want 0.95, got %f", matches[0].Score)
	}

	if matches[1].PatternName != "presidio:PHONE_NUMBER" {
		t.Errorf("match[1] name: want %q, got %q", "presidio:PHONE_NUMBER", matches[1].PatternName)
	}
}

func TestPresidioDetector_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	det := NewPresidioDetector(srv.URL)
	matches := det.Detect("test@example.com")

	// Graceful degradation: returns nil on error
	if matches != nil {
		t.Errorf("want nil on server error, got %v", matches)
	}
}

func TestPresidioDetector_ConnectionError(t *testing.T) {
	det := NewPresidioDetector("http://localhost:1", WithPresidioTimeout(100*time.Millisecond))
	matches := det.Detect("test@example.com")

	if matches != nil {
		t.Errorf("want nil on connection error, got %v", matches)
	}
}

func TestPresidioDetector_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	det := NewPresidioDetector(srv.URL)
	matches := det.Detect("test@example.com")

	if matches != nil {
		t.Errorf("want nil on invalid JSON, got %v", matches)
	}
}

func TestPresidioDetector_EmptyResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]presidioResult{})
	}))
	defer srv.Close()

	det := NewPresidioDetector(srv.URL)
	matches := det.Detect("no PII here")

	if len(matches) != 0 {
		t.Errorf("want 0 matches, got %d", len(matches))
	}
}

func TestPresidioDetector_HealthCheck(t *testing.T) {
	t.Run("healthy", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		det := NewPresidioDetector(srv.URL)
		err := det.HealthCheck(context.Background())
		if err != nil {
			t.Errorf("want nil error, got %v", err)
		}
	})

	t.Run("unhealthy", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer srv.Close()

		det := NewPresidioDetector(srv.URL)
		err := det.HealthCheck(context.Background())
		if err == nil {
			t.Error("want error for unhealthy service, got nil")
		}
	})

	t.Run("unreachable", func(t *testing.T) {
		det := NewPresidioDetector("http://localhost:1", WithPresidioTimeout(100*time.Millisecond))
		err := det.HealthCheck(context.Background())
		if err == nil {
			t.Error("want error for unreachable service, got nil")
		}
	})
}

func TestPresidioDetector_Options(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req presidioRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Language != "ko" {
			t.Errorf("want language %q, got %q", "ko", req.Language)
		}
		if req.ScoreThreshold != 0.9 {
			t.Errorf("want threshold 0.9, got %f", req.ScoreThreshold)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]presidioResult{})
	}))
	defer srv.Close()

	det := NewPresidioDetector(srv.URL,
		WithPresidioLanguage("ko"),
		WithPresidioThreshold(0.9),
	)
	det.Detect("test input")
}

func TestPresidioDetector_UnknownEntityType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		results := []presidioResult{
			{EntityType: "UNKNOWN_TYPE", Start: 0, End: 5, Score: 0.7},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}))
	defer srv.Close()

	det := NewPresidioDetector(srv.URL)
	matches := det.Detect("hello")

	if len(matches) != 1 {
		t.Fatalf("want 1 match, got %d", len(matches))
	}
	// Unknown entity types default to identity category
	if matches[0].Category != PIICategoryIdentity {
		t.Errorf("category: want %q, got %q", PIICategoryIdentity, matches[0].Category)
	}
}
