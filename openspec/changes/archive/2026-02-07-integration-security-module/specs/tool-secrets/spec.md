## ADDED Requirements

### Requirement: Secrets Store Operation
The system SHALL provide a `secrets.store` tool operation that encrypts and stores a secret value.

#### Scenario: Store new secret
- **WHEN** AI agent calls `secrets.store` with name "api-key" and value "sk-12345"
- **THEN** the system encrypts the value using the default encryption key
- **AND** stores the encrypted value in the Secret entity
- **AND** returns success confirmation

#### Scenario: Update existing secret
- **WHEN** AI agent calls `secrets.store` with an existing secret name
- **THEN** the system encrypts the new value
- **AND** updates the existing Secret entity
- **AND** increments the access count

### Requirement: Secrets Get Operation
The system SHALL provide a `secrets.get` tool operation that retrieves and decrypts a secret value.

#### Scenario: Get existing secret with approval
- **WHEN** AI agent calls `secrets.get` with name "api-key"
- **AND** user approves the operation
- **THEN** the system retrieves the encrypted value from Secret entity
- **AND** decrypts using the associated key
- **AND** returns the plaintext value

#### Scenario: Get secret denied
- **WHEN** AI agent calls `secrets.get` with name "api-key"
- **AND** user denies the operation
- **THEN** the system returns an error "execution denied by user"

#### Scenario: Get non-existent secret
- **WHEN** AI agent calls `secrets.get` with name "unknown-key"
- **THEN** the system returns an error "secret not found"

### Requirement: Secrets List Operation
The system SHALL provide a `secrets.list` tool operation that returns stored secret names.

#### Scenario: List secrets
- **WHEN** AI agent calls `secrets.list`
- **THEN** the system returns array of secret names (not values)
- **AND** includes metadata (createdAt, accessCount) for each

### Requirement: Secrets Delete Operation
The system SHALL provide a `secrets.delete` tool operation that removes a secret.

#### Scenario: Delete existing secret
- **WHEN** AI agent calls `secrets.delete` with name "api-key"
- **THEN** the system removes the Secret entity
- **AND** returns success confirmation

#### Scenario: Delete non-existent secret
- **WHEN** AI agent calls `secrets.delete` with name "unknown-key"
- **THEN** the system returns an error "secret not found"

### Requirement: Approval Required for Get
The system SHALL require user approval before returning secret values.

#### Scenario: Approval middleware intercepts get
- **WHEN** `secrets.get` is registered as a tool
- **THEN** it SHALL be included in ApprovalMiddleware's SensitiveTools list
