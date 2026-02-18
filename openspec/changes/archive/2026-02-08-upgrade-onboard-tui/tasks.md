# Upgrade Onboard TUI Tasks

## Core TUI Framework
- [x] Refactor `Wizard` model to hold `config.Config` state instead of simplified `WizardConfig`
- [x] Implement `CategoryList` component for main menu navigation
- [x] Create base `Form` component structure and input handling

## Configuration Forms
- [x] Implement `AgentForm` (Provider, Model, MaxTokens, Temp)
- [x] Implement `ServerForm` (Host, Port, HTTP/WS toggles)
- [x] Implement `ChannelsForm` (Telegram, Discord, Slack toggles and tokens)
- [x] Implement `ToolsForm` (Exec timeouts, Browser settings)
- [x] Implement `SecurityForm` (PII, Approval, Passphrase)

## Integration & Logic
- [x] Wire up navigation between Menu and Forms
- [x] Implement real-time validation for integer fields (Port, MaxTokens)
- [x] Implement "Save & Exit" logic to write `lango.json`
- [x] Implement `.lango.env` template generation for secrets

## Verification
- [x] Manual verification: Launch `lango onboard`, navigate to all sections, modify values, save, and verify `lango.json` output.
