# Passphrase Management

## Purpose
This capability defines how the user's passphrase for the Local Crypto Provider is securely handled, validated, and migrated. It ensures passphrases are never stored in plain text configuration files and provides mechanisms for key rotation.

## Requirements

### Requirement: Interactive Passphrase Prompt
The system SHALL prompt for passphrase via terminal instead of reading from config file.

#### Scenario: First-time passphrase setup
- **WHEN** LocalCryptoProvider is initialized
- **AND** no salt/checksum exists in database
- **THEN** system prompts "Enter new passphrase: " with hidden input
- **AND** prompts "Confirm passphrase: " for verification
- **AND** stores salt and checksum in database

#### Scenario: Subsequent passphrase entry
- **WHEN** LocalCryptoProvider is initialized
- **AND** salt/checksum exists in database
- **THEN** system prompts "Enter passphrase: " with hidden input
- **AND** validates against stored checksum

#### Scenario: Passphrase mismatch
- **WHEN** user enters incorrect passphrase
- **THEN** system displays "Passphrase does not match. Please try again."
- **AND** allows up to 3 retry attempts
- **AND** exits with error after 3 failures

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
The system SHALL ignore passphrase in config file and warn user.

#### Scenario: Passphrase in config detected
- **WHEN** config file contains security.passphrase field
- **THEN** system logs WARNING: "security.passphrase in config is deprecated and ignored. Passphrase will be prompted interactively."
- **AND** does NOT use the config value
