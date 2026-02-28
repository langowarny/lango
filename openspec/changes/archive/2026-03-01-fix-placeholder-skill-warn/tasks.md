## 1. Core Fix

- [x] 1.1 Add `strings` import to `internal/skill/file_store.go`
- [x] 1.2 Update `ListActive()` to skip directories starting with `.`
- [x] 1.3 Update `EnsureDefaults()` to skip embedded paths whose directory name starts with `.`

## 2. Verification

- [x] 2.1 Run `go build ./...` and confirm no build errors
- [x] 2.2 Run `go test ./internal/skill/...` and confirm all tests pass
