## Why

Docker 컨테이너에서 lango를 실행할 때 3가지 문제가 발생한다: (1) exec 도구의 approval chain이 TTY 없는 환경에서 전부 실패하여 도구 사용이 거부됨, (2) go-rod가 시스템 chromium 바이너리를 찾지 못해 브라우저 자동화 불가, (3) WORKDIR `/app`이 root 소유여서 non-root 사용자가 쓸 수 없음.

## What Changes

- Add `HeadlessProvider` to the approval system that auto-approves tool executions in headless environments with WARN-level audit logging
- Add `HeadlessAutoApprove` config field to `InterceptorConfig` (default: `false`, fail-closed)
- Add `BrowserBin` config field to `BrowserToolConfig` for explicit browser binary path
- Update `ensureBrowser()` to use `launcher.LookPath()` for automatic system chromium detection
- Change Dockerfile `WORKDIR` from `/app` to `/home/lango` (user home directory, writable)
- Remove unused `ENV ROD_BROWSER` from Dockerfile

## Capabilities

### New Capabilities
- `headless-approval`: Auto-approve provider for headless/Docker environments where no TTY or companion is available

### Modified Capabilities
- `ai-privacy-interceptor`: Add `HeadlessAutoApprove` config option for headless fallback behavior
- `tool-browser`: Add `BrowserBin` config and `LookPath()` auto-detection for system-installed browsers
- `docker-deployment`: Fix WORKDIR permissions and remove unused ROD_BROWSER env var

## Impact

- `internal/config/types.go`: New fields in `InterceptorConfig` and `BrowserToolConfig`
- `internal/approval/headless.go`: New file — `HeadlessProvider` implementation
- `internal/approval/headless_test.go`: New file — tests
- `internal/app/app.go`: Wiring for `HeadlessProvider` and `BrowserBin`
- `internal/tools/browser/browser.go`: `BrowserBin` field and `LookPath()` logic
- `Dockerfile`: WORKDIR change, ROD_BROWSER removal
