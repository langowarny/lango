# Fix Serve Encryption Spec

## Goal
Ensure `lango serve` can successfully open and use an encrypted session database.

## Requirements

### Requirement: Encryption Support
The application SHALL be able to open an encrypted SQLite database if the correct passphrase is provided.

#### Scenario: Server Startup with Passphrase
- **GIVEN** an encrypted session database exists at configured path
- **AND** `LANGO_PASSPHRASE` environment variable is set to the correct key
- **WHEN** `lango serve` is executed
- **THEN** the application starts successfully
- **AND** the session store is accessible (no "out of memory" or "not a database" errors)

### Requirement: Passphrase Configuration
The application SHALL prioritize the passphrase from environment variables over configuration files (standard security practice).

### Requirement: Path Expansion
The application SHALL verify that configuration paths using `~` are correctly expanded to the user's home directory.

#### Scenario: Tilde Expansion
- **GIVEN** `databasePath` is configured as `~/.lango/sessions.db`
- **WHEN** the application initializes storage
- **THEN** it expands `~` to the current user's home directory
- **AND** successfully locates the file/directory
