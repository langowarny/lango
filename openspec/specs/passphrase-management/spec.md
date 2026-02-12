# Passphrase Management

## Purpose
This capability defines how the user's passphrase for the Local Crypto Provider is securely handled, validated, and migrated. It ensures passphrases are never stored in plain text configuration files and provides mechanisms for key rotation.

## Requirements

### Requirement: Interactive Passphrase Prompt
The passphrase prompt SHALL only be triggered when `security.signer.provider` is explicitly set to `"local"`. When security is not configured, the system SHALL skip passphrase initialization entirely without error. When security IS configured but passphrase initialization fails, the system SHALL log a warning and continue startup without security tools (non-blocking).

#### Scenario: First-time passphrase setup
- **WHEN** `security.signer.provider` is `"local"` and no salt exists and environment is interactive
- **THEN** the system SHALL prompt for passphrase creation
- **THEN** the system SHALL store salt and checksum

#### Scenario: Security not configured
- **WHEN** `security.signer.provider` is empty or not set
- **THEN** the system SHALL skip all passphrase initialization
- **THEN** the system SHALL log an info message about security being disabled
- **THEN** the agent SHALL start normally without security tools

#### Scenario: Passphrase initialization failure
- **WHEN** `security.signer.provider` is `"local"` but passphrase cannot be obtained (non-interactive, no env var)
- **THEN** the system SHALL log a warning
- **THEN** the system SHALL continue startup without security tools
- **THEN** the system SHALL NOT return an error or block startup

### Requirement: Passphrase Checksum Validation
The system SHALL store a checksum to detect incorrect passphrase early.

#### Scenario: Checksum storage
- **WHEN** new passphrase is set
- **THEN** system stores SHA256(passphrase + salt) in security_config table

#### Scenario: Checksum verification
- **WHEN** passphrase is entered
- **THEN** system computes SHA256(passphrase + salt)
- **AND** compares with stored checksum
- **AND** rejects if mismatch before attempting any decryption

### Requirement: Passphrase Migration Command
The system SHALL provide a CLI command to migrate encrypted data to a new passphrase.

#### Scenario: Successful migration
- **WHEN** user runs `lango security migrate-passphrase`
- **AND** enters correct current passphrase
- **AND** enters new passphrase (with confirmation)
- **THEN** system decrypts all secrets with old key
- **AND** re-encrypts with new key
- **AND** updates salt and checksum

#### Scenario: Migration with wrong current passphrase
- **WHEN** user enters incorrect current passphrase
- **THEN** system displays error and aborts without modifying data

#### Scenario: Migration rollback on failure
- **WHEN** re-encryption fails mid-process
- **THEN** system rolls back all changes
- **AND** preserves original encrypted data

### Requirement: Config Passphrase Deprecation
The `security.passphrase` config field SHALL be ignored. The `LANGO_PASSPHRASE` environment variable SHALL be the only supported method for providing passphrase non-interactively. A deprecation warning SHALL be logged if `security.passphrase` is set in config.

#### Scenario: Passphrase in config detected
- **WHEN** `security.passphrase` is set in config file
- **THEN** the system SHALL log a deprecation warning
- **THEN** the system SHALL NOT use the config value
