## MODIFIED Requirements

### Requirement: Passphrase acquisition priority chain
The system SHALL acquire a passphrase using the following priority: (1) hardware keyring (Touch ID / TPM), (2) keyfile at `~/.lango/keyfile`, (3) interactive terminal prompt, (4) stdin pipe. The system SHALL return an error if no source is available.

#### Scenario: Keyring provider returns passphrase
- **WHEN** `Acquire()` is called with a non-nil `KeyringProvider` that returns a valid passphrase
- **THEN** the passphrase is returned with `SourceKeyring` and no further sources are tried

#### Scenario: Keyfile exists with correct permissions
- **WHEN** a keyfile exists at the configured path with 0600 permissions
- **THEN** the passphrase is read from the file and `SourceKeyfile` is returned

#### Scenario: Keyfile has wrong permissions
- **WHEN** a keyfile exists but does not have 0600 permissions
- **THEN** the keyfile is skipped and the next source is tried

#### Scenario: Interactive terminal available
- **WHEN** no keyfile is available and stdin is a terminal
- **THEN** the user is prompted for a passphrase via hidden input and `SourceInteractive` is returned

#### Scenario: New passphrase creation
- **WHEN** `AllowCreation` is true and interactive terminal is used
- **THEN** the user is prompted twice (entry + confirmation) and the passphrase must match

#### Scenario: Stdin pipe
- **WHEN** no keyfile is available and stdin is a pipe (not a terminal)
- **THEN** one line is read from stdin and `SourceStdin` is returned

#### Scenario: No source available
- **WHEN** no keyfile exists, stdin is not a terminal, and stdin pipe is empty
- **THEN** the system returns an error
