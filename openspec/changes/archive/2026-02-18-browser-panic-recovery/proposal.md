## Why

Docker slim + `chromedp/headless-shell` sidecar environment에서 브라우저 도구 사용 시, `go-rod/rod` 라이브러리가 Chrome과의 CDP/WebSocket 연결 끊김에 대해 error 대신 panic을 발생시켜 프로세스 전체가 크래시된다. Tool 실행 경로와 WebSocket goroutine에 panic recovery가 없어 단일 브라우저 장애가 전체 서비스 중단으로 이어진다.

## What Changes

- Add `ErrBrowserPanic` sentinel error and `safeRodCall`/`safeRodCallValue` panic recovery wrappers in the browser tool core layer
- Wrap all rod/CDP method calls (Navigate, Screenshot, Click, Type, GetText, GetSnapshot, GetElementInfo, Eval, WaitForSelector, NewSession, Close) with panic recovery
- Add auto-reconnect logic in `SessionManager.EnsureSession()` on `ErrBrowserPanic` detection (close + retry once)
- Add `wrapBrowserHandler` in the application layer to catch panics and retry on `ErrBrowserPanic` at the tool handler level
- Add panic recovery to WebSocket `readPump`/`writePump` goroutines and isolate RPC handler panics so a single handler crash does not tear down the connection
- Add Chrome sidecar healthcheck and `service_healthy` dependency condition in docker-compose.yml

## Capabilities

### New Capabilities

### Modified Capabilities
- `tool-browser`: Add panic recovery layer around all rod/CDP calls and auto-reconnect on connection loss
- `docker-deployment`: Add Chrome sidecar healthcheck and healthy dependency condition
- `gateway-server`: Add panic recovery to WebSocket read/write pumps and RPC handler invocations

## Impact

- `internal/tools/browser/browser.go` — new error type, panic recovery wrappers, all methods wrapped
- `internal/tools/browser/session_manager.go` — auto-reconnect on ErrBrowserPanic
- `internal/app/tools.go` — wrapBrowserHandler applied to all browser tools
- `internal/gateway/server.go` — readPump/writePump/handleRPC panic recovery
- `docker-compose.yml` — Chrome service healthcheck
- `internal/tools/browser/panic_recovery_test.go` — new test file
