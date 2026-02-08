## ADDED Requirements

### Requirement: Local Fallback Provider
The system SHALL provide a local encryption fallback when companion app is unavailable.

#### Scenario: Initialize local provider with passphrase
- **WHEN** no companion is connected
- **AND** user provides a passphrase
- **THEN** the system SHALL derive an AES-256 key using PBKDF2
- **AND** use this key for local encryption/decryption

#### Scenario: Local encryption
- **WHEN** `Encrypt` is called with local provider
- **THEN** the system SHALL encrypt using AES-256-GCM
- **AND** prepend a random nonce to ciphertext

#### Scenario: Local decryption
- **WHEN** `Decrypt` is called with local provider
- **THEN** the system SHALL extract nonce from ciphertext
- **AND** decrypt using AES-256-GCM

### Requirement: Composite Provider Strategy
The system SHALL use a composite provider that tries companion first, then falls back to local.

#### Scenario: Companion available
- **WHEN** companion is connected
- **THEN** the system SHALL delegate crypto operations to companion via RPCProvider

#### Scenario: Companion unavailable with fallback
- **WHEN** companion is not connected
- **AND** local fallback is configured
- **THEN** the system SHALL use local provider

#### Scenario: No providers available
- **WHEN** companion is not connected
- **AND** local fallback is not configured
- **THEN** the system SHALL return an error "no crypto provider available"
