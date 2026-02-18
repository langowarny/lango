## Why

Cron jobs, background tasks, and workflows execute successfully but results are silently dropped when no delivery channel is configured. Users must explicitly specify `deliver_to`/`channel` every time, even when calling from a known channel like Telegram. This creates a poor UX where automation runs but users never see the output.

## What Changes

- Add `DefaultDeliverTo []string` config field to `CronConfig`, `BackgroundConfig`, and `WorkflowConfig`
- Auto-detect delivery channel from session key prefix (e.g., `telegram:123` → `telegram`) when not explicitly specified
- Fall back to config defaults when auto-detection is not possible (e.g., API/CLI calls)
- Upgrade executor warning logs from Debug to Warn when no delivery channel is configured
- Add Settings TUI form fields for configuring default delivery channels
- Update automation prompts to inform the agent that delivery channels are optional

## Capabilities

### New Capabilities
- `automation-delivery-fallback`: Automatic delivery channel detection and fallback chain for cron, background, and workflow systems

### Modified Capabilities
- `cron-scheduling`: cron_add handler now auto-detects and falls back to default delivery channels
- `background-execution`: bg_submit handler now auto-detects and falls back to default delivery channel
- `workflow-engine`: workflow_run handler now auto-injects delivery channels when YAML lacks deliver_to
- `config-system`: CronConfig, BackgroundConfig, WorkflowConfig gain DefaultDeliverTo field
- `cli-settings`: Cron, Background, Workflow forms gain "Default Deliver To" input field
- `automation-agent-tools`: Prompt sections updated to mention auto-detection

## Impact

- **Config**: `internal/config/types.go`, `internal/config/loader.go` — new fields and defaults
- **TUI**: `internal/cli/settings/forms_impl.go`, `internal/cli/tuicore/state_update.go` — new form fields
- **Tools**: `internal/app/tools.go` — new helper + modified build*Tools signatures and handlers
- **App**: `internal/app/app.go` — updated call sites for build*Tools
- **Executors**: `internal/cron/executor.go`, `internal/background/notification.go`, `internal/workflow/engine.go` — warning logs
- **Prompts**: `internal/app/wiring.go` — updated automation prompt sections
