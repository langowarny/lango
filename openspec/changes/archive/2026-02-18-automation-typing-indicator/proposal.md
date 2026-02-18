## Why

Automation systems (cron, background) execute agent tasks without providing visual feedback to the target channel. Interactive sessions already show typing indicators (Telegram 4s refresh, Discord 8s refresh, Slack placeholder), but automation runs only send a text notification before `runner.Run()` and remain silent until completion. Users have no way to know the agent is actively processing.

## What Changes

- Expose public `StartTyping` methods on Telegram, Discord, and Slack channel adapters for external callers (automation systems)
- Add `StartTyping` method to `channelSender` to dispatch typing indicators by channel target
- Define `TypingIndicator` interface in both `cron` and `background` packages (following existing AgentRunner/ChannelSender pattern)
- Wire typing indicators into `Delivery` (cron) and `Notification` (background) structs
- Call `StartTyping` before `runner.Run()` and `stopTyping()` after completion in both executor flows

## Capabilities

### New Capabilities

- `automation-typing-indicator`: Typing indicator integration for automation systems during agent execution

### Modified Capabilities

- `cron-scheduling`: Delivery gains typing indicator support during job execution
- `background-execution`: Notification gains typing indicator support during task execution
- `channel-telegram`: Public StartTyping method with context cancellation and sync.Once safety
- `channel-discord`: Public StartTyping method with context cancellation and sync.Once safety
- `channel-slack`: Public StartTyping method using placeholder message pattern with DeleteMessage

## Impact

- **Channel adapters**: `telegram.go`, `discord.go`, `slack.go` gain public `StartTyping` methods; Slack `Client` interface gains `DeleteMessage` method
- **App layer**: `sender.go` gains `StartTyping` dispatch method
- **Cron**: `delivery.go` gains `TypingIndicator` interface and typing field; `executor.go` calls typing around `runner.Run()`
- **Background**: `notification.go` gains `TypingIndicator` interface and typing field; `manager.go` calls typing around `runner.Run()`
- **Wiring**: `wiring.go` passes `channelSender` as typing indicator to both systems
- **Tests**: Slack mock clients need `DeleteMessage` method added
