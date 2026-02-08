## Requirement: System Prompt Support
The system SHALL support a configurable system prompt to guide agent behavior.

#### Scenario: Prepend system prompt to new session
- **WHEN** a new agent session is started
- **THEN** the system SHALL prepend a message with `role: system` containing the configured prompt to the conversation history

#### Scenario: Default identity prompt
- **WHEN** no custom system prompt is provided
- **THEN** a default prompt describing the agent's identity and tools SHALL be used
