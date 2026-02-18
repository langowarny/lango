## Context

The Telegram channel uses a single-goroutine event loop that reads from `bot.GetUpdatesChan()` and dispatches updates. Regular messages are handled by `c.handleUpdate(ctx, update)` which calls `c.handler(ctx, incoming)` synchronously. When the handler blocks (e.g., waiting for tool approval via `RequestApproval`), the event loop stalls and cannot process CallbackQuery updates. This creates a deadlock between the approval-waiting handler and the callback-dispatching event loop.

## Goals / Non-Goals

**Goals:**
- Keep the event loop non-blocking so CallbackQuery updates are processed while message handlers execute.
- Maintain graceful shutdown semantics — all spawned handler goroutines complete before `Stop()` returns.

**Non-Goals:**
- Adding concurrency limits or rate limiting for message handlers.
- Changing the approval provider's internal architecture.

## Decisions

**Decision: Spawn `handleUpdate` in a goroutine tracked by `sync.WaitGroup`**

```go
c.wg.Add(1)
go func() {
    defer c.wg.Done()
    c.handleUpdate(ctx, update)
}()
```

**Alternatives considered:**
1. **Dedicated callback goroutine**: Split callbacks and messages into two separate goroutines reading from the same channel. Rejected — `GetUpdatesChan` returns a single channel; splitting requires filtering, adding complexity without benefit.
2. **Buffered callback channel**: Route callbacks to a separate buffered channel processed by a second goroutine. Rejected — over-engineering for a problem solved by simply not blocking the event loop.

**Rationale**: The goroutine approach is minimal, idiomatic Go, and the existing `sync.WaitGroup` on `Channel` already supports tracking multiple goroutines. The handler is stateless per-request, so concurrent execution is safe.

## Risks / Trade-offs

- **[Low] Concurrent handler execution** → Handlers now run in parallel. Since each handler operates on an independent message context and the bot's `Send` method is thread-safe, this is safe. If ordering becomes a concern in the future, per-chat serialization can be added.
- **[Low] Goroutine leak on slow handlers** → Mitigated by the existing context cancellation propagated through `ctx`. On `Stop()`, context is cancelled and `wg.Wait()` ensures all handlers finish.
