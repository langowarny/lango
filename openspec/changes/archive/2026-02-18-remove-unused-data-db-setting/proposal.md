## Why

The Settings TUI displays a "Database Path" field (`~/.lango/data.db`) under Session Configuration, but this value is never used in production. The bootstrap process always opens `~/.lango/lango.db` and the session store reuses that Ent client. The `data.db` file is never created, causing user confusion.

## What Changes

- Remove the `db_path` field from the Session Configuration form in the Settings TUI.
- Change the default `session.databasePath` from `~/.lango/data.db` to `~/.lango/lango.db` so standalone CLI commands (doctor, memory list) reference the correct file.
- Update config types documentation to clarify that `lango.db` is the primary database and `DatabasePath` is a fallback for standalone CLI access.
- Update `config.json` to reflect the corrected default.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `cli-settings`: Remove the dead `db_path` field from Session Configuration form.
- `config-system`: Change default `session.databasePath` to `~/.lango/lango.db` and clarify documentation.

## Impact

- `internal/cli/settings/forms_impl.go` — field removed
- `internal/cli/tuicore/state_update.go` — `db_path` case removed
- `internal/config/loader.go` — default value changed
- `internal/config/types.go` — comment updated
- `config.json` — default value updated
- `internal/cli/settings/forms_impl_test.go` — test expectation updated
