## 1. Core Layer — Browser Tool Panic Recovery

- [x] 1.1 Add `ErrBrowserPanic` sentinel error to `internal/tools/browser/browser.go`
- [x] 1.2 Add `safeRodCall(fn func() error) error` panic recovery wrapper
- [x] 1.3 Add `safeRodCallValue[T any](fn func() (T, error)) (T, error)` generic panic recovery wrapper
- [x] 1.4 Wrap all rod API calls in browser.go methods: Navigate, Screenshot, Click, Type, GetText, GetSnapshot, GetElementInfo, Eval, WaitForSelector, NewSession, Close

## 2. Application Layer — Session Manager Auto-Reconnect

- [x] 2.1 Update `SessionManager.EnsureSession()` to detect `ErrBrowserPanic` and retry once after closing

## 3. Application Layer — Browser Tool Handler Wrapper

- [x] 3.1 Add `wrapBrowserHandler` function to `internal/app/tools.go` with panic recovery and retry on `ErrBrowserPanic`
- [x] 3.2 Apply `wrapBrowserHandler` to all browser tools in `buildTools()`

## 4. Infrastructure Layer — WebSocket Panic Recovery

- [x] 4.1 Add panic recovery to `readPump()` defer block in `internal/gateway/server.go`
- [x] 4.2 Add panic recovery to `writePump()` defer block
- [x] 4.3 Extract RPC handler invocation into `handleRPC()` with isolated panic recovery

## 5. Docker — Chrome Sidecar Health Check

- [x] 5.1 Add healthcheck to Chrome service in `docker-compose.yml` (curl localhost:9222/json/version)
- [x] 5.2 Update lango-sidecar `depends_on` to use `condition: service_healthy`

## 6. Tests

- [x] 6.1 Add `TestSafeRodCall_RecoversPanic` — verify panic converts to ErrBrowserPanic error
- [x] 6.2 Add `TestSafeRodCall_PassesNormalError` — verify normal errors pass through unchanged
- [x] 6.3 Add `TestSafeRodCallValue_RecoversPanic` — verify generic version recovers panics
- [x] 6.4 Add `TestSafeRodCallValue_ReturnsValueOnSuccess` — verify normal returns work
- [x] 6.5 Add `TestErrBrowserPanic_Unwrap` — verify errors.Is works with wrapped errors

## 7. Verification

- [x] 7.1 Run `go build ./...` — verify build passes
- [x] 7.2 Run `go test ./internal/tools/browser/...` — verify panic recovery tests pass
- [x] 7.3 Run `go test ./...` — verify no regressions
