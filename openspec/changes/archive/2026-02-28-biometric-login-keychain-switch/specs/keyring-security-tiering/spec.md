## MODIFIED Requirements

### Requirement: BiometricProvider uses macOS Keychain with Touch ID ACL
The system SHALL provide a `BiometricProvider` that stores secrets in the macOS login Keychain (NOT the Data Protection Keychain) using `kSecAccessControlBiometryCurrentSet` access control with `kSecAttrAccessibleWhenPasscodeSetThisDeviceOnly` protection. All Keychain queries SHALL set `kSecUseDataProtectionKeychain = kCFBooleanFalse` to explicitly target the login Keychain. This provider SHALL require Touch ID authentication for every read operation, and SHALL invalidate stored items when biometric enrollment changes.

#### Scenario: Store and retrieve with biometric
- **WHEN** a secret is stored via `BiometricProvider.Set()` and later retrieved via `BiometricProvider.Get()`
- **THEN** the Set SHALL create a login Keychain item with `BiometryCurrentSet` ACL and `kSecUseDataProtectionKeychain = false`, and Get SHALL trigger Touch ID before returning the value

#### Scenario: Biometric not available on non-Darwin platform
- **WHEN** `NewBiometricProvider()` is called on a non-Darwin or non-CGO platform
- **THEN** it SHALL return `ErrBiometricNotAvailable`

#### Scenario: Ad-hoc signed binary works without entitlement
- **WHEN** a `go build` ad-hoc signed binary calls `BiometricProvider.Set()` or `BiometricProvider.Get()`
- **THEN** the operation SHALL succeed without requiring `keychain-access-groups` entitlement

#### Scenario: Fingerprint enrollment change invalidates stored items
- **WHEN** a user changes their biometric enrollment (adds or removes fingerprints) after storing a secret
- **THEN** attempts to retrieve the secret SHALL fail because `BiometryCurrentSet` invalidates the access control

#### Scenario: Device passcode not set
- **WHEN** the device does not have a passcode configured
- **THEN** `NewBiometricProvider()` SHALL return `ErrBiometricNotAvailable` because the Keychain probe will fail

## ADDED Requirements

### Requirement: BiometricProvider availability probe uses real Keychain write
The `keychain_biometric_available` function SHALL verify biometric support by performing a real `SecItemAdd` probe to the login Keychain with biometric ACL, rather than only checking `SecAccessControlCreateWithFlags`. The probe item SHALL be cleaned up immediately after the test.

#### Scenario: Probe succeeds on capable hardware
- **WHEN** `keychain_biometric_available()` is called on a macOS device with Touch ID and device passcode set
- **THEN** it SHALL add a probe item to the login Keychain, delete it, and return 1

#### Scenario: Probe fails without passcode
- **WHEN** `keychain_biometric_available()` is called on a macOS device without a passcode
- **THEN** the `SecItemAdd` SHALL fail and the function SHALL return 0

#### Scenario: Probe does not trigger Touch ID
- **WHEN** the probe item is added via `SecItemAdd`
- **THEN** it SHALL NOT trigger a Touch ID prompt because Keychain writes bypass ACL evaluation

### Requirement: All Keychain queries target login Keychain explicitly
Every Keychain query dictionary (set, get, has, delete) SHALL include `kSecUseDataProtectionKeychain = kCFBooleanFalse` to ensure operations target the login Keychain and never fall through to the Data Protection Keychain.

#### Scenario: Set targets login Keychain
- **WHEN** `keychain_set_biometric()` builds its query dictionaries
- **THEN** both the delete-existing and add-new dictionaries SHALL include `kSecUseDataProtectionKeychain = kCFBooleanFalse`

#### Scenario: Get targets login Keychain
- **WHEN** `keychain_get_biometric()` builds its query dictionary
- **THEN** it SHALL include `kSecUseDataProtectionKeychain = kCFBooleanFalse`

#### Scenario: Has targets login Keychain
- **WHEN** `keychain_has_biometric()` builds its query dictionary
- **THEN** it SHALL include `kSecUseDataProtectionKeychain = kCFBooleanFalse`

#### Scenario: Delete targets login Keychain
- **WHEN** `keychain_delete_biometric()` builds its query dictionary
- **THEN** it SHALL include `kSecUseDataProtectionKeychain = kCFBooleanFalse`
