## Why

When a session TTL expires in channel adapters (Telegram/Slack/Discord), users receive a repeated `session expired: <key>` error with no recovery path. The root cause is that `EntStore.Get()` returns a plain string error on expiry, but `SessionServiceAdapter.Get()` only matches `ErrSessionNotFound` for auto-create — expired errors pass through unhandled, blocking the user permanently.

## What Changes

- Add `ErrSessionExpired` sentinel error to the session package for programmatic matching
- Wrap the TTL expiry error in `EntStore.Get()` with `ErrSessionExpired` using `%w`
- Add an expired-session branch in `SessionServiceAdapter.Get()` that deletes the stale session and auto-creates a fresh one
- Strengthen tests: sentinel error matching in TTL tests, mock store expiry simulation, auto-renew integration tests

## Capabilities

### New Capabilities

### Modified Capabilities
- `sentinel-errors`: Add `ErrSessionExpired` sentinel for session TTL expiry
- `session-auto-create`: Extend auto-create logic to handle expired sessions via delete-and-recreate

## Impact

- `internal/session/errors.go` — new sentinel
- `internal/session/ent_store.go` — TTL error wrapping
- `internal/adk/session_service.go` — expired branch in Get()
- `internal/session/store_test.go` — TTL test strengthening
- `internal/adk/state_test.go` — mockStore expiry support
- `internal/adk/session_service_test.go` — auto-renew tests
