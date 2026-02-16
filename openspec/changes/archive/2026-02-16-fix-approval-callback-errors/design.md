## Context

All four approval channels (Telegram, Slack, Discord, Gateway) have independent implementations of the approval request/response flow. Each stores pending approvals in a sync.Map or mutex-guarded map and uses channels to deliver results. Analysis revealed shared bugs: buttons persisting after timeout, TOCTOU races on duplicate clicks, and hardcoded timeouts.

## Goals / Non-Goals

**Goals:**
- Ensure approval buttons are removed on all exit paths (approve, deny, timeout, context cancellation)
- Eliminate TOCTOU race conditions where duplicate clicks can deliver conflicting results
- Make Gateway approval timeout configurable via Config struct
- Downgrade noisy expired-callback errors to Debug level in Telegram

**Non-Goals:**
- Redesigning the approval architecture or introducing a shared base type
- Adding persistent storage for approval state
- Changing the approval UI/UX beyond button removal and expired messages

## Decisions

### 1. `approvalPending` struct pattern across all channels
**Decision**: Store message metadata (chatID/messageID, channelID/timestamp) alongside the response channel in a struct, rather than bare `chan bool`.
**Rationale**: Timeout and cancellation paths need message coordinates to edit the message. Bare channels lose this context. Discord already removes buttons on interaction but not on timeout.

### 2. `LoadAndDelete`-first in callback handlers
**Decision**: Call `sync.Map.LoadAndDelete` as the first operation in all callback/interaction handlers, before any message editing.
**Rationale**: This is the atomic operation that prevents TOCTOU races. If two callbacks arrive simultaneously, only the first `LoadAndDelete` succeeds. The previous Slack pattern (`Load` → edit → `LoadAndDelete`) had a window where the timeout's `defer Delete` could execute between Load and LoadAndDelete, silently losing the approval result.
**Alternative**: Using a separate mutex per request — rejected as over-engineering given sync.Map already provides atomic operations.

### 3. Atomic delete in Gateway `handleApprovalResponse`
**Decision**: Move `delete(s.pendingApprovals, req.RequestID)` inside the same lock scope as the lookup.
**Rationale**: The previous pattern (Lock → lookup → Unlock → send) left the entry in the map, allowing a duplicate response to send a second value to the channel. Moving delete inside the lock ensures exactly-once delivery.

### 4. Config-driven approval timeout for Gateway
**Decision**: Add `ApprovalTimeout time.Duration` to `gateway.Config` with a 30s default fallback.
**Rationale**: Hardcoded timeouts prevent deployment-specific tuning. The other channels already accept timeout as a constructor parameter.

## Risks / Trade-offs

- **[Risk]** Editing expired Telegram messages may fail if the message was already deleted → Mitigation: `isMessageNotModifiedErr` filter suppresses benign errors
- **[Risk]** Gateway Config change is additive but callers not setting `ApprovalTimeout` get 0 → Mitigation: Explicit fallback to 30s when value is <= 0
- **[Trade-off]** Using string matching for Telegram error classification (`"query is too old"`, `"message is not modified"`) is fragile → Accepted: Telegram Bot API does not expose structured error codes
