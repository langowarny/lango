## Context

The Settings TUI editor (`internal/cli/settings/`) provides a menu-driven interface for configuring all Lango features. While cron, background, and workflow systems were implemented with full config structs and runtime initialization, the Settings editor had no menu entries for them. Users had to edit YAML config files manually to enable these features.

Separately, when agents attempt to run `lango cron add` via exec, the spawned process requires passphrase authentication and fails. In-process tools (cron_add, bg_submit, workflow_run) already exist and work correctly.

## Goals / Non-Goals

**Goals:**
- Allow users to configure automation features (cron, background, workflow) through the Settings TUI
- Prevent agents from using exec to invoke lango automation CLI commands
- Guide agents toward using built-in in-process tools via prompt and guardrail

**Non-Goals:**
- Changing the underlying config structure (CronConfig, BackgroundConfig, WorkflowConfig already exist)
- Adding new automation capabilities or modifying existing tool behavior
- Blocking ALL lango CLI invocations via exec (only automation subcommands)

## Decisions

### 1. Form builder pattern for settings forms
**Decision**: Follow existing pattern — one `NewXxxForm()` function per category in `forms_impl.go`, returning `*tuicore.FormModel`.

**Rationale**: Consistent with all existing forms (Agent, Server, Channels, etc.). Each form builder maps config struct fields to form fields with appropriate input types.

### 2. Guardrail returns guidance instead of error
**Decision**: `blockLangoExec()` returns a map with `blocked: true` and a guidance `message` instead of returning an error.

**Rationale**: Returning structured data lets the agent understand why the command was blocked and what to do instead. An error would be less informative and could trigger retry loops.

### 3. Guardrail checks automation availability
**Decision**: Pass `automationAvailable map[string]bool` to exec tools so the guardrail can give context-aware guidance — either "use the built-in tools" (if enabled) or "enable the feature in Settings" (if disabled).

**Rationale**: Without this context, the agent would get a generic error. With it, the agent can either use built-in tools directly or instruct the user to enable the feature.

### 4. Prompt reinforcement as defense-in-depth
**Decision**: Add explicit exec-prohibition text to `buildAutomationPromptSection()` in addition to the runtime guardrail.

**Rationale**: The prompt instruction prevents the agent from even attempting the exec call in most cases. The runtime guardrail catches edge cases where the prompt is ignored.

## Risks / Trade-offs

- **String matching fragility**: `blockLangoExec` uses `strings.HasPrefix` on lowercased commands. Unusual invocations (e.g., `/usr/local/bin/lango cron`) bypass the check. → Acceptable: the prompt instruction handles the common case, and direct binary paths are rare in agent-generated commands.
- **Form field count growth**: Each new category adds 3-5 form fields to the switch statement in `UpdateConfigFromForm`. → Manageable: the switch is already large and follows a clear pattern.
