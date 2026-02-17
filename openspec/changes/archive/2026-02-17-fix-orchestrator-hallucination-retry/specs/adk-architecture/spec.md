## ADDED Requirements

### Requirement: Agent hallucination retry in RunAndCollect
`RunAndCollect` SHALL detect "failed to find agent" errors, extract the hallucinated agent name, send a correction message with valid sub-agent names, and retry once. If the retry also fails, the original error SHALL be returned.

#### Scenario: Hallucinated agent name triggers retry
- **WHEN** a `RunAndCollect` call yields an error matching `"failed to find agent: <name>"`
- **AND** the agent has sub-agents registered
- **THEN** the system SHALL send a correction message: `[System: Agent "<name>" does not exist. Valid agents: <list>. Please retry using one of the valid agent names listed above.]`
- **AND** retry the run exactly once with the correction message

#### Scenario: Retry succeeds
- **WHEN** the correction message retry produces a successful response
- **THEN** `RunAndCollect` SHALL return the successful response with no error

#### Scenario: Retry also fails
- **WHEN** the correction message retry also produces an error
- **THEN** `RunAndCollect` SHALL return the retry error

#### Scenario: Non-hallucination error is not retried
- **WHEN** `RunAndCollect` encounters an error that does not match "failed to find agent"
- **THEN** the error SHALL be returned immediately without retry

#### Scenario: No sub-agents means no retry
- **WHEN** `RunAndCollect` encounters a "failed to find agent" error
- **AND** the agent has no sub-agents
- **THEN** the error SHALL be returned immediately without retry
