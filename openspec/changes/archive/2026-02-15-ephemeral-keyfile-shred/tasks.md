## 1. Core: ShredKeyfile Function

- [x] 1.1 Add `ShredKeyfile(path string) error` to `internal/passphrase/keyfile.go` — zero-overwrite, sync, remove; return nil for nonexistent files
- [x] 1.2 Add unit tests for `ShredKeyfile` in `internal/passphrase/keyfile_test.go` — shred existing file, nonexistent file returns nil

## 2. Bootstrap: KeepKeyfile Option and Shred Integration

- [x] 2.1 Add `KeepKeyfile bool` field to `Options` in `internal/bootstrap/bootstrap.go`
- [x] 2.2 Capture passphrase `source` in `Run()` (change `pass, _, err` to `pass, source, err`)
- [x] 2.3 Call `ShredKeyfile()` after crypto init + checksum verification when `source == SourceKeyfile && !opts.KeepKeyfile`; log warning to stderr on failure

## 3. Bootstrap: Integration Tests

- [x] 3.1 Add `TestRun_ShredsKeyfileAfterCryptoInit` — verify keyfile deleted after bootstrap
- [x] 3.2 Add `TestRun_KeepsKeyfileWhenOptedOut` — verify keyfile persists with `KeepKeyfile: true`

## 4. Verification

- [x] 4.1 Run `go build ./...` and verify no build errors
- [x] 4.2 Run `go test ./internal/passphrase/... ./internal/bootstrap/...` and verify all tests pass
