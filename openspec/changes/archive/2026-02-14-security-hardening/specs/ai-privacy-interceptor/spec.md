## MODIFIED Requirements

### Requirement: Approval Workflow
The system SHALL block execution of "sensitive" tools (configured list) and require explicit approval before proceeding. The system SHALL use fail-closed semantics: without explicit approval, execution is denied.

#### Scenario: Companion approval granted
- **WHEN** the agent attempts to call a sensitive tool
- **AND** a companion is connected
- **AND** the companion approves the request
- **THEN** the tool execution SHALL proceed

#### Scenario: Companion approval denied
- **WHEN** the agent attempts to call a sensitive tool
- **AND** a companion is connected
- **AND** the companion denies the request
- **THEN** the tool execution SHALL be denied with an error

#### Scenario: Companion approval error
- **WHEN** the agent attempts to call a sensitive tool
- **AND** a companion is connected
- **AND** the approval request encounters an error
- **THEN** the tool execution SHALL be denied (fail-closed)

#### Scenario: TTY fallback approval
- **WHEN** the agent attempts to call a sensitive tool
- **AND** no companion is connected
- **AND** stdin is a terminal (TTY)
- **THEN** the system SHALL prompt the user via stderr with "Allow? [y/N]"
- **AND** proceed only if user responds "y" or "yes"

#### Scenario: No approval source available
- **WHEN** the agent attempts to call a sensitive tool
- **AND** no companion is connected
- **AND** stdin is not a terminal
- **THEN** the tool execution SHALL be denied with error "no approval source available"
