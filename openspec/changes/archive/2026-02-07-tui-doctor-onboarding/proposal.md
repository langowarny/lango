## Why

Lango currently lacks user-friendly CLI tools for onboarding and diagnostics. Unlike OpenClaw which provides `openclaw onboard` and `openclaw doctor` commands, new users must manually create configuration files and troubleshoot issues without guidance. This creates a high barrier to entry and poor developer experience.

## What Changes

- Add `lango doctor` command with TUI interface for diagnosing configuration and connectivity issues
- Add `lango onboard` command with TUI wizard for initial setup (API key, model, channels)
- Use bubbletea library for rich interactive terminal UI
- Support `--fix` flag for auto-repair in doctor command
- Support `--json` flag for machine-readable output
- Implement minimal onboarding flow: API key → Model → Single channel

## Capabilities

### New Capabilities

- `cli-doctor`: Diagnostic command that checks configuration validity, API key verification, channel token validation, database accessibility, and port availability. Supports `--fix` for auto-repair and `--json` for scripted usage.
- `cli-onboard`: Interactive TUI wizard using bubbletea for first-time setup. Guides users through API key configuration, model selection, and channel setup with minimal steps.

### Modified Capabilities

- `config-system`: Add validation helpers for doctor command to verify configuration completeness

## Impact

- **New files**: `internal/cli/doctor/`, `internal/cli/onboard/`, `internal/cli/tui/`
- **Modified**: `cmd/lango/main.go` - add doctor and onboard subcommands
- **Dependencies**: Add `github.com/charmbracelet/bubbletea` and `github.com/charmbracelet/lipgloss`
- **APIs affected**: None (internal CLI only)
