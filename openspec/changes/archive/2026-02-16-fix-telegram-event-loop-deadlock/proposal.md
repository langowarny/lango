## Why

The Telegram channel's event loop processes both incoming messages and CallbackQuery updates on a single goroutine. When a message handler blocks waiting for tool approval (via `RequestApproval`), the event loop cannot process the resulting CallbackQuery from the user's button click. This creates a deadlock: the handler waits for a callback response that can never arrive because the event loop is blocked by the handler itself. The approval always times out regardless of whether the user clicks the button.

## What Changes

- Run `handleUpdate` in a separate goroutine so the event loop remains free to process CallbackQuery updates while the message handler is blocking on approval.
- Track the spawned goroutine with `sync.WaitGroup` to ensure graceful shutdown.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `channel-telegram`: Message handling SHALL be non-blocking to the event loop, allowing concurrent processing of CallbackQuery updates during handler execution.

## Impact

- **Code**: `internal/channels/telegram/telegram.go` — `Start` method event loop (single-line change: `go c.handleUpdate` with WaitGroup tracking).
- **Concurrency**: Message handlers now run concurrently. The existing handler is already stateless per-request, so this is safe.
- **Shutdown**: `Stop()` already waits on `c.wg`, so spawned goroutines are tracked correctly.
