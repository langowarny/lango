# Fix Serve Encryption Tasks

- [x] Update `internal/session/store.go` (or equivalent) to accept passphrase options
- [x] Update `internal/app/app.go` to retrieve passphrase from config/env and pass to `session.NewEntStore`
- [x] Verify `lango serve` startup with `LANGO_PASSPHRASE` set
- [x] Verify `lango serve` startup fails gracefully if passphrase is missing for encrypted DB
