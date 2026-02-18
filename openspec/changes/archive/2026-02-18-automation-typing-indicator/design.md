## Context

Interactive sessions already display typing indicators while the agent processes messages: Telegram sends `ChatTyping` actions every 4 seconds, Discord calls `ChannelTyping` every 8 seconds, and Slack posts a `_Thinking..._` placeholder that gets replaced. However, automation systems (cron jobs, background tasks) call `runner.Run()` with no visual feedback — users only see a start notification and then silence until the result arrives.

The existing `startTyping` methods on channel adapters are private and lack context-based cancellation, making them unsuitable for external callers.

## Goals / Non-Goals

**Goals:**
- Show typing indicators on target channels while automation systems execute agent tasks
- Reuse existing channel typing mechanisms (Telegram ChatTyping, Discord ChannelTyping, Slack placeholder)
- Ensure typing goroutines are automatically cleaned up via context cancellation and sync.Once safety
- Follow the existing interface-per-consumer pattern (TypingIndicator defined in each consumer package)

**Non-Goals:**
- Modifying the existing interactive session typing behavior
- Adding typing indicators to the workflow engine (can be done as a follow-up)
- Implementing real-time progress updates or streaming responses

## Decisions

### 1. Public StartTyping on channel adapters (not extracting to a shared service)

Each channel adapter gets a public `StartTyping` method alongside the existing private `startTyping`. This keeps the interactive path unchanged and avoids a new abstraction layer. The public method adds `context.Context` monitoring and `sync.Once` double-close protection.

**Alternative**: Create a shared `TypingService` that abstracts across channels. Rejected because it adds unnecessary indirection — the `channelSender` already dispatches by channel type.

### 2. Interface-per-consumer pattern for TypingIndicator

Both `cron` and `background` packages define their own `TypingIndicator` interface (same as existing `AgentRunner`, `ChannelSender` patterns). This avoids import cycles and keeps each package self-contained.

**Alternative**: Define a shared interface in a common package. Rejected because it contradicts the established pattern and would create a dependency direction issue.

### 3. channelSender satisfies both ChannelSender and TypingIndicator

The `channelSender` struct in `app/sender.go` implements both interfaces. In wiring, the same `sender` instance is passed as both arguments. This keeps wiring simple with no new types.

### 4. Slack uses DeleteMessage for placeholder cleanup

Slack's `StartTyping` posts a `_Processing..._` placeholder and returns a stop function that deletes it. This required adding `DeleteMessage` to the `Client` interface. The existing `postThinking` / `updateThinking` pattern for interactive sessions remains unchanged.

**Alternative**: No-op for Slack (start notification serves as feedback). Rejected because a deletable placeholder provides cleaner UX — the placeholder disappears when the actual result arrives.

## Risks / Trade-offs

- **[Typing failure blocks execution]** → Mitigated: all typing errors are logged at Warn level and return no-op stop functions. Typing never blocks the agent run.
- **[Goroutine leak on context cancel]** → Mitigated: public `StartTyping` methods monitor `ctx.Done()` in addition to the `done` channel. Goroutines always exit.
- **[Double-close panic]** → Mitigated: `sync.Once` wraps all stop functions. Callers can safely call stop multiple times.
- **[Slack DeleteMessage permission]** → If the bot lacks `chat:write` scope for deletion, the error is logged and the placeholder persists. Minimal impact.
