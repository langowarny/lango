## ADDED Requirements

### Requirement: P2P node key encrypted storage
The system SHALL store P2P Ed25519 node keys in `SecretsStore` (AES-256-GCM) under the key name `p2p.node.privatekey` instead of as plaintext files.

#### Scenario: New node key generation with SecretsStore available
- **WHEN** a P2P node starts for the first time and `SecretsStore` is available
- **THEN** the system SHALL generate an Ed25519 key, store it encrypted in SecretsStore under `p2p.node.privatekey`, and NOT create a plaintext `node.key` file

#### Scenario: Existing key loaded from SecretsStore
- **WHEN** a P2P node starts and SecretsStore contains `p2p.node.privatekey`
- **THEN** the system SHALL load and decrypt the key from SecretsStore without checking the filesystem

### Requirement: Legacy key auto-migration
The system SHALL automatically migrate plaintext `node.key` files to SecretsStore when both a legacy file exists and SecretsStore is available.

#### Scenario: Auto-migration of legacy node key
- **WHEN** a P2P node starts, SecretsStore is available, SecretsStore does NOT contain `p2p.node.privatekey`, and a plaintext `node.key` file exists
- **THEN** the system SHALL store the key in SecretsStore, delete the plaintext file, and log an info message confirming migration

#### Scenario: Migration failure is non-fatal
- **WHEN** migration to SecretsStore fails (e.g., DB locked)
- **THEN** the system SHALL log a warning, continue using the legacy file, and retry migration on next startup

### Requirement: Fallback to file-based storage
The system SHALL fall back to file-based key storage when `SecretsStore` is nil (not available).

#### Scenario: New key without SecretsStore
- **WHEN** a P2P node starts for the first time and `SecretsStore` is nil
- **THEN** the system SHALL generate an Ed25519 key and write it to `keyDir/node.key` with `0600` permissions

#### Scenario: Existing key loaded from file without SecretsStore
- **WHEN** a P2P node starts, `SecretsStore` is nil, and `keyDir/node.key` exists
- **THEN** the system SHALL load the key from the file

### Requirement: Key material memory cleanup
The system SHALL zero all key material byte slices from memory immediately after use using the `zeroBytes()` pattern.

#### Scenario: Key bytes zeroed after load
- **WHEN** node key bytes are loaded from SecretsStore or file
- **THEN** the raw byte slice SHALL be overwritten with zeros via `defer zeroBytes(data)` before the function returns
