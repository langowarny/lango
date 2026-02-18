## 1. Core Fix

- [x] 1.1 Move handler call in `handleMessage` to a separate goroutine with `c.wg.Add(1)` / `defer c.wg.Done()`

## 2. Verification

- [x] 2.1 Run `go build ./...` and confirm no compilation errors
- [x] 2.2 Run `go test ./internal/channels/slack/... -v` and confirm all tests pass
- [x] 2.3 Run `go test ./internal/channels/telegram/... -v` and confirm no regressions
