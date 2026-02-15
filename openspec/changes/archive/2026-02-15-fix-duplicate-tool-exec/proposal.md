## Why

When the knowledge system is enabled, base tools (exec, exec_bg, etc.) are added twice to the ADK runner's tool list—once as wrapped tools and once via `Registry.AllTools()` which includes the original base tools. This causes the ADK runner to reject the duplicate tool names with `"agent error: duplicate tool: \"exec\""`, breaking Telegram message handling in Docker.

## What Changes

- Add `LoadedSkills()` method to `skill.Registry` that returns only dynamically loaded skill tools (excluding base tools)
- Change `app.go` to call `LoadedSkills()` instead of `AllTools()` when appending skill tools, preventing base tool duplication

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `skill-system`: `Registry` gains a `LoadedSkills()` method to return only dynamic skills without base tools

## Impact

- `internal/skill/registry.go` — new `LoadedSkills()` method
- `internal/app/app.go` — call site change from `AllTools()` to `LoadedSkills()`
- Fixes runtime crash in Docker when knowledge system is active and messages are received via Telegram
