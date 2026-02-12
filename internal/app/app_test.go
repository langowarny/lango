package app

import (
	"testing"

	"github.com/langowarny/lango/internal/config"
)

func TestNew_MinimalConfig(t *testing.T) {
	t.Skip("requires provider credentials; run manually with GOOGLE_API_KEY set")

	cfg := config.DefaultConfig()
	cfg.Agent.Provider = "google"
	cfg.Agent.Model = "gemini-2.0-flash"
	cfg.Providers = map[string]config.ProviderConfig{
		"google": {
			Type:   "gemini",
			APIKey: "test-key",
		},
	}

	app, err := New(cfg)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	if app.Agent == nil {
		t.Fatal("expected agent to be initialized")
	}
	if app.Gateway == nil {
		t.Fatal("expected gateway to be initialized")
	}
	if app.Store == nil {
		t.Fatal("expected store to be initialized")
	}
}

func TestNew_SecurityDisabledByDefault(t *testing.T) {
	t.Skip("requires provider credentials; run manually with GOOGLE_API_KEY set")

	cfg := config.DefaultConfig()
	cfg.Agent.Provider = "google"
	cfg.Providers = map[string]config.ProviderConfig{
		"google": {
			Type:   "gemini",
			APIKey: "test-key",
		},
	}

	// Security is not configured — should not block startup
	_, err := New(cfg)
	if err != nil {
		t.Fatalf("New() should succeed without security config, got: %v", err)
	}
}
