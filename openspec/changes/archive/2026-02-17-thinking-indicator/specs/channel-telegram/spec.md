## ADDED Requirements

### Requirement: Typing indicator during processing
The Telegram channel SHALL show a typing action indicator while the message handler processes a user message. The typing action SHALL be sent via `Request(ChatActionConfig)` immediately and refreshed every 4 seconds until the handler returns.

#### Scenario: Typing indicator during processing
- **WHEN** a user sends a message to the Telegram bot
- **THEN** the bot SHALL send a `ChatTyping` action to the chat before calling the handler
- **AND** SHALL refresh the typing action every 4 seconds
- **AND** SHALL stop refreshing when the handler returns

#### Scenario: Typing indicator API failure
- **WHEN** the `Request` call for typing action fails
- **THEN** the bot SHALL log a warning and continue processing normally
