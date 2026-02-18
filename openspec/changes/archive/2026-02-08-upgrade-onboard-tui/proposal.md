# Upgrade Onboard TUI

## Goal
Upgrade the `lango onboard` command from a simple setup wizard to a comprehensive Terminal User Interface (TUI) configuration editor. The new logical "Advanced" mode will allow users to configure all aspects of `lango.json` interactively, closing the gap between the CLI wizard and manual file editing.

## Background
The current "Advanced" mode in the onboarding wizard is a placeholder that behaves identically to "QuickStart". Users currently have to manually edit `lango.json` to configure essential features like:
- Server settings (Host, Port, HTTP/WS toggles)
- Agent parameters (Temperature, MaxTokens, System Prompt)
- Tool configuration (Exec timeouts, Browser headless mode)
- Security settings (PII redaction, Approval workflows)

This upgrade aims to provide a unified, interactive interface for all these settings.

## Proposed Changes

### UI Overhaul
- **Category-based Navigation**: Move from a linear step-by-step wizard to a navigable menu of configuration categories (Agent, Server, Tools, Channels, Security).
- **Interactive Forms**: Use comprehensive forms for each category, supporting text inputs, toggles (checkboxes), and lists.
- **Review & Save**: A final summary screen before saving the configuration.

### Enhanced Configuration Coverage
The TUI will support editing of:
1.  **Agent**: Provider, Model, MaxTokens, Temperature, SystemPromptPath.
2.  **Server**: Host, Port, HTTPEnabled, WSEnabled.
3.  **Channels**: Enable/Disable multiple channels (Telegram, Discord, Slack) and set their tokens.
4.  **Tools**: Configure Exec, Filesystem, and Browser tool settings.
5.  **Security**: Configure PII redaction, Approval steps, and Passphrase setup.

### Validation
- **Real-time Validation**: Validate inputs such as port numbers (1-65535), file paths, and URL formats.

## Capabilities

### Modified Capabilities
- `cli-onboard`: The interactive setup wizard will be significantly expanded to support full configuration editing.

## Impact
- **Codebase**: Major refactor of `internal/cli/onboard` package.
- **Dependencies**: minimal/no new dependencies (using existing `bubbletea`/`lipgloss`).
- **User Experience**: Significantly improved initial setup experience and discoverability of features.
