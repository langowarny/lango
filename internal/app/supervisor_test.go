package app

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/langowarny/lango/internal/config"
	"github.com/stretchr/testify/require"
)

func TestSupervisorExecWhitelist(t *testing.T) {
	// 1. Setup Environment
	secretVar := "SUPERVISOR_TEST_SECRET"
	secretVal := "should_not_be_seen"
	os.Setenv(secretVar, secretVal)
	defer os.Unsetenv(secretVar)

	// 2. Initialize App
	cfg := &config.Config{
		Tools: config.ToolsConfig{
			Exec: config.ExecToolConfig{
				DefaultTimeout: 5 * time.Second,
				WorkDir:        ".",
			},
		},
		Agent: config.AgentConfig{
			Provider: "noop", // Only testing tools
			Model:    "noop",
		},
		Session: config.SessionConfig{
			DatabasePath: ":memory:",
		},
	}

	app, err := New(cfg)
	require.NoError(t, err)

	// 3. Execute "env" command via Agent Tool
	ctx := context.Background()
	params := map[string]interface{}{
		"command": "env",
	}

	result, err := app.Agent.ExecuteTool(ctx, "exec", params)
	require.NoError(t, err)

	output, ok := result.(string)
	require.True(t, ok, "exec tool should return string output")

	// 4. Verification
	// The secret variable should NOT be present
	if strings.Contains(output, secretVar) || strings.Contains(output, secretVal) {
		t.Errorf("Security Leak! Found leaked environment variable in output:\n%s", output)
	}

	// Whitelisted variables SHOULD be present (e.g., PATH)
	if !strings.Contains(output, "PATH=") {
		t.Errorf("Expected whitelisted variable PATH to be present, but it was missing.\nOutput:\n%s", output)
	}
}
