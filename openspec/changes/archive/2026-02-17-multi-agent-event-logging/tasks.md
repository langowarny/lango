## 1. Logger Setup

- [x] 1.1 Add `logging` import and `logger()` helper function to `internal/adk/agent.go`

## 2. Event Stream Logging

- [x] 2.1 Add `hasText()` helper to check if an event contains non-empty text parts
- [x] 2.2 Add delegation logging in `RunAndCollect()` — log when `Author` and `TransferToAgent` are set
- [x] 2.3 Add response logging in `RunAndCollect()` — log when `Author` is set and event has text content

## 3. Verification

- [x] 3.1 Run `go build ./...` to verify compilation
- [x] 3.2 Run `go test ./internal/adk/...` to verify tests pass
