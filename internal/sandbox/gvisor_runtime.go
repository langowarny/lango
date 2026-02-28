package sandbox

import (
	"context"
)

// GVisorRuntime is a stub for future gVisor-based container isolation.
// It always reports as unavailable and returns ErrRuntimeUnavailable on Run.
type GVisorRuntime struct{}

// NewGVisorRuntime creates a GVisorRuntime stub.
func NewGVisorRuntime() *GVisorRuntime {
	return &GVisorRuntime{}
}

// Run returns ErrRuntimeUnavailable — gVisor support is not yet implemented.
func (r *GVisorRuntime) Run(_ context.Context, _ ContainerConfig) (*ExecutionResult, error) {
	return nil, ErrRuntimeUnavailable
}

// Cleanup is a no-op for the gVisor stub.
func (r *GVisorRuntime) Cleanup(_ context.Context, _ string) error {
	return nil
}

// IsAvailable always returns false for the gVisor stub.
func (r *GVisorRuntime) IsAvailable(_ context.Context) bool {
	return false
}

// Name returns the runtime name.
func (r *GVisorRuntime) Name() string {
	return "gvisor"
}
