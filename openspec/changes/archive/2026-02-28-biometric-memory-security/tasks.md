## 1. C-Level Memory Security

- [x] 1.1 Add `secure_free` C helper with volatile pointer zeroing to `biometric_darwin.go`
- [x] 1.2 Update `Get()` to use `C.GoBytes` + `C.secure_free` instead of `C.GoStringN` + `C.free`
- [x] 1.3 Add Go `[]byte` zeroing loop after string extraction in `Get()`

## 2. Set() Buffer Zeroing

- [x] 2.1 Replace `defer C.free(cValue)` with `defer` that calls `C.memset` + `C.free` in `Set()`

## 3. Documentation

- [x] 3.1 Update `SourceKeyring` comment in `passphrase/acquire.go` to "hardware keyring (Touch ID / TPM)"

## 4. Verification

- [x] 4.1 Run `go build ./...` — confirm no compilation errors
- [x] 4.2 Run `go test ./internal/keyring/...` — confirm all tests pass
- [x] 4.3 Run `go test ./internal/security/passphrase/...` — confirm all tests pass
