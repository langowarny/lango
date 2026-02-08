## RENAMED Requirements

FROM: ### Requirement: Hardware Signing Interface
TO: ### Requirement: Hardware Crypto Interface
The system SHALL provide an interface for cryptographic operations that delegates the actual signing, encryption, and decryption operations to an external provider (Secure Enclave).

#### Scenario: Message Signing
- **WHEN** the application needs to sign a payload
- **THEN** it invokes the `CryptoProvider.Sign(payload)` method
- **AND** the payload is sent to the configured provider (e.g., macOS host app via RPC)
- **AND** the valid signature is returned

## ADDED Requirements

### Requirement: Hardware Encryption
The system SHALL provide a method to encrypt data using hardware-backed keys without exposing the private key to the application.

#### Scenario: Data Encryption
- **WHEN** the application invokes `CryptoProvider.Encrypt(plaintext)`
- **THEN** an `encrypt.request` is sent via RPC to the host provider
- **AND** the provider returns the encrypted ciphertext
- **AND** the application never sees the encryption key

### Requirement: Hardware Decryption
The system SHALL provide a method to decrypt data that was encrypted by the hardware key, returning the plaintext only to the authorized caller.

#### Scenario: Data Decryption
- **WHEN** the application invokes `CryptoProvider.Decrypt(ciphertext)`
- **THEN** a `decrypt.request` is sent via RPC to the host provider
- **AND** the provider returns the decrypted plaintext

## MODIFIED Requirements

### Requirement: RPC Signer Provider
The system SHALL implement a CryptoProvider that communicates with a local host process (e.g., Swift app) via WebSocket/IPC to perform signing, encryption, and decryption.

#### Scenario: RPC Delegation
- **WHEN** `Sign`, `Encrypt`, or `Decrypt` is called with the RPC provider configured
- **THEN** a request message with the corresponding type (`sign.request`, `encrypt.request`, `decrypt.request`) is sent over the IPC channel
- **AND** the system waits for a response matching the request ID
