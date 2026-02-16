## Context

The Slack channel uses Socket Mode with a single event loop goroutine (`handleEvents`) that processes all incoming events sequentially. When a message triggers a handler that blocks (e.g., waiting for tool approval via `RequestApproval()`), the event loop cannot process subsequent interactive events (button clicks). This creates a deadlock: the handler waits for an approval response that can never arrive.

The Telegram channel already solved this by spawning `handleUpdate` in a separate goroutine (commit 912cb5c). Discord is unaffected because `discordgo` runs each handler in its own goroutine.

## Goals / Non-Goals

**Goals:**
- Eliminate the event loop deadlock in the Slack channel
- Allow concurrent processing of message events and interactive events
- Maintain graceful shutdown behavior via `sync.WaitGroup`

**Non-Goals:**
- Concurrency limiting (rate limiting, worker pools) — not needed at current scale
- Refactoring the event loop architecture beyond the minimal fix
- Changes to Discord or Telegram channels

## Decisions

### Decision 1: Spawn goroutine in `handleMessage`

Move the `c.handler(ctx, incoming)` call and its response handling into a new goroutine within `handleMessage`, tracked by `c.wg`.

**Why here (not in `handleCallbackEvent` or `handleEventsAPI`)**:
- `handleMessage` is the exact point where blocking occurs
- Keeps event parsing/routing synchronous (easier to reason about)
- Matches the Telegram pattern (`handleUpdate` in goroutine)

**Alternatives considered**:
- Worker pool with bounded concurrency: Over-engineering for current usage; can be added later if needed.
- Spawning goroutine at the event loop level: Would make all event types concurrent, but interactive events are already lightweight and don't need it.

## Risks / Trade-offs

- **[Concurrent handler calls]** → Multiple messages may invoke the handler simultaneously. This is acceptable because the handler is already designed for concurrent use (the approval system uses per-session channels).
- **[Goroutine leak on shutdown]** → Mitigated by `c.wg.Add(1)` / `defer c.wg.Done()` and context cancellation propagation via `ctx`.
