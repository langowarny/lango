## MODIFIED Requirements

### Requirement: Secrets get requires approval
The secrets.get operation SHALL require user approval via ApprovalMiddleware.

#### Scenario: Approval required
- **WHEN** AI calls secrets.get
- **THEN** ApprovalMiddleware intercepts and requests approval

#### Scenario: Approval granted
- **WHEN** user approves
- **THEN** an opaque reference token SHALL be returned (not the plaintext value)

### Requirement: Secrets get returns reference token
The secrets.get operation SHALL return an opaque reference token instead of the plaintext secret value. The plaintext SHALL be stored in the RefStore and resolved at execution time by the exec tool.

#### Scenario: Get secret returns reference
- **WHEN** AI calls secrets_get with name "api-key"
- **AND** the secret exists
- **THEN** the response SHALL contain `value: "{{secret:api-key}}"`
- **AND** the response SHALL contain a `note` field explaining the token usage
- **AND** the plaintext SHALL NOT appear in the response

#### Scenario: Get secret registers with scanner
- **WHEN** AI calls secrets_get with name "api-key"
- **AND** the secret exists
- **AND** a SecretScanner is configured
- **THEN** the plaintext value SHALL be registered with the SecretScanner for output scanning

### Requirement: Crypto decrypt returns reference token
The crypto.decrypt operation SHALL return an opaque reference token instead of the plaintext decrypted data. The plaintext SHALL be stored in the RefStore and resolved at execution time.

#### Scenario: Decrypt returns reference
- **WHEN** AI calls crypto_decrypt with valid ciphertext
- **THEN** the response SHALL contain `data: "{{decrypt:<uuid>}}"` where uuid is a unique identifier
- **AND** the response SHALL contain a `note` field explaining the token usage
- **AND** the plaintext SHALL NOT appear in the response

#### Scenario: Decrypt registers with scanner
- **WHEN** AI calls crypto_decrypt with valid ciphertext
- **AND** a SecretScanner is configured
- **THEN** the plaintext value SHALL be registered with the SecretScanner for output scanning
