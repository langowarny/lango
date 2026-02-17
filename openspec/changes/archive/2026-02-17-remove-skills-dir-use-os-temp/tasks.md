## 1. Executor Simplification

- [x] 1.1 Remove `skillsDir` field from `Executor` struct in `internal/skill/executor.go`
- [x] 1.2 Simplify `NewExecutor` to return `*Executor` (no error): remove home dir lookup and `MkdirAll`
- [x] 1.3 Rewrite `executeScript` to use `os.CreateTemp("", "lango-skill-*.sh")` with proper write/close/defer cleanup
- [x] 1.4 Remove `"path/filepath"` from imports

## 2. Registry Simplification

- [x] 2.1 Change `NewRegistry` return type from `(*Registry, error)` to `*Registry`
- [x] 2.2 Remove error wrapping around `NewExecutor` call in `NewRegistry`

## 3. Caller Updates

- [x] 3.1 Update `initKnowledge` in `internal/app/wiring.go` to call `NewRegistry` without error handling

## 4. Test Updates

- [x] 4.1 Remove `t.Setenv("HOME", ...)` and error check from `newTestExecutor` in `executor_test.go`
- [x] 4.2 Remove `t.Setenv("HOME", ...)` and error check from `newTestRegistry` in `registry_test.go`

## 5. Verification

- [x] 5.1 Run `go build ./...` — confirm clean build
- [x] 5.2 Run `go test ./internal/skill/...` — confirm all tests pass
