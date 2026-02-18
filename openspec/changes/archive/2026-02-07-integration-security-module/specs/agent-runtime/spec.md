## ADDED Requirements

### Requirement: Security Tools Registration
The system SHALL register secrets and crypto tools during agent initialization.

#### Scenario: Tools registered on startup
- **WHEN** agent runtime is initialized
- **THEN** the system SHALL register `secrets` tool with store, get, list, delete operations
- **AND** SHALL register `crypto` tool with encrypt, decrypt, sign, hash, keys operations

### Requirement: Sensitive Tool Configuration
The system SHALL configure ApprovalMiddleware with security-sensitive tools.

#### Scenario: secrets.get requires approval
- **WHEN** agent runtime configures ApprovalMiddleware
- **THEN** `secrets.get` SHALL be included in SensitiveTools list

#### Scenario: crypto operations without approval
- **WHEN** AI calls crypto.encrypt or crypto.decrypt
- **THEN** the operation SHALL proceed without user approval
