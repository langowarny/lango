## Why

The `Executor` creates `~/.lango/skills/` during initialization, but this directory is only used as a temporary file location for script skill execution. Skill storage is handled by SQLite DB via `knowledge.Store`. This conflation of concerns makes the purpose unclear and introduces unnecessary filesystem side-effects during initialization.

## What Changes

- Remove `skillsDir` field from `Executor` struct and eliminate home directory lookup / `MkdirAll` from `NewExecutor`
- Switch script execution temp files from `filepath.Join(skillsDir, ...)` to `os.CreateTemp("", "lango-skill-*.sh")`
- Simplify `NewExecutor` and `NewRegistry` return types by removing `error` (no more fallible initialization)
- Update callers in `wiring.go` and test helpers to match simplified signatures

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `skill-system`: Executor no longer creates or depends on `~/.lango/skills/` directory; uses OS temp directory for script execution temp files

## Impact

- `internal/skill/executor.go` — struct field removal, constructor simplification, `executeScript` rewrite
- `internal/skill/registry.go` — `NewRegistry` return type change (no error)
- `internal/app/wiring.go` — caller update (remove error handling)
- `internal/skill/executor_test.go`, `registry_test.go` — remove `HOME` env override and error checks
