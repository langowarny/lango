package sandbox

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSubprocessExecutor(t *testing.T) {
	cfg := Config{TimeoutPerTool: 10 * time.Second}
	exec := NewSubprocessExecutor(cfg)
	require.NotNil(t, exec)
	assert.Equal(t, 10*time.Second, exec.cfg.TimeoutPerTool)
}

func TestCleanEnv(t *testing.T) {
	env := cleanEnv()
	for _, e := range env {
		assert.True(t,
			len(e) > 5 && (e[:5] == "PATH=" || e[:5] == "HOME="),
			"unexpected env var: %s", e,
		)
	}
}

func TestWorkerFlag_Constant(t *testing.T) {
	assert.Equal(t, "--sandbox-worker", workerFlag)
}

func TestIsWorkerMode_Default(t *testing.T) {
	// In normal test execution, --sandbox-worker is not passed.
	assert.False(t, IsWorkerMode())
}
