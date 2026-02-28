## Why

Biometric passphrase storage on macOS fails with `-34018 (errSecMissingEntitlement)` for ad-hoc signed binaries built with `go build`. Users must re-enter the passphrase on every launch because the Data Protection Keychain with biometric ACL requires proper Apple Developer code signing. Since `BiometricProvider` and `OSProvider` share the same macOS Keychain with the same service/account key, falling back to `OSProvider` (no biometric ACL) allows passphrase persistence without code signing while still enabling biometric protection for properly signed release builds.

## What Changes

- Add `ErrEntitlement` sentinel error for `-34018` detection via `errors.Is()`
- Wrap OSStatus `-34018` as `ErrEntitlement` in `BiometricProvider.Get/Set/Delete`
- Add `FallbackProvider` field to `passphrase.Options` for plain OS keyring fallback on read
- Bootstrap detects entitlement errors on biometric store and falls back to `OSProvider` with user-facing guidance
- CLI `keyring store` command applies the same fallback logic
- Add `build/entitlements.plist` with Keychain access groups for release code signing
- Add `make codesign` target for signing binaries with biometric Keychain entitlements

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `os-keyring`: Add `ErrEntitlement` sentinel; `BiometricProvider` wraps `-34018` as `ErrEntitlement`
- `passphrase-acquisition`: Add `FallbackProvider` for plain OS keyring read fallback
- `bootstrap-lifecycle`: Entitlement-aware fallback storage with user messaging

## Impact

- `internal/keyring/keyring.go` — `ErrEntitlement` sentinel
- `internal/keyring/biometric_darwin.go` — `-34018` → `ErrEntitlement` wrapping in Get/Set/Delete
- `internal/security/passphrase/acquire.go` — `FallbackProvider` field + fallback read path
- `internal/bootstrap/bootstrap.go` — fallback store logic + `FallbackProvider` wiring
- `internal/cli/security/keyring.go` — `keyring store` fallback
- `build/entitlements.plist` — NEW: macOS entitlements for Keychain access
- `Makefile` — `codesign` target
