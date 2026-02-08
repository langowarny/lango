package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Server.Port != 18789 {
		t.Errorf("expected default port 18789, got %d", cfg.Server.Port)
	}

	if cfg.Agent.Provider != "anthropic" {
		t.Errorf("expected default provider anthropic, got %s", cfg.Agent.Provider)
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("expected default log level info, got %s", cfg.Logging.Level)
	}
}

func TestExpandEnvVars(t *testing.T) {
	os.Setenv("TEST_API_KEY", "sk-test-123")
	defer os.Unsetenv("TEST_API_KEY")

	result := expandEnvVars("${TEST_API_KEY}")
	if result != "sk-test-123" {
		t.Errorf("expected sk-test-123, got %s", result)
	}

	// Test non-existent variable (should keep original)
	result = expandEnvVars("${NON_EXISTENT_VAR}")
	if result != "${NON_EXISTENT_VAR}" {
		t.Errorf("expected ${NON_EXISTENT_VAR}, got %s", result)
	}
}

func TestValidate(t *testing.T) {
	// Valid config
	cfg := DefaultConfig()
	if err := Validate(cfg); err != nil {
		t.Errorf("expected valid config, got error: %v", err)
	}

	// Invalid port
	cfg.Server.Port = 0
	if err := Validate(cfg); err == nil {
		t.Error("expected error for invalid port")
	}
	cfg.Server.Port = 18789

	// Invalid provider
	cfg.Agent.Provider = "invalid"
	if err := Validate(cfg); err == nil {
		t.Error("expected error for invalid provider")
	}
	cfg.Agent.Provider = "anthropic"

	// Invalid log level
	cfg.Logging.Level = "invalid"
	if err := Validate(cfg); err == nil {
		t.Error("expected error for invalid log level")
	}
}
