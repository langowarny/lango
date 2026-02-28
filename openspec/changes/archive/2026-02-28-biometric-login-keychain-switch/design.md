## Context

The `BiometricProvider` currently targets the macOS Data Protection Keychain by using `kSecAttrAccessibleWhenUnlockedThisDeviceOnly` as the protection level in `SecAccessControlCreateWithFlags`. The Data Protection Keychain requires the binary to be signed with a `keychain-access-groups` entitlement, which means ad-hoc signed binaries from `go build` fail with OSStatus `-34018 (errSecMissingEntitlement)`.

macOS provides an alternative: the login Keychain, which accepts `kSecAttrAccessControl` with biometric flags without requiring entitlements. By explicitly opting out of the Data Protection Keychain (`kSecUseDataProtectionKeychain = false`) and using `kSecAttrAccessibleWhenPasscodeSetThisDeviceOnly` + `kSecAccessControlBiometryCurrentSet`, Touch ID protection is achieved without code signing requirements.

## Goals / Non-Goals

**Goals:**
- Make biometric keyring storage work with ad-hoc signed (`go build`) binaries
- Improve security by switching from `BiometryAny` to `BiometryCurrentSet` (invalidates items on fingerprint enrollment changes)
- Add a real Keychain probe in `keychain_biometric_available` to detect entitlement issues at detection time rather than at first use
- Keep codesign as an optional enhancement for Data Protection Keychain access in release builds

**Non-Goals:**
- Changing the Provider interface or adding new SecurityTier values
- Modifying TPMProvider behavior
- Re-introducing `go-keyring` dependency
- Changing the biometric provider's Go-level API

## Decisions

### Decision 1: Login Keychain via `kSecUseDataProtectionKeychain = false`

All Keychain queries (set/get/has/delete) explicitly set `kSecUseDataProtectionKeychain = kCFBooleanFalse` to force operations to the login Keychain.

**Rationale**: The login Keychain does not require `keychain-access-groups` entitlement. Biometric ACL (`kSecAttrAccessControl`) still enforces Touch ID authentication for reads. This is the simplest path to making ad-hoc binaries work.

**Alternative considered**: Keeping Data Protection Keychain and documenting that codesign is required. Rejected because it creates a poor developer experience and blocks immediate use after `go build`.

### Decision 2: `BiometryCurrentSet` instead of `BiometryAny`

Changed from `kSecAccessControlBiometryAny` to `kSecAccessControlBiometryCurrentSet`.

**Rationale**: `BiometryCurrentSet` invalidates stored items when the biometric enrollment changes (fingerprints added/removed). This prevents a scenario where an attacker adds their fingerprint and accesses previously stored secrets. Strictly more secure.

### Decision 3: `kSecAttrAccessibleWhenPasscodeSetThisDeviceOnly` protection level

Changed from `kSecAttrAccessibleWhenUnlockedThisDeviceOnly` to `kSecAttrAccessibleWhenPasscodeSetThisDeviceOnly`.

**Rationale**: This requires a device passcode to be set, which is a prerequisite for biometric authentication anyway. Items are still excluded from backups (`ThisDeviceOnly`). The key difference is that if the user removes their passcode, the items become inaccessible — a desirable security property.

### Decision 4: Real Keychain probe in availability check

The `keychain_biometric_available` function now performs a real `SecItemAdd` + `SecItemDelete` probe instead of just checking `SecAccessControlCreateWithFlags`.

**Rationale**: `SecAccessControlCreateWithFlags` succeeds even when the Keychain won't accept items (e.g., missing entitlements, passcode not set). The probe catches these issues at detection time, so `DetectSecureProvider` returns `(nil, TierNone)` instead of returning a provider that fails on first use.

## Risks / Trade-offs

- **[Risk] Login Keychain items are not hardware-encrypted by Secure Enclave** → The Keychain file itself is still encrypted. Touch ID ACL still enforces biometric authentication. The practical security difference is minimal for passphrase storage.
- **[Risk] Device passcode not set** → `kSecAttrAccessibleWhenPasscodeSetThisDeviceOnly` will fail if no passcode is configured. The probe detects this, and error messages now mention this requirement.
- **[Risk] Fingerprint enrollment change invalidates stored passphrase** → This is intentional security behavior. Users will need to re-store the passphrase after changing fingerprints. Error messages should guide users.
- **[Trade-off] Probe adds ~10ms to startup** → Only runs once during `DetectSecureProvider`. Acceptable for a one-time detection.
