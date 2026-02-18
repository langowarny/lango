## ADDED Requirements

### Requirement: Typing indicator during processing
The Discord channel SHALL show a typing indicator while the message handler processes a user message. The `Session` interface SHALL include `ChannelTyping`. The typing indicator SHALL be sent via `ChannelTyping` immediately and refreshed every 8 seconds until the handler returns.

#### Scenario: Typing indicator during processing
- **WHEN** a user sends a message (DM or mention) to the Discord bot
- **THEN** the bot SHALL call `ChannelTyping` on the channel before calling the handler
- **AND** SHALL refresh the typing indicator every 8 seconds
- **AND** SHALL stop refreshing when the handler returns

#### Scenario: Typing indicator API failure
- **WHEN** the `ChannelTyping` call fails
- **THEN** the bot SHALL log a warning and continue processing normally
