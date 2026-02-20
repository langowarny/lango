## Why

Background tasks and workflow executions fail with `context canceled` errors because they inherit the parent request's context. When the originating agent request completes and its context is cancelled, all background goroutines using derived contexts are also cancelled, causing Gemini/LLM API calls to fail mid-execution.

## What Changes

- Introduce a `ctxutil.Detach()` utility that creates a context independent of the parent's cancellation while preserving context values (session keys, approval targets).
- Modify `background.Manager.Submit()` to detach from the parent context and apply a configurable task timeout instead.
- Modify `workflow.Engine.Run()` to detach from the parent context before DAG execution.
- Add `workflow.Engine.RunAsync()` for non-blocking workflow execution from tool handlers.
- Add `BackgroundConfig.TaskTimeout` field for configurable per-task timeout (default: 30m).
- Change `workflow_run` tool handler to use `RunAsync()` for immediate response with `run_id`.

## Capabilities

### New Capabilities
- `context-detach`: Utility to detach a context from its parent's cancellation while preserving values, enabling long-running goroutines to survive their originating request.

### Modified Capabilities
- `background-execution`: Background task context is now detached from the parent request with a configurable task timeout.
- `workflow-engine`: Workflow execution context is now detached; added RunAsync for non-blocking execution from tool handlers.

## Impact

- `internal/ctxutil/` — new package
- `internal/config/types.go` — new `TaskTimeout` field on `BackgroundConfig`
- `internal/background/manager.go` — detached context, new `taskTimeout` parameter
- `internal/workflow/engine.go` — detached context, new `RunAsync()` method
- `internal/app/tools.go` — `workflow_run` handler uses `RunAsync()`
- `internal/app/wiring.go` — passes `taskTimeout` to `NewManager()`
