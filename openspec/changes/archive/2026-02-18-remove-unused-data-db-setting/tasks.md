## 1. Settings TUI Cleanup

- [x] 1.1 Remove `db_path` field from `NewSessionForm()` in `internal/cli/settings/forms_impl.go`
- [x] 1.2 Remove `db_path` case from state update handler in `internal/cli/tuicore/state_update.go`
- [x] 1.3 Update `TestNewSessionForm_AllFields` to expect 2 fields (no `db_path`) in `internal/cli/settings/forms_impl_test.go`

## 2. Config Default Alignment

- [x] 2.1 Change default `session.databasePath` from `~/.lango/data.db` to `~/.lango/lango.db` in `internal/config/loader.go`
- [x] 2.2 Update `SessionConfig` struct comment in `internal/config/types.go` to document bootstrap reuse behavior
- [x] 2.3 Update `config.json` default `databasePath` to `~/.lango/lango.db`

## 3. Verification

- [x] 3.1 Run `go build ./...` — no compilation errors
- [x] 3.2 Run `go test ./internal/cli/settings/... ./internal/config/...` — all tests pass
