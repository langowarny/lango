## Context

When users send messages to the agent via any channel (Telegram, Discord, Slack, Gateway WebSocket), the agent may take seconds to minutes to process. During this time, users see no feedback and may assume the message was lost. Each platform provides native mechanisms to indicate activity.

## Goals / Non-Goals

**Goals:**
- Show immediate visual feedback on every channel when the agent starts processing
- Use each platform's native mechanism (typing indicators or placeholder messages)
- Scope Gateway events to the authenticated user's session
- Gracefully handle indicator failures without affecting message delivery

**Non-Goals:**
- Progress updates during processing (e.g., "Searching knowledge base...")
- Streaming/partial responses
- Custom animated indicators beyond platform defaults

## Decisions

### 1. Per-platform native approach (vs. unified abstraction)

Each channel uses its platform-native mechanism rather than a shared abstraction layer:
- **Telegram**: `sendChatAction("typing")` — auto-expires after 5s, refresh at 4s
- **Discord**: `ChannelTyping()` — auto-expires after 10s, refresh at 8s
- **Slack**: Post placeholder message then update — Slack bots cannot trigger typing indicators
- **Gateway**: WebSocket events (`agent.thinking` / `agent.done`)

**Rationale**: The mechanisms are fundamentally different (ticker-based vs message-replace vs events). A shared abstraction would add complexity without benefit.

### 2. Goroutine + ticker for auto-refresh (Telegram/Discord)

A `startTyping()` method launches a goroutine with `time.Ticker` and returns a `func()` to stop it. The stop function closes a `done` channel to terminate the goroutine.

**Rationale**: Clean lifecycle — the caller just defers `stopThinking()`. The goroutine is guaranteed to stop when the handler returns.

### 3. Session-scoped Gateway broadcast

`BroadcastToSession` sends events only to UI clients matching the session key. When auth is disabled (empty session keys), it broadcasts to all UI clients.

**Rationale**: Prevents leaking thinking state to other users in multi-user setups.

### 4. Slack placeholder fallback

If `postThinking` fails, the flow falls back to the existing `Send` path. If `updateThinking` fails, a new message is sent instead.

**Rationale**: The thinking indicator is best-effort UX — it must never prevent response delivery.

## Risks / Trade-offs

- **[API rate limits]** Rapid typing indicator refreshes could hit rate limits under heavy load → Mitigation: conservative refresh intervals (4s/8s) well within platform limits
- **[Slack double-message on update failure]** If UpdateMessage fails, both placeholder and new response are visible → Mitigation: acceptable UX tradeoff; user sees the response regardless
- **[Goroutine leak if stopThinking not called]** A panic in the handler could skip the stop call → Mitigation: handler panics are recovered at a higher level; the goroutine ticker also stops if the channel is garbage collected
