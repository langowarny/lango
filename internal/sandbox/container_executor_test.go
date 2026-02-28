package sandbox

import (
	"context"
	"testing"
	"time"

	"github.com/langoai/lango/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRuntime is a test double for ContainerRuntime.
type mockRuntime struct {
	name      string
	available bool
	runResult *ExecutionResult
	runErr    error
}

func (m *mockRuntime) Run(_ context.Context, _ ContainerConfig) (*ExecutionResult, error) {
	return m.runResult, m.runErr
}

func (m *mockRuntime) Cleanup(_ context.Context, _ string) error {
	return nil
}

func (m *mockRuntime) IsAvailable(_ context.Context) bool {
	return m.available
}

func (m *mockRuntime) Name() string {
	return m.name
}

func TestContainerExecutor_FallbackToNative(t *testing.T) {
	// When runtime is "auto" and Docker/gVisor are unavailable, should fall back to native.
	cfg := Config{
		Enabled:        true,
		TimeoutPerTool: 30 * time.Second,
		MaxMemoryMB:    256,
	}
	containerCfg := config.ContainerSandboxConfig{
		Runtime:     "auto",
		Image:       "test-image:latest",
		NetworkMode: "none",
	}

	exec, err := NewContainerExecutor(cfg, containerCfg)
	require.NoError(t, err)
	// On CI/local without Docker, should fall back to native.
	assert.Contains(t, []string{"docker", "native"}, exec.RuntimeName())
}

func TestContainerExecutor_RuntimeName(t *testing.T) {
	exec := &ContainerExecutor{
		runtime: &mockRuntime{name: "test-runtime", available: true},
	}
	assert.Equal(t, "test-runtime", exec.RuntimeName())
}

func TestContainerExecutor_Execute_Success(t *testing.T) {
	mock := &mockRuntime{
		name:      "mock",
		available: true,
		runResult: &ExecutionResult{
			Output: map[string]interface{}{"status": "ok"},
		},
	}

	exec := &ContainerExecutor{
		runtime:     mock,
		cfg:         Config{TimeoutPerTool: 10 * time.Second, MaxMemoryMB: 128},
		image:       "test:latest",
		networkMode: "none",
		readOnly:    true,
	}

	result, err := exec.Execute(context.Background(), "echo", map[string]interface{}{"msg": "hello"})
	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}

func TestContainerExecutor_Execute_Error(t *testing.T) {
	mock := &mockRuntime{
		name:      "mock",
		available: true,
		runResult: nil,
		runErr:    ErrContainerTimeout,
	}

	exec := &ContainerExecutor{
		runtime:     mock,
		cfg:         Config{TimeoutPerTool: 10 * time.Second},
		image:       "test:latest",
		networkMode: "none",
	}

	_, err := exec.Execute(context.Background(), "slow-tool", nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrContainerTimeout)
}

func TestContainerExecutor_NativeRuntime_Explicit(t *testing.T) {
	cfg := Config{
		Enabled:        true,
		TimeoutPerTool: 5 * time.Second,
	}
	containerCfg := config.ContainerSandboxConfig{
		Runtime: "native",
		Image:   "unused",
	}

	// "native" is not docker/gvisor, so it falls through to native fallback.
	exec, err := NewContainerExecutor(cfg, containerCfg)
	require.NoError(t, err)
	assert.Equal(t, "native", exec.RuntimeName())
}

func TestContainerExecutor_DockerUnavailable_Explicit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Docker check in short mode")
	}

	cfg := Config{Enabled: true}
	containerCfg := config.ContainerSandboxConfig{
		Runtime: "docker",
		Image:   "test:latest",
	}

	exec, err := NewContainerExecutor(cfg, containerCfg)
	if err != nil {
		// Docker requested but unavailable — expected on some machines.
		assert.ErrorIs(t, err, ErrRuntimeUnavailable)
		return
	}
	// Docker is available on this machine.
	assert.Equal(t, "docker", exec.RuntimeName())
	_ = exec
}

func TestContainerExecutor_GVisorUnavailable_Explicit(t *testing.T) {
	cfg := Config{Enabled: true}
	containerCfg := config.ContainerSandboxConfig{
		Runtime: "gvisor",
		Image:   "test:latest",
	}

	_, err := NewContainerExecutor(cfg, containerCfg)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRuntimeUnavailable)
}

func TestContainerExecutor_Runtime(t *testing.T) {
	mock := &mockRuntime{name: "mock", available: true}
	exec := &ContainerExecutor{runtime: mock}
	assert.Equal(t, mock, exec.Runtime())
}
