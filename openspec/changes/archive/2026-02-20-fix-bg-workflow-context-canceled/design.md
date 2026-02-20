## Context

Background tasks (`internal/background/`) and workflows (`internal/workflow/`) execute agent prompts in goroutines. Currently, both inherit the parent request's context, which is cancelled when the originating agent turn completes. This causes `context canceled` errors during LLM API calls.

The cron system (`internal/cron/scheduler.go:239`) already uses `context.Background()` and does not have this bug, establishing a precedent for the fix.

## Goals / Non-Goals

**Goals:**
- Background tasks and workflows survive parent context cancellation
- Context values (session key, approval target) are preserved in detached contexts
- Background tasks have a configurable maximum execution timeout
- Workflow tool handler returns immediately with a run ID (non-blocking)

**Non-Goals:**
- Changing cron system behavior (already works correctly)
- Adding retry/resume logic for failed tasks
- Persistent background task state across server restarts

## Decisions

### 1. Custom detached context type vs. plain `context.Background()`

**Decision**: Custom `detachedCtx` that delegates `Value()` to the parent.

**Rationale**: `context.Background()` loses all values. Background tasks need the session key for the agent runner and the approval target for tool approval routing. A custom type preserves values while decoupling the lifecycle.

**Alternative**: Manually extracting and re-injecting values — fragile and breaks when new value types are added.

### 2. Reusable `ctxutil` package vs. inline implementation

**Decision**: New `internal/ctxutil/` package with `Detach()` function.

**Rationale**: Both background and workflow need the same pattern. A shared utility avoids duplication and is available for future long-running subsystems.

### 3. Workflow `RunAsync()` vs. modifying `Run()` signature

**Decision**: Add `RunAsync()` that returns `(string, error)` alongside existing `Run()`.

**Rationale**: `Run()` is used by the resume path and potentially by tests. Adding `RunAsync()` preserves backward compatibility while enabling non-blocking tool handler behavior.

### 4. Task timeout as config field vs. hardcoded default

**Decision**: `BackgroundConfig.TaskTimeout` field with 30-minute default.

**Rationale**: Different deployments may need different timeout values. A config field follows the existing pattern (e.g., `WorkflowConfig.DefaultTimeout`, `AgentConfig.RequestTimeout`).

## Risks / Trade-offs

- **[Orphaned goroutines]** Detached contexts mean tasks can outlive the server shutdown signal. → Mitigated by existing `Manager.Shutdown()` and `Engine.Shutdown()` which cancel all active tasks.
- **[Duplicate step records]** `RunAsync()` pre-creates step records, then `runDAG()` also attempts creation. → Mitigated by treating `CreateStepRun` as idempotent (log + continue on error).
- **[Breaking API change in `NewManager`]** Adding `taskTimeout` parameter changes the signature. → This is an internal API; only `initBackground()` in `wiring.go` calls it.
