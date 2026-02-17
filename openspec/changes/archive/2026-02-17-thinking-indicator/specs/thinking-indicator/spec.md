## ADDED Requirements

### Requirement: Typing indicator during agent processing
All channel adapters SHALL show a typing or thinking indicator immediately when a user message is received, and SHALL stop the indicator when the agent response is ready.

#### Scenario: Indicator starts before handler
- **WHEN** a user sends a message to any channel adapter
- **THEN** the adapter SHALL activate the platform-native thinking indicator before invoking the message handler

#### Scenario: Indicator stops after handler completes
- **WHEN** the message handler returns (success or error)
- **THEN** the adapter SHALL stop the thinking indicator

#### Scenario: Indicator failure does not block response
- **WHEN** the thinking indicator API call fails
- **THEN** the adapter SHALL log a warning and continue processing the message normally
