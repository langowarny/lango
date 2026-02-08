package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/langowarny/lango/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityToolsRegistration(t *testing.T) {
	// Create temp dir for db
	tempDir, err := os.MkdirTemp("", "lango_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set passphrase env var for test
	os.Setenv("LANGO_PASSPHRASE", "test-passphrase")
	defer os.Unsetenv("LANGO_PASSPHRASE")

	dbPath := filepath.Join(tempDir, "sessions.db")

	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Agent: config.AgentConfig{
			Provider:  "gemini",
			Model:     "gemini-test",
			APIKey:    "test-key",
			MaxTokens: 100,
		},
		Session: config.SessionConfig{
			DatabasePath: dbPath,
		},
		Security: config.SecurityConfig{
			Signer: config.SignerConfig{
				Provider: "local",
			},
			Passphrase: "test-passphrase", // Should trigger initialization
		},
	}

	app, err := New(cfg)
	require.NoError(t, err)

	// Check if security tools are registered
	// The tools are registered in New(), so checking app.Agent is sufficient

	// Secrets tool
	tool, ok := app.Agent.GetTool("secrets")
	assert.True(t, ok, "secrets tool should be registered")
	if ok {
		assert.Equal(t, "secrets", tool.Name)
	}

	// Crypto tool
	tool, ok = app.Agent.GetTool("crypto")
	assert.True(t, ok, "crypto tool should be registered")
	if ok {
		assert.Equal(t, "crypto", tool.Name)
	}
}

func TestSecurityToolsDisabled(t *testing.T) {
	// Create temp dir for db
	tempDir, err := os.MkdirTemp("", "lango_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "sessions.db")

	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8081,
		},
		Agent: config.AgentConfig{
			Provider:  "gemini",
			Model:     "gemini-test",
			APIKey:    "test-key",
			MaxTokens: 100,
		},
		Session: config.SessionConfig{
			DatabasePath: dbPath,
		},
		Security: config.SecurityConfig{
			Signer: config.SignerConfig{
				// No provider set, or one that doesn't initialize
				Provider: "",
			},
		},
	}

	app, err := New(cfg)
	require.NoError(t, err)

	// Tools should NOT be registered
	_, ok := app.Agent.GetTool("secrets")
	assert.False(t, ok, "secrets tool should NOT be registered when security is disabled")

	_, ok = app.Agent.GetTool("crypto")
	assert.False(t, ok, "crypto tool should NOT be registered when security is disabled")
}

func TestApprovalMiddleware_Secrets(t *testing.T) {
	// Create temp dir for db
	tempDir, err := os.MkdirTemp("", "lango_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set passphrase env var for test
	os.Setenv("LANGO_PASSPHRASE", "test-passphrase")
	defer os.Unsetenv("LANGO_PASSPHRASE")

	dbPath := filepath.Join(tempDir, "sessions.db")

	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8082,
		},
		Agent: config.AgentConfig{
			Provider:  "gemini",
			Model:     "gemini-test",
			APIKey:    "test-key",
			MaxTokens: 100,
		},
		Session: config.SessionConfig{
			DatabasePath: dbPath,
		},
		Security: config.SecurityConfig{
			Signer: config.SignerConfig{
				Provider: "local",
			},
			Passphrase: "test-passphrase",
			Interceptor: config.InterceptorConfig{
				Enabled:          true,
				ApprovalRequired: true,
				// "secrets.get" is added by default in app.go
			},
		},
	}

	app, err := New(cfg)
	require.NoError(t, err)

	ctx := context.Background()

	// 1. secrets.list should NOT require approval
	listParams := map[string]interface{}{
		"operation": "list",
	}
	_, err = app.Agent.ExecuteTool(ctx, "secrets", listParams)
	require.NoError(t, err, "secrets.list should not require approval")
	// Verify result structure (it returns ListResult struct or map?)
	// secrets.go returns ListResult struct which might be marshalled?
	// The tool handler returns interface{}, likely parsing happens inside.
	// We just check no error for now.

	// 2. secrets.get SHOULD require approval
	getParams := map[string]interface{}{
		"operation": "get",
		"name":      "non-existent-secret", // Secret existence doesn't matter for approval check
	}
	_, err = app.Agent.ExecuteTool(ctx, "secrets", getParams)
	require.Error(t, err)
	// The error message depends on whether gateway is wired and companions are connected
	// Since in test gateway might not have companions, "no companion connected" is expected wrapped error
	assert.Contains(t, err.Error(), "no companion connected", "secrets.get should fail with connectivity error if no companion")
}
