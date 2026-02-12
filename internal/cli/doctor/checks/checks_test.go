package checks

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/langowarny/lango/internal/config"
)

func TestConfigCheck_Run_ValidConfig(t *testing.T) {
	// Create a temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "lango.json")
	configContent := `{
		"server": {"host": "localhost", "port": 18789},
		"agent": {"provider": "gemini", "model": "gemini-2.0-flash-exp"}
	}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to temp directory
	oldCwd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldCwd)

	check := &ConfigCheck{}
	result := check.Run(context.Background(), nil)

	// Accept both pass and warn (validation warnings are acceptable)
	if result.Status != StatusPass && result.Status != StatusWarn {
		t.Errorf("expected StatusPass or StatusWarn, got %v: %s", result.Status, result.Message)
	}
}

func TestConfigCheck_Run_MissingConfig(t *testing.T) {
	// Change to empty temp directory
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldCwd)

	check := &ConfigCheck{}
	result := check.Run(context.Background(), nil)

	if result.Status != StatusFail {
		t.Errorf("expected StatusFail, got %v", result.Status)
	}
	if !result.Fixable {
		t.Error("expected Fixable to be true")
	}
}

func TestProvidersCheck_Run_NoConfig(t *testing.T) {
	check := &ProvidersCheck{}
	result := check.Run(context.Background(), nil)

	if result.Status != StatusSkip {
		t.Errorf("expected StatusSkip, got %v", result.Status)
	}
}

func TestProvidersCheck_Run_ExplicitProvider(t *testing.T) {
	cfg := &config.Config{
		Providers: map[string]config.ProviderConfig{
			"openai": {Type: "openai", APIKey: "sk-test"},
		},
	}

	check := &ProvidersCheck{}
	result := check.Run(context.Background(), cfg)

	if result.Status != StatusPass {
		t.Errorf("expected StatusPass, got %v: %s", result.Status, result.Message)
	}
}

func TestProvidersCheck_Run_ProviderInMap(t *testing.T) {
	cfg := &config.Config{
		Agent: config.AgentConfig{
			Provider: "gemini",
		},
		Providers: map[string]config.ProviderConfig{
			"gemini": {Type: "gemini", APIKey: "test-key"},
		},
	}

	check := &ProvidersCheck{}
	result := check.Run(context.Background(), cfg)

	if result.Status != StatusPass {
		t.Errorf("expected StatusPass, got %v: %s", result.Status, result.Message)
	}
}

func TestProvidersCheck_Run_MissingKey(t *testing.T) {
	cfg := &config.Config{
		Providers: map[string]config.ProviderConfig{
			"anthropic": {Type: "anthropic", APIKey: ""},
		},
	}

	check := &ProvidersCheck{}
	result := check.Run(context.Background(), cfg)

	if result.Status != StatusFail {
		t.Errorf("expected StatusFail, got %v", result.Status)
	}
}

func TestChannelCheck_Run_NoChannelsEnabled(t *testing.T) {
	cfg := &config.Config{}

	check := &ChannelCheck{}
	result := check.Run(context.Background(), cfg)

	if result.Status != StatusWarn {
		t.Errorf("expected StatusWarn for no channels, got %v", result.Status)
	}
}

func TestChannelCheck_Run_TelegramEnabledNoToken(t *testing.T) {
	cfg := &config.Config{}
	cfg.Channels.Telegram.Enabled = true

	check := &ChannelCheck{}
	result := check.Run(context.Background(), cfg)

	if result.Status != StatusFail {
		t.Errorf("expected StatusFail for missing token, got %v", result.Status)
	}
}

func TestNetworkCheck_Run_PortAvailable(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 19999, // Use a high port unlikely to be in use
		},
	}

	check := &NetworkCheck{}
	result := check.Run(context.Background(), cfg)

	if result.Status != StatusPass {
		t.Errorf("expected StatusPass, got %v: %s", result.Status, result.Message)
	}
}

func TestDatabaseCheck_Run_DirectoryNotExist(t *testing.T) {
	cfg := &config.Config{
		Session: config.SessionConfig{
			DatabasePath: "/nonexistent/path/sessions.db",
		},
	}

	check := &DatabaseCheck{}
	result := check.Run(context.Background(), cfg)

	if result.Status != StatusFail {
		t.Errorf("expected StatusFail, got %v", result.Status)
	}
	if !result.Fixable {
		t.Error("expected Fixable to be true")
	}
}

func TestNewSummary(t *testing.T) {
	results := []Result{
		{Name: "Test1", Status: StatusPass},
		{Name: "Test2", Status: StatusPass},
		{Name: "Test3", Status: StatusWarn},
		{Name: "Test4", Status: StatusFail},
		{Name: "Test5", Status: StatusSkip},
	}

	summary := NewSummary(results)

	if summary.Passed != 2 {
		t.Errorf("expected 2 passed, got %d", summary.Passed)
	}
	if summary.Warnings != 1 {
		t.Errorf("expected 1 warning, got %d", summary.Warnings)
	}
	if summary.Failed != 1 {
		t.Errorf("expected 1 failed, got %d", summary.Failed)
	}
	if summary.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", summary.Skipped)
	}
	if !summary.HasErrors() {
		t.Error("expected HasErrors to be true")
	}
	if !summary.HasWarnings() {
		t.Error("expected HasWarnings to be true")
	}
}
