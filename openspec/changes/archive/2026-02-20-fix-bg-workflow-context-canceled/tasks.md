## 1. Context Detach Utility

- [x] 1.1 Create `internal/ctxutil/detach.go` with `Detach()` function and `detachedCtx` type
- [x] 1.2 Create `internal/ctxutil/detach_test.go` with tests for cancellation isolation, value preservation, deadline absence, and child wrapping

## 2. Background Task Context Fix

- [x] 2.1 Add `TaskTimeout` field to `BackgroundConfig` in `internal/config/types.go`
- [x] 2.2 Add `taskTimeout` field to `Manager` struct and update `NewManager()` signature in `internal/background/manager.go`
- [x] 2.3 Replace `context.WithCancel(ctx)` with `ctxutil.Detach(ctx)` + `context.WithTimeout()` in `Manager.Submit()`
- [x] 2.4 Update `initBackground()` in `internal/app/wiring.go` to pass `taskTimeout` to `NewManager()`

## 3. Workflow Engine Context Fix

- [x] 3.1 Detach context in `Engine.Run()` before creating run record
- [x] 3.2 Extract DAG execution logic into `runDAG()` method
- [x] 3.3 Add `RunAsync()` method that creates records and launches `runDAG()` in a goroutine
- [x] 3.4 Update `workflow_run` tool handler in `internal/app/tools.go` to use `RunAsync()`

## 4. Verification

- [x] 4.1 `go build ./...` passes
- [x] 4.2 `go test ./internal/ctxutil/...` passes
- [x] 4.3 `go test ./internal/background/...` passes (no test files)
- [x] 4.4 `go test ./internal/workflow/...` passes (no test files)
- [x] 4.5 `go test ./internal/app/...` passes
