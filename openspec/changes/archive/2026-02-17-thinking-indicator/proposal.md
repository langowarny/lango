## Why

When the agent processes a user message (which can take seconds to minutes), there is no visual feedback in the chat channel. Users may think their question was lost and resend it, or feel frustrated by the silent wait. Each channel platform provides native mechanisms to show a "thinking" or "typing" state that should be leveraged.

## What Changes

- Add typing indicator to **Telegram** using `sendChatAction("typing")` refreshed every 4 seconds
- Add typing indicator to **Discord** using `ChannelTyping()` refreshed every 8 seconds
- Add placeholder message to **Slack** ("Thinking...") that gets replaced with the actual response, since Slack bots cannot use the typing API
- Add `agent.thinking` / `agent.done` WebSocket events to the **Gateway** server, scoped to the user's session
- Add `BroadcastToSession` method to Gateway for session-scoped event delivery

## Capabilities

### New Capabilities
- `thinking-indicator`: Typing/thinking feedback across all channel adapters during agent processing

### Modified Capabilities
- `channel-telegram`: Add typing action indicator while handler processes messages
- `channel-discord`: Add typing indicator while handler processes messages
- `channel-slack`: Add thinking placeholder message flow (post → update)
- `gateway-server`: Add session-scoped broadcast and agent thinking/done events

## Impact

- **Code**: `internal/channels/telegram/telegram.go`, `internal/channels/discord/discord.go`, `internal/channels/discord/session.go`, `internal/channels/slack/slack.go`, `internal/gateway/server.go`
- **Interfaces**: Discord `Session` interface gains `ChannelTyping` method
- **Tests**: Updated mocks and new tests for all four adapters
- **Dependencies**: None (uses existing platform APIs)
