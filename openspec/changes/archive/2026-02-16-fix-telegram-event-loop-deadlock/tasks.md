## 1. Fix Event Loop Deadlock

- [x] 1.1 Wrap `c.handleUpdate(ctx, update)` call in `Start()` with `c.wg.Add(1)` and `go func()` with `defer c.wg.Done()` (`internal/channels/telegram/telegram.go`)

## 2. Verification

- [x] 2.1 Run existing telegram tests (`go test ./internal/channels/telegram/... -v`)
- [x] 2.2 Verify `go build ./...` succeeds with no errors
