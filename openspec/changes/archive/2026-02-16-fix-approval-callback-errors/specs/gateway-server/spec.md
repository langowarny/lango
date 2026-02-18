## MODIFIED Requirements

### Requirement: Atomic delete in approval response handler
The Gateway server SHALL delete the pending approval entry within the same lock scope as the lookup to prevent duplicate response delivery.

#### Scenario: First response delivers result and deletes entry
- **WHEN** a companion sends an approval response for a pending request
- **THEN** the system SHALL atomically look up and delete the entry under the same mutex lock, then deliver the result

#### Scenario: Duplicate response has no effect
- **WHEN** a second approval response arrives for an already-deleted request
- **THEN** the system SHALL find no entry and skip delivery without error

### Requirement: Configurable approval timeout
The Gateway server SHALL use `Config.ApprovalTimeout` instead of a hardcoded 30-second timeout for approval requests. If the configured value is zero or negative, it SHALL default to 30 seconds.

#### Scenario: Custom timeout from config
- **WHEN** Config.ApprovalTimeout is set to 60 seconds
- **THEN** the system SHALL wait up to 60 seconds for an approval response before timing out

#### Scenario: Default timeout when not configured
- **WHEN** Config.ApprovalTimeout is zero
- **THEN** the system SHALL use the default 30-second timeout
