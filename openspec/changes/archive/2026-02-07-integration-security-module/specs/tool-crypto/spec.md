## ADDED Requirements

### Requirement: Crypto Encrypt Operation
The system SHALL provide a `crypto.encrypt` tool operation for data encryption.

#### Scenario: Encrypt data with specified key
- **WHEN** AI agent calls `crypto.encrypt` with data and keyId
- **THEN** the system encrypts the data using the specified key via CryptoProvider
- **AND** returns base64-encoded ciphertext

#### Scenario: Encrypt with default key
- **WHEN** AI agent calls `crypto.encrypt` with data but no keyId
- **THEN** the system uses the default encryption key
- **AND** returns base64-encoded ciphertext

### Requirement: Crypto Decrypt Operation
The system SHALL provide a `crypto.decrypt` tool operation for data decryption.

#### Scenario: Decrypt data with specified key
- **WHEN** AI agent calls `crypto.decrypt` with ciphertext and keyId
- **THEN** the system decrypts the data using the specified key via CryptoProvider
- **AND** returns plaintext

#### Scenario: Decrypt with wrong key
- **WHEN** AI agent calls `crypto.decrypt` with ciphertext and incorrect keyId
- **THEN** the system returns a decryption error

### Requirement: Crypto Sign Operation
The system SHALL provide a `crypto.sign` tool operation for digital signatures.

#### Scenario: Sign data
- **WHEN** AI agent calls `crypto.sign` with data and keyId
- **THEN** the system generates a signature using the specified key via CryptoProvider
- **AND** returns base64-encoded signature

### Requirement: Crypto Hash Operation
The system SHALL provide a `crypto.hash` tool operation for hash generation.

#### Scenario: Generate SHA-256 hash
- **WHEN** AI agent calls `crypto.hash` with data and algorithm "sha256"
- **THEN** the system computes the hash locally
- **AND** returns hex-encoded hash

#### Scenario: Unsupported algorithm
- **WHEN** AI agent calls `crypto.hash` with unsupported algorithm
- **THEN** the system returns an error listing supported algorithms

### Requirement: Key Listing
The system SHALL provide a `crypto.keys` tool operation to list available keys.

#### Scenario: List available keys
- **WHEN** AI agent calls `crypto.keys`
- **THEN** the system returns array of key metadata (id, name, type)
- **AND** does NOT return actual key material
