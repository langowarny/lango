## Why

The Settings TUI editor lacks Cron/Background/Workflow configuration menus, so users cannot enable or configure automation features through the UI. Additionally, when the Telegram bot receives a cron scheduling request, the agent attempts `exec` → `lango cron add` which fails because spawning a new lango process requires passphrase authentication. Built-in in-process tools (`cron_add`, `bg_submit`, `workflow_run`) exist but the agent defaults to exec because no guardrail prevents it.

## What Changes

- Add three new settings menu categories (Cron Scheduler, Background Tasks, Workflow Engine) to the TUI editor with corresponding forms and config mappings
- Add exec guardrail that blocks `lango cron/bg/workflow` commands via exec/exec_bg tools, returning guidance to use built-in automation tools instead
- Enhance the automation system prompt to explicitly instruct the agent not to use exec for lango automation commands

## Capabilities

### New Capabilities

### Modified Capabilities
- `cli-settings`: Add three new menu categories (cron, background, workflow) with form builders and handleMenuSelection cases
- `cli-tuicore`: Add form→config field mappings for cron/background/workflow settings in UpdateConfigFromForm
- `tool-exec`: Add blockLangoExec guardrail to exec and exec_bg handlers that detects and blocks lango automation CLI commands
- `automation-agent-tools`: Add exec-prohibition instruction to the automation prompt section

## Impact

- **UI**: `internal/cli/settings/menu.go`, `editor.go`, `forms_impl.go` — 3 new menu items, 3 handler cases, 3 form builders
- **State**: `internal/cli/tuicore/state_update.go` — 12 new field→config mappings
- **Tools**: `internal/app/tools.go` — `blockLangoExec()` helper, `buildExecTools` signature change (added `automationAvailable` map)
- **Wiring**: `internal/app/app.go` — builds and passes `automationAvailable` map to `buildTools`
- **Prompts**: `internal/app/wiring.go` — added exec-prohibition note to `buildAutomationPromptSection`
