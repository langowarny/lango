## 1. ErrEntitlement Sentinel

- [x] 1.1 Add `ErrEntitlement` sentinel error to `internal/keyring/keyring.go`
- [x] 1.2 Wrap OSStatus `-34018` as `ErrEntitlement` in `BiometricProvider.Get`
- [x] 1.3 Wrap OSStatus `-34018` as `ErrEntitlement` in `BiometricProvider.Set`
- [x] 1.4 Wrap OSStatus `-34018` as `ErrEntitlement` in `BiometricProvider.Delete`

## 2. Passphrase Acquisition Fallback

- [x] 2.1 Add `FallbackProvider keyring.Provider` field to `passphrase.Options`
- [x] 2.2 Add fallback read path in `Acquire()` between primary keyring and keyfile

## 3. Bootstrap Fallback Storage

- [x] 3.1 Wire `OSProvider` as `FallbackProvider` on macOS when biometric is detected
- [x] 3.2 Detect `ErrEntitlement` on biometric store failure and fall back to `OSProvider.Set()`
- [x] 3.3 Emit user-facing messages: warning, fallback result, codesign guidance

## 4. CLI Keyring Store Fallback

- [x] 4.1 Add `ErrEntitlement` detection to `keyring store` command
- [x] 4.2 Fall back to `OSProvider.Set()` with user-facing messages

## 5. Codesign Infrastructure

- [x] 5.1 Create `build/entitlements.plist` with Keychain access groups
- [x] 5.2 Add `codesign` target to Makefile with `APPLE_IDENTITY` requirement

## 6. Verification

- [x] 6.1 Run `go build ./...` and confirm no compilation errors
- [x] 6.2 Run `go test ./internal/keyring/... ./internal/bootstrap/... ./internal/security/...` and confirm all pass
