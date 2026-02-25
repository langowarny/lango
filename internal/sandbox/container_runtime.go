package sandbox

import (
	"context"
	"errors"
	"time"
)

// ContainerConfig defines parameters for a containerized tool execution.
type ContainerConfig struct {
	// Image is the Docker image to use for the container.
	Image string

	// ToolName is the name of the tool to execute.
	ToolName string

	// NetworkMode is the Docker network mode (e.g. "none", "bridge").
	NetworkMode string

	// Params are the tool invocation parameters.
	Params map[string]interface{}

	// MemoryLimitMB is the hard memory limit in megabytes.
	MemoryLimitMB int64

	// CPUQuotaUS is the CPU quota in microseconds.
	CPUQuotaUS int64

	// ReadOnlyRootfs mounts the root filesystem as read-only.
	ReadOnlyRootfs bool

	// Timeout is the maximum execution duration.
	Timeout time.Duration
}

// ContainerRuntime provides an execution environment for isolated tool runs.
type ContainerRuntime interface {
	// Run executes a tool inside a container and returns the result.
	Run(ctx context.Context, cfg ContainerConfig) (*ExecutionResult, error)

	// Cleanup removes containers associated with the given container ID.
	Cleanup(ctx context.Context, containerID string) error

	// IsAvailable checks whether the runtime is operational.
	IsAvailable(ctx context.Context) bool

	// Name returns the human-readable runtime name.
	Name() string
}

var (
	// ErrRuntimeUnavailable indicates the container runtime is not installed or accessible.
	ErrRuntimeUnavailable = errors.New("container runtime unavailable")

	// ErrContainerTimeout indicates the container execution exceeded its deadline.
	ErrContainerTimeout = errors.New("container execution timed out")

	// ErrContainerOOM indicates the container was killed due to out-of-memory.
	ErrContainerOOM = errors.New("container killed due to out-of-memory")
)
