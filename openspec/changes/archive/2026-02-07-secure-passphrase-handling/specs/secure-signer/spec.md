## MODIFIED Requirements

### Requirement: Local Fallback Provider
The system SHALL provide a local encryption fallback when companion app is unavailable.

#### Scenario: Initialize local provider with passphrase
- **WHEN** no companion is connected
- **AND** security.signer.provider is "local"
- **THEN** system prompts for passphrase via interactive terminal
- **AND** derives an AES-256 key using PBKDF2
- **AND** uses this key for local encryption/decryption

#### Scenario: Local encryption
- **WHEN** `Encrypt` is called with local provider
- **THEN** the system SHALL encrypt using AES-256-GCM
- **AND** prepend a random nonce to ciphertext

#### Scenario: Local decryption
- **WHEN** `Decrypt` is called with local provider
- **THEN** the system SHALL extract nonce from ciphertext
- **AND** decrypt using AES-256-GCM

#### Scenario: Headless environment detection
- **WHEN** terminal is not interactive (no TTY)
- **AND** local provider is configured
- **THEN** system exits with error "LocalCryptoProvider requires interactive terminal. Use RPCProvider with Companion for headless environments."
