package checks

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/langowarny/lango/internal/config"
)

// ConfigCheck validates the configuration file.
type ConfigCheck struct{}

// Name returns the check name.
func (c *ConfigCheck) Name() string {
	return "Configuration File"
}

// Run checks if the configuration file exists and is valid.
func (c *ConfigCheck) Run(ctx context.Context, cfg *config.Config) Result {
	configPath := findConfigPath()

	if configPath == "" {
		return Result{
			Name:      c.Name(),
			Status:    StatusFail,
			Message:   "Configuration file not found",
			Details:   "No lango.json found in current directory or ~/.lango/",
			Fixable:   true,
			FixAction: "Create default configuration file",
		}
	}

	// Check if file is readable
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: "Cannot read configuration file",
			Details: err.Error(),
		}
	}

	// Validate JSON syntax
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: "Invalid JSON syntax",
			Details: err.Error(),
		}
	}

	// Try to load with full validation
	if _, err := config.Load(configPath); err != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusWarn,
			Message: "Configuration has validation warnings",
			Details: err.Error(),
		}
	}

	return Result{
		Name:    c.Name(),
		Status:  StatusPass,
		Message: "Configuration file valid",
		Details: configPath,
	}
}

// Fix creates a default configuration file if missing.
func (c *ConfigCheck) Fix(ctx context.Context, cfg *config.Config) Result {
	configPath := findConfigPath()
	if configPath != "" {
		return Result{
			Name:    c.Name(),
			Status:  StatusPass,
			Message: "Configuration file already exists",
			Details: configPath,
		}
	}

	// Create default config in current directory
	defaultConfig := config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 18789,
		},
		Agent: config.AgentConfig{
			Provider: "gemini",
			Model:    "gemini-2.0-flash-exp",
		},
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "console",
		},
	}

	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: "Failed to create default config",
			Details: err.Error(),
		}
	}

	if err := os.WriteFile("lango.json", data, 0644); err != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: "Failed to write config file",
			Details: err.Error(),
		}
	}

	return Result{
		Name:    c.Name(),
		Status:  StatusPass,
		Message: "Created default configuration file",
		Details: "lango.json",
	}
}

// findConfigPath looks for lango.json in common locations.
func findConfigPath() string {
	// Check current directory
	if _, err := os.Stat("lango.json"); err == nil {
		return "lango.json"
	}

	// Check home directory
	home, err := os.UserHomeDir()
	if err == nil {
		homeConfig := filepath.Join(home, ".lango", "lango.json")
		if _, err := os.Stat(homeConfig); err == nil {
			return homeConfig
		}
	}

	return ""
}
