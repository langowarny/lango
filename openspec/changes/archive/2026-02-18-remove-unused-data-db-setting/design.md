## Context

The bootstrap process opens `~/.lango/lango.db` and passes its Ent client to `initSessionStore()`. The session store reuses this client, so the `session.databasePath` config (defaulting to `~/.lango/data.db`) is never opened in production. Standalone CLI commands (doctor, memory list) do use `session.databasePath` as a fallback, but the default points to a file that doesn't exist.

## Goals / Non-Goals

**Goals:**
- Remove the dead `db_path` field from Settings TUI to avoid user confusion.
- Align the default `session.databasePath` to `~/.lango/lango.db` so standalone CLI commands open the correct database.
- Clarify documentation in config types.

**Non-Goals:**
- Splitting session data into a separate database file.
- Removing the `DatabasePath` field from `SessionConfig` entirely (still needed for standalone CLI and tests).

## Decisions

1. **Remove UI field, keep config field**: The `DatabasePath` struct field remains in `SessionConfig` because standalone CLI commands and tests use it. Only the Settings TUI field is removed since users cannot meaningfully configure it.
2. **Default to `lango.db`**: Changing the default ensures standalone CLI paths (doctor checks, memory list) reference the actual database without requiring explicit user configuration.

## Risks / Trade-offs

- [Risk] Users with custom `session.databasePath` in `config.json` → No impact; their override still works. The change only affects the default value.
- [Risk] Existing `config.json` files still reference `data.db` → Updated `config.json` in repo; user-local configs are not auto-migrated but the default fallback is correct.
