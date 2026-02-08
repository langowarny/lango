## ADDED Requirements

### Requirement: Security Provider Mode Check
The system SHALL check security provider configuration and provide appropriate warnings.

#### Scenario: Local provider warning
- **WHEN** security.signer.provider is "local"
- **THEN** doctor displays WARNING: "Using LocalCryptoProvider (dev/test only). For production, use RPCProvider with Companion app."
- **AND** check status is "warn"

#### Scenario: RPC provider configured
- **WHEN** security.signer.provider is "rpc"
- **THEN** doctor displays "Security: Using RPCProvider (production mode)"
- **AND** check status is "pass"

#### Scenario: Companion connectivity with RPC
- **WHEN** security.signer.provider is "rpc"
- **AND** no companion is connected
- **THEN** doctor displays WARNING: "RPCProvider configured but no companion connected. Crypto operations will fail."
- **AND** check status is "warn"

### Requirement: Passphrase Checksum Integrity Check
The system SHALL verify that passphrase checksum exists when local provider is configured.

#### Scenario: Checksum present
- **WHEN** local provider is configured
- **AND** security_config contains valid checksum
- **THEN** doctor displays "Passphrase checksum: configured"
- **AND** check status is "pass"

#### Scenario: Checksum missing
- **WHEN** local provider is configured
- **AND** security_config has no checksum
- **THEN** doctor displays WARNING: "Passphrase not initialized. Run application to set up."
- **AND** check status is "warn"
