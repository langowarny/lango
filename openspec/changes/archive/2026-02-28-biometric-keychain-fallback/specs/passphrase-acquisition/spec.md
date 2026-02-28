## ADDED Requirements

### Requirement: FallbackProvider for plain OS keyring read
The `Options` struct SHALL include a `FallbackProvider keyring.Provider` field. When set, `Acquire()` SHALL attempt to read from `FallbackProvider` after `KeyringProvider` fails and before trying keyfile. Non-`ErrNotFound` errors from `FallbackProvider` SHALL be logged to stderr as `warning: fallback keyring read failed: <error>`.

#### Scenario: Primary fails, fallback succeeds
- **WHEN** `KeyringProvider.Get()` returns `ErrNotFound` and `FallbackProvider.Get()` returns a passphrase
- **THEN** `Acquire()` SHALL return the passphrase with `SourceKeyring`

#### Scenario: Primary fails, fallback also fails
- **WHEN** both `KeyringProvider.Get()` and `FallbackProvider.Get()` return `ErrNotFound`
- **THEN** `Acquire()` SHALL proceed to keyfile → interactive → stdin

#### Scenario: Fallback provider is nil
- **WHEN** `FallbackProvider` is nil
- **THEN** the fallback step SHALL be skipped entirely

#### Scenario: Fallback read error logged
- **WHEN** `FallbackProvider.Get()` returns a non-`ErrNotFound` error
- **THEN** stderr SHALL contain `warning: fallback keyring read failed: <error detail>`
- **AND** acquisition SHALL continue to keyfile
