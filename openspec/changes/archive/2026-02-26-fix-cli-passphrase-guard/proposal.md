## Why

The agent's `exec` tool only blocked 3 `lango` CLI subcommands (cron, bg, workflow) from being invoked via subprocess. All other `lango` CLI commands (security, memory, graph, p2p, config, doctor, etc.) could still be attempted, but would fail silently because every `lango` command requires passphrase authentication via `bootstrap.Run()` — which hangs or errors in non-interactive subprocess contexts.

## What Changes

- Expand `blockLangoExec()` to block ALL `lango` CLI subcommands, not just 3
- Add a catch-all guard for any `lango` prefix that doesn't match specific subcommands
- Provide per-subcommand guidance messages pointing to in-process tool equivalents (graph, memory, p2p, security, payment)
- For subcommands without in-process equivalents (config, doctor, settings), instruct the agent to ask the user to run them directly
- Update `TOOL_USAGE.md` prompt to explicitly warn against using exec with any `lango` command
- Broaden the automation prompt section in `wiring.go` to cover all `lango` subcommands

## Capabilities

### New Capabilities

### Modified Capabilities
- `agent-tools`: Expanded CLI exec guard to block all lango subcommands with per-command tool alternatives

## Impact

- `internal/app/tools.go` — `blockLangoExec()` function rewritten with comprehensive guard list and catch-all
- `internal/app/tools_test.go` — New test cases for all subcommands and catch-all behavior
- `prompts/TOOL_USAGE.md` — Added exec safety rules at top of Exec Tool section
- `internal/app/wiring.go` — Broadened automation prompt warning to cover all lango CLI commands
