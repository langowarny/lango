## Context

Channel adapters (Telegram/Slack/Discord) generate session keys like `telegram:123:456`. When the session TTL expires, `EntStore.Get()` returns `fmt.Errorf("session expired: %s", key)` — a plain string error. `SessionServiceAdapter.Get()` only checks for `ErrSessionNotFound` to trigger auto-creation, so expired errors pass through unhandled. Users see a permanent `session expired` error with no recovery path.

## Goals / Non-Goals

**Goals:**
- Expired sessions auto-recover transparently — users never see expiry errors
- Reuse existing patterns: sentinel errors, `Delete()`, `getOrCreate()` — no new interfaces
- Maintain concurrent safety via existing `getOrCreate()` retry logic

**Non-Goals:**
- Session data migration or history preservation across expiry boundaries
- Changing TTL configuration or adding per-session TTL overrides
- Adding explicit "renew" or "extend" session operations

## Decisions

### Use sentinel error + errors.Is matching (not string matching)
Add `ErrSessionExpired` to the session error catalog and wrap it with `%w` in `EntStore.Get()`. This follows the established sentinel pattern (`ErrSessionNotFound`, `ErrDuplicateSession`) and enables type-safe matching via `errors.Is()`.

**Alternative rejected**: Adding an `IsExpired()` method or using string matching. Sentinel errors are the established project pattern and integrate cleanly with Go's `errors.Is()`.

### Delete-then-recreate (not in-place renewal)
On expiry, delete the stale session record first, then call `getOrCreate()` to create fresh. This avoids adding new Store interface methods and leverages the existing concurrent-safe `getOrCreate()` flow (handles `ErrDuplicateSession` retry).

**Alternative rejected**: Adding `Renew(key)` to the Store interface. This would break all Store implementations and mock stores for a feature that delete+create achieves identically.

### Log at Info level on auto-renewal
The expired→renewed transition is logged via `logger().Infow()` for operational visibility. This is not an error condition — it's expected lifecycle behavior.

## Risks / Trade-offs

- [Session history lost on expiry] → By design: TTL expiry means the session is stale, and preserving old history could confuse the AI context. Fresh start is the intended behavior.
- [Delete failure blocks recovery] → The delete error is propagated with context. This is correct: if the DB can't delete, it likely can't create either, and masking the error would cause silent failures.
- [Tiny race window between Delete and getOrCreate] → Another goroutine could create the session in this window. `getOrCreate()` handles this via `ErrDuplicateSession` retry, so it's safe.
