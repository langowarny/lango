## Purpose

Debug-level logging of multi-agent event streams for runtime observability of agent delegation and response routing in `RunAndCollect()`.

## Requirements

### Requirement: Agent delegation logging
The system SHALL log a debug-level message when an agent delegates to another agent, including the source agent name, target agent name, and session ID.

#### Scenario: Orchestrator delegates to sub-agent
- **WHEN** an event in `RunAndCollect()` has a non-empty `Author` and a non-empty `Actions.TransferToAgent`
- **THEN** the system logs a debug message with key `"agent delegation"`, fields `from` (Author), `to` (TransferToAgent), and `session` (session ID)

### Requirement: Agent response logging
The system SHALL log a debug-level message when an agent produces a text response, including the agent name and session ID.

#### Scenario: Agent produces text content
- **WHEN** an event in `RunAndCollect()` has a non-empty `Author`, no `TransferToAgent`, and contains at least one non-empty text part
- **THEN** the system logs a debug message with key `"agent response"`, fields `agent` (Author) and `session` (session ID)

#### Scenario: Agent event with no text content
- **WHEN** an event in `RunAndCollect()` has a non-empty `Author`, no `TransferToAgent`, and no text parts
- **THEN** the system SHALL NOT log an agent response message

### Requirement: Debug level only
All agent event logging SHALL use debug log level to ensure zero noise in production environments.

#### Scenario: Production log level
- **WHEN** the log level is set to `info` or above
- **THEN** no agent event log messages are emitted

#### Scenario: Debug log level enabled
- **WHEN** the log level is set to `debug`
- **THEN** agent delegation and response log messages are emitted

### Requirement: Logging subsystem
Agent event logs SHALL use the `agent` logging subsystem via `logging.Agent()`.

#### Scenario: Log subsystem identification
- **WHEN** an agent event log message is emitted
- **THEN** the log entry is tagged with the `agent` subsystem name
