## MODIFIED Requirements

### Requirement: Event handling
The system SHALL process incoming Slack events using the Events API. Message handler invocations SHALL NOT block the event loop, allowing concurrent processing of interactive events.

#### Scenario: App mention event
- **WHEN** a user mentions the bot in a channel
- **THEN** the event SHALL be forwarded to the agent

#### Scenario: Direct message event
- **WHEN** a user sends a DM to the bot
- **THEN** the message SHALL be processed by the agent

#### Scenario: Concurrent message and interactive event processing
- **WHEN** a message handler is blocking (e.g., waiting for tool approval)
- **THEN** the event loop SHALL remain free to process interactive events (button clicks, approval callbacks)

#### Scenario: Graceful shutdown with active handlers
- **WHEN** the channel is stopped while message handlers are running
- **THEN** the system SHALL wait for all active handlers to complete before shutting down
