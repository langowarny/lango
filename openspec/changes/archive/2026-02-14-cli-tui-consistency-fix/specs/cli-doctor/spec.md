## MODIFIED Requirements

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

#### Scenario: Enclave provider configured
- **WHEN** security.signer.provider is "enclave"
- **THEN** doctor SHALL NOT display any warnings
- **AND** check status is "pass"

#### Scenario: Unknown provider
- **WHEN** security.signer.provider is an unrecognized value
- **THEN** doctor displays error: "Unknown security provider: <value>"
- **AND** check status is "fail"

#### Scenario: Companion connectivity with RPC
- **WHEN** security.signer.provider is "rpc"
- **AND** no companion is connected
- **THEN** doctor displays WARNING: "RPCProvider configured but no companion connected. Crypto operations will fail."
- **AND** check status is "warn"

## ADDED Requirements

### Requirement: Doctor command description lists all checks
The `lango doctor` command Long description SHALL enumerate all diagnostic checks performed.

#### Scenario: Long description content
- **WHEN** user runs `lango doctor --help`
- **THEN** the description SHALL list: Configuration file validity, API key and provider configuration, Channel token validation, Session database accessibility, Server port availability, Security configuration, Companion connectivity
