package sandbox

import (
	"context"
)

// NativeRuntime wraps the SubprocessExecutor as a ContainerRuntime fallback.
// It provides process-level isolation without containerization.
type NativeRuntime struct {
	executor *SubprocessExecutor
}

// NewNativeRuntime creates a NativeRuntime backed by a SubprocessExecutor.
func NewNativeRuntime(cfg Config) *NativeRuntime {
	return &NativeRuntime{
		executor: NewSubprocessExecutor(cfg),
	}
}

// Run executes the tool via the subprocess executor, adapting ContainerConfig
// parameters to the subprocess model.
func (r *NativeRuntime) Run(ctx context.Context, cfg ContainerConfig) (*ExecutionResult, error) {
	output, err := r.executor.Execute(ctx, cfg.ToolName, cfg.Params)
	if err != nil {
		return &ExecutionResult{Error: err.Error()}, err
	}
	return &ExecutionResult{Output: output}, nil
}

// Cleanup is a no-op for native runtime — subprocesses are cleaned up on exit.
func (r *NativeRuntime) Cleanup(_ context.Context, _ string) error {
	return nil
}

// IsAvailable always returns true for native runtime.
func (r *NativeRuntime) IsAvailable(_ context.Context) bool {
	return true
}

// Name returns the runtime name.
func (r *NativeRuntime) Name() string {
	return "native"
}
