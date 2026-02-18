## MODIFIED Requirements

### Requirement: Message reception
The system SHALL receive and process incoming messages from Telegram chats. Message handling SHALL be dispatched to a separate goroutine so that the event loop remains non-blocking and can continue processing CallbackQuery updates concurrently.

#### Scenario: Direct message received
- **WHEN** a user sends a direct message to the bot
- **THEN** the message SHALL be forwarded to the agent with sender context

#### Scenario: Group message with mention
- **WHEN** a user mentions the bot in a group chat
- **THEN** the message SHALL be processed if mention-gating allows

#### Scenario: Concurrent callback processing during handler execution
- **WHEN** a message handler is blocking (e.g., waiting for tool approval)
- **AND** a CallbackQuery update arrives from a button click
- **THEN** the event loop SHALL process the CallbackQuery immediately without waiting for the handler to complete

#### Scenario: Graceful shutdown with active handlers
- **WHEN** the channel is stopped while message handlers are still executing
- **THEN** the system SHALL wait for all active handler goroutines to complete before returning from Stop()
