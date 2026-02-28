## 1. Sentinel Error

- [x] 1.1 Add `ErrSessionExpired` to `internal/session/errors.go`

## 2. Error Wrapping

- [x] 2.1 Wrap TTL expiry in `EntStore.Get()` with `fmt.Errorf("get session %q: %w", key, ErrSessionExpired)`

## 3. Auto-Recovery Logic

- [x] 3.1 Add `ErrSessionExpired` branch in `SessionServiceAdapter.Get()` that deletes expired session and calls `getOrCreate()`

## 4. Tests

- [x] 4.1 Strengthen `TestEntStore_TTL` with `errors.Is(err, ErrSessionExpired)` assertion
- [x] 4.2 Add `TestEntStore_TTL_DeleteAndRecreate` for delete→recreate flow
- [x] 4.3 Add `expiredKeys` and `deleteErr` fields to `mockStore` in `state_test.go`
- [x] 4.4 Add `TestSessionServiceAdapter_Get_ExpiredSession_AutoRenews` test
- [x] 4.5 Add `TestSessionServiceAdapter_Get_ExpiredSession_DeleteFails` test

## 5. Verification

- [x] 5.1 Run `go build ./...` — no compilation errors
- [x] 5.2 Run `go test ./internal/session/... ./internal/adk/...` — all tests pass
