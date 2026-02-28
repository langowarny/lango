## Context

The agent's `exec` tool spawns shell commands as subprocesses. Every `lango` CLI command calls `bootstrap.Run()` which invokes `passphrase.Acquire()`. In subprocess contexts, stdin is not a TTY, so passphrase acquisition hangs or fails unless a keyring/keyfile is pre-configured. The existing `blockLangoExec()` guard only covered 3 automation subcommands (cron, bg, workflow), leaving the agent free to attempt `lango security`, `lango graph`, `lango p2p`, etc. — all of which fail identically.

## Goals / Non-Goals

**Goals:**
- Block ALL `lango` CLI invocations from the agent's exec tool
- Provide actionable guidance: for subcommands with in-process equivalents, list the tools; for others, tell the agent to ask the user
- Maintain backward compatibility with the existing guard structure
- Cover both `exec` and `exec_bg` paths (both already call `blockLangoExec`)

**Non-Goals:**
- Adding new in-process tools for commands that lack them (config, doctor, settings, serve, onboard)
- Modifying the passphrase acquisition pipeline itself
- Changing the bootstrap flow to support non-interactive mode

## Decisions

1. **Two-phase guard structure** — Specific subcommand guards (Phase 1) take priority over the catch-all (Phase 2). This ensures precise tool alternative messages for known subcommands while still catching unknown/future ones.

2. **Feature flag awareness for automation guards only** — Automation guards (cron, bg, workflow) check `automationAvailable` to provide "enable the feature" guidance. Non-automation guards (graph, memory, p2p, security, payment) are always-on since these tools are always available when the feature itself is enabled.

3. **Prompt-level reinforcement** — In addition to the runtime guard, both `TOOL_USAGE.md` (system prompt) and the automation prompt section (dynamic) explicitly warn against `lango` CLI exec. Defense in depth: prompt prevents the attempt, guard catches it if attempted.

## Risks / Trade-offs

- **False positives for non-CLI lango references** — Commands like `echo lango` or `cat lango.yaml` are not blocked because the guard checks `lango ` (with space) or exact `lango`. Risk is minimal.
- **New subcommands not in Phase 1** — Future `lango` subcommands will be caught by the catch-all but won't get specific tool alternative messages. Acceptable: the catch-all message is still helpful.
