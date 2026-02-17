## Context

The `Executor` creates `~/.lango/skills/` at initialization time using `os.UserHomeDir()` + `os.MkdirAll()`. This directory serves only as a location for temporary script files during `executeScript`. Skill persistence is fully handled by SQLite via `knowledge.Store`. The directory creation makes `NewExecutor` fallible, which propagates error handling through `NewRegistry` and `initKnowledge` in `wiring.go`.

## Goals / Non-Goals

**Goals:**
- Remove the `~/.lango/skills/` directory dependency from executor initialization
- Use `os.CreateTemp` for script execution temp files (standard Go pattern)
- Simplify `NewExecutor` and `NewRegistry` to infallible constructors

**Non-Goals:**
- Changing skill storage (remains in SQLite)
- Changing script execution behavior or security validation
- Migrating or cleaning up existing `~/.lango/skills/` directories on user machines

## Decisions

**Use `os.CreateTemp` instead of a fixed directory**
- OS temp directory (`/tmp` or equivalent) is the standard location for short-lived process files
- `os.CreateTemp` provides unique filenames without race conditions
- No initialization side-effects — temp files are created on-demand during execution
- Alternative considered: keeping the directory but making creation lazy. Rejected because OS temp is simpler and more idiomatic.

**Remove error from constructor return types**
- With no filesystem operations at init time, `NewExecutor` cannot fail
- This simplifies `NewRegistry` (which wraps it) and `initKnowledge` in `wiring.go`
- Callers no longer need to handle initialization errors for the skill subsystem

## Risks / Trade-offs

**Temp file location is less predictable** → Acceptable trade-off. `defer os.Remove(f.Name())` ensures cleanup. The file exists only for the duration of script execution.

**Existing `~/.lango/skills/` directories remain on disk** → No migration needed. The directory was never documented or user-facing. Users can delete it manually if desired.
