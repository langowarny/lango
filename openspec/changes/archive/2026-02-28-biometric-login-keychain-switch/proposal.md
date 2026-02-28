## Why

The `BiometricProvider` targets the macOS Data Protection Keychain via `kSecAttrAccessibleWhenUnlockedThisDeviceOnly`, which requires `keychain-access-groups` entitlement. This means `go build` ad-hoc signed binaries fail with `-34018 (errSecMissingEntitlement)`. Switching to the login Keychain with `kSecAttrAccessControl` + `BiometryCurrentSet` provides Touch ID protection without entitlement requirements, making biometric storage work out of the box.

## What Changes

- Switch all Keychain queries from Data Protection Keychain to login Keychain (`kSecUseDataProtectionKeychain = false`)
- Change access control from `kSecAttrAccessibleWhenUnlockedThisDeviceOnly` + `BiometryAny` to `kSecAttrAccessibleWhenPasscodeSetThisDeviceOnly` + `BiometryCurrentSet`
- Replace simple `SecAccessControlCreateWithFlags` availability check with a real Keychain probe (SecItemAdd + cleanup) for accurate detection
- Update error messages to mention device passcode requirement
- Update Makefile `codesign` target description from required to optional enhancement

## Capabilities

### New Capabilities

(none)

### Modified Capabilities
- `keyring-security-tiering`: BiometricProvider now targets login Keychain instead of Data Protection Keychain; access control changed to BiometryCurrentSet; entitlement no longer required for biometric tier

## Impact

- `internal/keyring/biometric_darwin.go` — All C functions updated (set/get/has/delete/available)
- `internal/keyring/keyring.go` — ErrEntitlement doc comment updated
- `internal/cli/security/keyring.go` — Error messages improved
- `internal/bootstrap/bootstrap.go` — Error messages improved
- `Makefile` — codesign target description changed
- `build/entitlements.plist` — Retained for optional release builds
