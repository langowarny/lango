## 1. Core Keychain C Functions

- [x] 1.1 Update `keychain_set_biometric` to use `kSecAttrAccessibleWhenPasscodeSetThisDeviceOnly` + `kSecAccessControlBiometryCurrentSet` and add `kSecUseDataProtectionKeychain = false` to all query dictionaries
- [x] 1.2 Update `keychain_get_biometric` to add `kSecUseDataProtectionKeychain = false`
- [x] 1.3 Update `keychain_has_biometric` to add `kSecUseDataProtectionKeychain = false`
- [x] 1.4 Update `keychain_delete_biometric` to add `kSecUseDataProtectionKeychain = false`
- [x] 1.5 Replace `keychain_biometric_available` with real Keychain probe (SecItemAdd + cleanup)

## 2. Go Layer Updates

- [x] 2.1 Update `BiometricProvider` struct doc comment to reflect login Keychain and BiometryCurrentSet
- [x] 2.2 Update `Set` method doc comment
- [x] 2.3 Update `osStatusDescription` to add `-25291` (passcode not set) and improve `-25293` description
- [x] 2.4 Update `ErrEntitlement` doc comment in `keyring.go`

## 3. Error Messages

- [x] 3.1 Update `internal/cli/security/keyring.go` — add passcode requirement note to entitlement error
- [x] 3.2 Update `internal/bootstrap/bootstrap.go` — add passcode requirement note to entitlement warning

## 4. Build System

- [x] 4.1 Update Makefile `codesign` target description from required to optional enhancement

## 5. Verification

- [x] 5.1 Run `go build ./...` — verify compilation succeeds
- [x] 5.2 Run `go test ./internal/keyring/...` — verify all tests pass
