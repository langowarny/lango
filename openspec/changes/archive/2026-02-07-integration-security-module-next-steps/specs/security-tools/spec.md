## ADDED Requirements

### Requirement: Secrets tool registration
The app SHALL register secrets tool with store/get/list/delete operations.

#### Scenario: Tool available
- **WHEN** agent runtime is initialized
- **THEN** secrets tool is available with operations: store, get, list, delete

### Requirement: Crypto tool registration
The app SHALL register crypto tool with encrypt/decrypt/sign/hash/keys operations.

#### Scenario: Tool available
- **WHEN** agent runtime is initialized
- **THEN** crypto tool is available with operations: encrypt, decrypt, sign, hash, keys

### Requirement: Secrets get requires approval
The secrets.get operation SHALL require user approval via ApprovalMiddleware.

#### Scenario: Approval required
- **WHEN** AI calls secrets.get
- **THEN** ApprovalMiddleware intercepts and requests approval

#### Scenario: Approval granted
- **WHEN** user approves
- **THEN** secret value is returned to AI

### Requirement: LocalCryptoProvider initialization
The app SHALL initialize LocalCryptoProvider with passphrase when security.signer.provider is "local".

#### Scenario: First-time setup
- **WHEN** no salt exists in database
- **THEN** prompt user for passphrase and store salt

#### Scenario: Subsequent startup
- **WHEN** salt exists in database
- **THEN** prompt for passphrase and derive key
