## MODIFIED Requirements

### Requirement: Report biometric passphrase store outcome
When the bootstrap flow stores a passphrase in the secure keyring provider, it SHALL report the outcome to stderr. On entitlement error (`ErrEntitlement`), the system SHALL fall back to `OSProvider` (plain macOS Keychain without biometric ACL) and report the fallback status. On other failures, the message SHALL be `warning: store passphrase failed: <error>`. On success, the message SHALL be `Passphrase stored successfully.`.

#### Scenario: Biometric store succeeds
- **WHEN** `secureProvider.Set()` returns nil
- **THEN** stderr SHALL contain `Passphrase stored successfully.`

#### Scenario: Biometric store fails with entitlement error
- **WHEN** `secureProvider.Set()` returns an error satisfying `errors.Is(err, keyring.ErrEntitlement)`
- **THEN** stderr SHALL contain `warning: biometric storage unavailable (binary not codesigned)`
- **AND** the system SHALL attempt `OSProvider.Set()` as fallback

#### Scenario: Entitlement fallback to OSProvider succeeds
- **WHEN** biometric store fails with `ErrEntitlement` and `OSProvider.Set()` succeeds
- **THEN** stderr SHALL contain `Passphrase stored in macOS Keychain (without biometric protection).`
- **AND** stderr SHALL contain `For biometric protection, codesign the binary: make codesign`

#### Scenario: Entitlement fallback to OSProvider fails
- **WHEN** biometric store fails with `ErrEntitlement` and `OSProvider.Set()` also fails
- **THEN** stderr SHALL contain `warning: fallback keychain store also failed: <error>`

#### Scenario: Biometric store fails with non-entitlement error
- **WHEN** `secureProvider.Set()` returns an error NOT satisfying `errors.Is(err, keyring.ErrEntitlement)`
- **THEN** stderr SHALL contain `warning: store passphrase failed: <error detail>`

## ADDED Requirements

### Requirement: FallbackProvider wiring for macOS
On macOS, when a secure hardware provider (biometric) is detected, bootstrap SHALL create an `OSProvider` as the `FallbackProvider` in passphrase acquisition options. This enables reading passphrase items stored without biometric ACL.

#### Scenario: macOS with biometric provider
- **WHEN** bootstrap runs on macOS with a detected biometric provider
- **THEN** `passphrase.Acquire()` SHALL receive an `OSProvider` as `FallbackProvider`

#### Scenario: Non-macOS or no secure provider
- **WHEN** bootstrap runs on Linux or with no secure provider
- **THEN** `FallbackProvider` SHALL be nil
