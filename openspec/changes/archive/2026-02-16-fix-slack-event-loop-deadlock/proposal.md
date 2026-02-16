## Why

The Slack channel's single-goroutine event loop blocks when a message handler waits for tool approval, preventing interactive events (button clicks) from being processed. This causes a deadlock: the handler waits for approval, but the approval button click event can never be delivered because the event loop is stuck in the handler call.

## What Changes

- Move the handler invocation in `handleMessage` to a separate goroutine, so the event loop remains free to process interactive events (approval button clicks) concurrently.
- Track the spawned goroutine with the existing `sync.WaitGroup` for graceful shutdown.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `channel-slack`: Event handling requirement updated — message handler calls SHALL NOT block the event loop, enabling concurrent processing of interactive events (approval callbacks).

## Impact

- **Code**: `internal/channels/slack/slack.go` — `handleMessage` method
- **Behavior**: Message handling becomes concurrent; multiple messages can be processed in parallel while the event loop continues to dispatch interactive events.
- **Risk**: Minimal — follows the same pattern already proven in the Telegram channel (`internal/channels/telegram/telegram.go:152-156`).
