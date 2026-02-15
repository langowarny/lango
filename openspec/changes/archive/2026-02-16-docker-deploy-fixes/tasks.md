## 1. Dockerfile Fixes

- [x] 1.1 Change WORKDIR from `/app` to `/home/lango`
- [x] 1.2 Remove unused `ENV ROD_BROWSER=/usr/bin/chromium`
- [x] 1.3 Add `chmod +x` for docker-entrypoint.sh

## 2. Browser Binary Auto-Detection

- [x] 2.1 Add `BrowserBin` field to `BrowserToolConfig` in `internal/config/types.go`
- [x] 2.2 Add `BrowserBin` field to `browser.Config` struct in `internal/tools/browser/browser.go`
- [x] 2.3 Update `ensureBrowser()` to use `launcher.LookPath()` with `BrowserBin` config fallback
- [x] 2.4 Wire `BrowserBin` from config to `browser.Config` in `internal/app/app.go`

## 3. Headless Approval Provider

- [x] 3.1 Add `HeadlessAutoApprove` field to `InterceptorConfig` in `internal/config/types.go`
- [x] 3.2 Create `HeadlessProvider` in `internal/approval/headless.go`
- [x] 3.3 Wire `HeadlessProvider` as TTY fallback when `HeadlessAutoApprove` is true in `internal/app/app.go`
- [x] 3.4 Create tests for `HeadlessProvider` in `internal/approval/headless_test.go`

## 4. Verification

- [x] 4.1 Run `go build ./...` — pass
- [x] 4.2 Run `go test ./internal/approval/...` — pass
