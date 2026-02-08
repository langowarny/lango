//go:build integration

package agent

import (
	"context"
	"os"
	"testing"

	"github.com/langowarny/lango/internal/session"
	// "google.golang.org/adk/runner" // Not used directly here yet, accessed via Runtime
)

func TestIntegration_Gemini(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		t.Skip("skipping integration test: GOOGLE_API_KEY is not set")
	}

	// Setup Config
	cfg := Config{
		Provider:             "gemini",
		Model:                "gemini-2.0-flash-exp",
		APIKey:               apiKey,
		MaxConversationTurns: 20,
	}

	// Setup Mock Session Store (or simple in-memory implementation of our interface)
	store := &mockSessionStore{
		sessions: make(map[string]*session.Session),
	}

	// Create Runtime
	r, err := New(cfg, store)
	if err != nil {
		t.Fatalf("failed to create runtime: %v", err)
	}

	// Test Run
	ctx := context.Background()
	sessionKey := "integration-test-session"
	input := "Hello, are you working?"

	events := make(chan StreamEvent)
	go func() {
		if err := r.Run(ctx, sessionKey, input, events); err != nil {
			t.Errorf("Run failed: %v", err)
		}
	}()

	// Read events
	gotResponse := false
	for evt := range events {
		if evt.Type == "error" {
			t.Errorf("received error event: %v", evt.Error)
		}
		if evt.Type == "text_delta" {
			gotResponse = true
			t.Logf("Received text: %s", evt.Text)
		}
	}

	if !gotResponse {
		t.Error("did not receive any text response")
	}
}

// Simple mock for session store
type mockSessionStore struct {
	sessions map[string]*session.Session
}

func (m *mockSessionStore) Get(key string) (*session.Session, error) {
	if s, ok := m.sessions[key]; ok {
		return s, nil
	}
	return nil, os.ErrNotExist // Or custom error
}

func (m *mockSessionStore) Create(sess *session.Session) error {
	m.sessions[sess.Key] = sess
	return nil
}

func (m *mockSessionStore) AppendMessage(key string, msg session.Message) error {
	if s, ok := m.sessions[key]; ok {
		s.History = append(s.History, msg)
		return nil
	}
	return os.ErrNotExist
}

// Implement other methods of session.Store interface if required...
// Assuming minimal interface usage in Runtime.Run
func (m *mockSessionStore) List(userID string) ([]*session.Session, error) { return nil, nil }
func (m *mockSessionStore) Delete(key string) error                        { return nil, nil }
func (m *mockSessionStore) Update(sess *session.Session) error             { return nil, nil }
