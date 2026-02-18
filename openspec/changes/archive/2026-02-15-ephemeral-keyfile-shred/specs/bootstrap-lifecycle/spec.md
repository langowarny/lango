## ADDED Requirements

### Requirement: Ephemeral keyfile shredding after crypto initialization
The system SHALL shred the passphrase keyfile after successful crypto initialization and checksum verification when the passphrase source is keyfile and `KeepKeyfile` is false (default). Shred failure SHALL emit a warning to stderr but SHALL NOT prevent bootstrap from completing.

#### Scenario: Keyfile shredded after successful bootstrap
- **WHEN** the passphrase source is `SourceKeyfile` and `KeepKeyfile` is false
- **AND** crypto initialization and checksum verification succeed
- **THEN** the keyfile is securely shredded and no longer exists on disk

#### Scenario: Keyfile kept when opted out
- **WHEN** the passphrase source is `SourceKeyfile` and `KeepKeyfile` is true
- **THEN** the keyfile remains on disk after bootstrap

#### Scenario: Non-keyfile source unaffected
- **WHEN** the passphrase source is `SourceInteractive` or `SourceStdin`
- **THEN** no shredding is attempted regardless of `KeepKeyfile` value

#### Scenario: Shred failure is non-fatal
- **WHEN** `ShredKeyfile()` returns an error during bootstrap
- **THEN** a warning is printed to stderr and bootstrap continues with the already-initialized crypto provider

## MODIFIED Requirements

### Requirement: Unified bootstrap sequence
The system SHALL execute a complete bootstrap sequence: ensure data directory → open database → acquire passphrase → initialize crypto → shred keyfile (if applicable) → load config profile. The result SHALL be a single struct containing all initialized components. The `Options` struct SHALL NOT include a `MigrationPath` field. The `Options` struct SHALL include a `KeepKeyfile bool` field that defaults to false (secure by default).

#### Scenario: First-run bootstrap
- **WHEN** no salt exists in the database (first run)
- **THEN** the system acquires a new passphrase (with confirmation), generates a salt, stores the checksum, shreds the keyfile if source is keyfile, creates a default config profile, and returns the Result

#### Scenario: Returning-user bootstrap
- **WHEN** salt and checksum exist in the database
- **THEN** the system acquires the passphrase, verifies it against the stored checksum, shreds the keyfile if source is keyfile, and loads the active profile

#### Scenario: Wrong passphrase on returning user
- **WHEN** the user provides an incorrect passphrase for an existing database
- **THEN** the system returns a "passphrase checksum mismatch" error and the keyfile is NOT shredded

#### Scenario: No profiles exist
- **WHEN** no profiles exist in the database
- **THEN** the system creates a default profile with `config.DefaultConfig()` and sets it as active
