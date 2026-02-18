## ADDED Requirements

### Requirement: Hardware Signing Interface
The system SHALL provide an interface for cryptographic signing that delegates the actual signing operation to an external provider (Secure Enclave).

#### Scenario: Message Signing
- **WHEN** the application needs to sign a payload
- **THEN** it invokes the `Signer.Sign(payload)` method
- **AND** the payload is sent to the configured provider (e.g., macOS host app via RPC)
- **AND** the valid signature is returned

### Requirement: RPC Signer Provider
The system SHALL implement a Signer provider that communicates with a local host process (e.g., Swift app) via WebSocket/IPC to perform the signing.

#### Scenario: RPC Delegation
- **WHEN** `Signer.Sign` is called with the RPC provider configured
- **THEN** a `sign.request` message is sent over the IPC channel
- **AND** the system waits for a `sign.response` matching the request ID
