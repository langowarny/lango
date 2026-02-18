## ADDED Requirements

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

### Requirement: Composite Provider Strategy
The system SHALL use a composite provider that tries companion first, then falls back to local.

#### Scenario: Companion available
- **WHEN** companion is connected
- **THEN** the system SHALL delegate crypto operations to companion via RPCProvider

#### Scenario: Companion unavailable with fallback
- **WHEN** companion is not connected
- **AND** local fallback is configured
- **AND** terminal is interactive (TTY available)
- **THEN** the system SHALL use local provider

#### Scenario: No providers available
- **WHEN** companion is not connected
- **AND** local fallback is not configured
- **THEN** the system SHALL return an error "no crypto provider available"

#### Scenario: Docker environment detection
- **WHEN** the system detects it is running in a Docker container (/.dockerenv exists OR cgroup contains "docker")
- **AND** no companion is connected
- **THEN** the system SHALL log error "Docker environment requires RPC Provider. Please connect Companion app."
- **AND** SHALL NOT attempt to use LocalCryptoProvider
