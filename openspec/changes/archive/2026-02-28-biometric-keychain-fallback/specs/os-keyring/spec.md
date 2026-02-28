## ADDED Requirements

### Requirement: ErrEntitlement sentinel for missing code signing
The keyring package SHALL export an `ErrEntitlement` sentinel error. `BiometricProvider.Get`, `Set`, and `Delete` SHALL wrap OSStatus `-34018` as `ErrEntitlement` using `fmt.Errorf %w`, allowing callers to match with `errors.Is(err, keyring.ErrEntitlement)`.

#### Scenario: BiometricProvider.Set returns ErrEntitlement on -34018
- **WHEN** `SecItemAdd` returns OSStatus `-34018`
- **THEN** the returned error SHALL satisfy `errors.Is(err, keyring.ErrEntitlement)`
- **AND** the error message SHALL contain `keychain biometric set:`

#### Scenario: BiometricProvider.Get returns ErrEntitlement on -34018
- **WHEN** `SecItemCopyMatching` returns OSStatus `-34018`
- **THEN** the returned error SHALL satisfy `errors.Is(err, keyring.ErrEntitlement)`
- **AND** the error message SHALL contain `keychain biometric get:`

#### Scenario: BiometricProvider.Delete returns ErrEntitlement on -34018
- **WHEN** `SecItemDelete` returns OSStatus `-34018`
- **THEN** the returned error SHALL satisfy `errors.Is(err, keyring.ErrEntitlement)`
- **AND** the error message SHALL contain `keychain biometric delete:`

#### Scenario: Non-entitlement OSStatus errors unchanged
- **WHEN** a biometric operation returns an OSStatus other than `-34018` (e.g., `-25308`)
- **THEN** the error SHALL NOT satisfy `errors.Is(err, keyring.ErrEntitlement)`
- **AND** the error message SHALL contain the numeric code and `osStatusDescription` output
