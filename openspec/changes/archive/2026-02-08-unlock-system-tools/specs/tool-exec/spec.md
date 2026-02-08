## MODIFIED Requirements

### Requirement: Enhanced execution feedback
The system SHALL provide more descriptive feedback when commands fail or time out.

#### Scenario: Detailed failure message
- **WHEN** a command fails with a non-zero exit code
- **THEN** the system SHALL return both stdout and stderr to the agent for debugging
