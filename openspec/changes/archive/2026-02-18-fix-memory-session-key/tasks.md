## 1. Core Fix — ContextAwareModelAdapter

- [x] 1.1 Remove `sessionKey` field from `ContextAwareModelAdapter` struct in `internal/adk/context_model.go`
- [x] 1.2 Change `WithMemory` signature from `WithMemory(provider, sessionKey)` to `WithMemory(provider)`, remove sessionKey assignment
- [x] 1.3 Add `session` package import to `context_model.go`
- [x] 1.4 In `GenerateContent`, resolve session key from context via `session.SessionKeyFromContext(ctx)` at method entry
- [x] 1.5 Replace all `m.sessionKey` references in `GenerateContent` with local `sessionKey` variable

## 2. Method Signature Updates

- [x] 2.1 Add `sessionKey string` parameter to `assembleMemorySection` and replace internal `m.sessionKey` usage
- [x] 2.2 Add `sessionKey string` parameter to `assembleRAGSection` and replace internal `m.sessionKey` usage
- [x] 2.3 Add `sessionKey string` parameter to `assembleGraphRAGSection` and replace internal `m.sessionKey` usage
- [x] 2.4 Update all call sites in `GenerateContent` to pass `sessionKey` to the assemble methods

## 3. Wiring Updates

- [x] 3.1 Update `wiring.go` line 675: change `ctxAdapter.WithMemory(mc.store, "")` to `ctxAdapter.WithMemory(mc.store)`
- [x] 3.2 Update `wiring.go` line 711: change `ctxAdapter.WithMemory(mc.store, "")` to `ctxAdapter.WithMemory(mc.store)`

## 4. Tests

- [x] 4.1 Create `internal/adk/context_model_test.go` with mock `MemoryProvider`
- [x] 4.2 Add test: session key from context is passed to memory provider
- [x] 4.3 Add test: no session key in context skips memory retrieval
- [x] 4.4 Add test: session key updates RuntimeContextAdapter
- [x] 4.5 Add test: memory content is injected into system prompt

## 5. Verification

- [x] 5.1 Run `go build ./...` — confirm no compilation errors
- [x] 5.2 Run `go test ./internal/adk/...` — confirm all tests pass
- [x] 5.3 Run `go test ./...` — confirm no regressions
