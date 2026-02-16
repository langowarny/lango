## Why

Telegram approval button clicks trigger errors, and analysis reveals all channels (Telegram, Slack, Discord, Gateway) share common approval handler bugs: buttons remaining after timeout, TOCTOU race conditions on duplicate clicks, and hardcoded timeouts. These issues degrade user experience and can cause approval results to be silently lost.

## What Changes

- **Telegram**: Add `approvalPending` struct to store message metadata, implement `editApprovalMessage` helper with error classification, edit messages on timeout/cancel to show "Expired" and remove buttons, use `LoadAndDelete`-first pattern to prevent duplicate callbacks
- **Slack**: Fix TOCTOU race by using single `LoadAndDelete` call instead of `Load` then `LoadAndDelete`, add `editExpiredMessage` helper for timeout/cancel, pass empty `MsgOptionBlocks()` to remove action buttons
- **Discord**: Add `approvalPending` struct to capture `messageID`, implement `editExpiredMessage` for timeout/cancel using `ChannelMessageEditComplex`, remove buttons on expiry
- **Gateway**: Fix `handleApprovalResponse` to atomically delete pending entry within the lock to prevent duplicate response delivery, replace hardcoded 30s timeout with configurable `ApprovalTimeout` from `Config`

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `channel-telegram`: Add timeout message editing, keyboard removal on all exit paths, and TOCTOU-safe callback handling
- `channel-slack`: Fix TOCTOU race in HandleInteractive, add timeout message editing with button removal
- `channel-discord`: Add timeout message editing with button removal, capture message ID for edit operations
- `gateway-server`: Atomic delete in approval response handler, configurable approval timeout

## Impact

- `internal/channels/telegram/approval.go` — struct change, new helpers, refactored flow
- `internal/channels/slack/approval.go` — TOCTOU fix, new helper, timeout editing
- `internal/channels/discord/approval.go` — struct change, new helper, message ID capture
- `internal/gateway/server.go` — Config field addition, atomic delete fix
- All corresponding `_test.go` files updated with new test cases
