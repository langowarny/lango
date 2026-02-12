## Why

Lango's Telegram bot is currently unresponsive because the application blocks on an interactive passphrase prompt during startup, even when the `LANGO_PASSPHRASE` environment variable is provided. Furthermore, the startup logs are too sparse, leaving users uncertain about the system's operational state (whether it's stuck, initializing, or ready).

## What Changes

1.  **Non-blocking Passphrase Initialization**: Modify the security initialization in `internal/app/app.go` to always check for the `LANGO_PASSPHRASE` environment variable first, even in interactive terminal sessions. Interactive prompting will only occur if the environment variable is missing.
2.  **Granular Startup Logging**: Introduce detailed logs in `cmd/lango/main.go` and `internal/app/app.go` to track the initialization of core components (Supervisor, Agent, Tools, Gateway, and Channels).
3.  **Standardized Readiness Feedback**: Implement a standardized "Ready" logging pattern that clearly indicates when the Gateway is listening and when each enabled channel (Telegram, Discord, etc.) is authorized and active.

## Capabilities

### New Capabilities
- `system-feedback`: Comprehensive logging and status feedback during system lifecycle events (startup, shutdown, component initialization).

### Modified Capabilities
- None (This change focuses on implementation details of existing security and channel startup logic rather than changing their functional requirements).

## Impact

- `internal/app/app.go`: Refactor `New()` and `Start()` for better logging and non-blocking security setup.
- `cmd/lango/main.go`: Enhance top-level logging during the `serve` command.
- `internal/channels/telegram`: Ensure the bot logs its authorized state.
- `internal/cli/prompt`: (Potentially) No changes, but its usage in `app.go` will be refined.
