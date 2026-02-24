package sandbox

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNativeRuntime_Name(t *testing.T) {
	rt := NewNativeRuntime(Config{TimeoutPerTool: 0})
	assert.Equal(t, "native", rt.Name())
}

func TestNativeRuntime_IsAvailable(t *testing.T) {
	rt := NewNativeRuntime(Config{})
	assert.True(t, rt.IsAvailable(context.Background()))
}

func TestNativeRuntime_Cleanup(t *testing.T) {
	rt := NewNativeRuntime(Config{})
	err := rt.Cleanup(context.Background(), "some-id")
	assert.NoError(t, err)
}

func TestGVisorRuntime_Name(t *testing.T) {
	rt := NewGVisorRuntime()
	assert.Equal(t, "gvisor", rt.Name())
}

func TestGVisorRuntime_IsAvailable(t *testing.T) {
	rt := NewGVisorRuntime()
	assert.False(t, rt.IsAvailable(context.Background()))
}

func TestGVisorRuntime_Run(t *testing.T) {
	rt := NewGVisorRuntime()
	_, err := rt.Run(context.Background(), ContainerConfig{})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRuntimeUnavailable)
}

func TestGVisorRuntime_Cleanup(t *testing.T) {
	rt := NewGVisorRuntime()
	err := rt.Cleanup(context.Background(), "some-id")
	assert.NoError(t, err)
}

func TestContainerConfig_Fields(t *testing.T) {
	cfg := ContainerConfig{
		Image:          "test-image:latest",
		ToolName:       "echo",
		NetworkMode:    "none",
		MemoryLimitMB:  256,
		CPUQuotaUS:     50000,
		ReadOnlyRootfs: true,
	}
	assert.Equal(t, "test-image:latest", cfg.Image)
	assert.Equal(t, "echo", cfg.ToolName)
	assert.Equal(t, "none", cfg.NetworkMode)
	assert.Equal(t, int64(256), cfg.MemoryLimitMB)
	assert.Equal(t, int64(50000), cfg.CPUQuotaUS)
	assert.True(t, cfg.ReadOnlyRootfs)
}

func TestErrorSentinels(t *testing.T) {
	assert.Error(t, ErrRuntimeUnavailable)
	assert.Error(t, ErrContainerTimeout)
	assert.Error(t, ErrContainerOOM)
	assert.Equal(t, "container runtime unavailable", ErrRuntimeUnavailable.Error())
	assert.Equal(t, "container execution timed out", ErrContainerTimeout.Error())
	assert.Equal(t, "container killed due to out-of-memory", ErrContainerOOM.Error())
}
