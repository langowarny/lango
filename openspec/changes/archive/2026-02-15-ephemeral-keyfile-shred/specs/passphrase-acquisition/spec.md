## MODIFIED Requirements

### Requirement: Keyfile management
The system SHALL read, write, and securely shred keyfiles with strict 0600 permission enforcement.

#### Scenario: Write keyfile
- **WHEN** a keyfile is written
- **THEN** the file is created with 0600 permissions and parent directories with 0700

#### Scenario: Read keyfile with valid permissions
- **WHEN** a keyfile is read with 0600 permissions
- **THEN** the passphrase is returned with trailing whitespace trimmed

#### Scenario: Read keyfile with invalid permissions
- **WHEN** a keyfile exists with permissions other than 0600
- **THEN** the system returns a permission validation error

#### Scenario: Shred keyfile
- **WHEN** `ShredKeyfile()` is called on an existing keyfile
- **THEN** the file content is overwritten with zero bytes, synced, and removed

#### Scenario: Shred nonexistent keyfile
- **WHEN** `ShredKeyfile()` is called on a nonexistent file
- **THEN** nil is returned without error
