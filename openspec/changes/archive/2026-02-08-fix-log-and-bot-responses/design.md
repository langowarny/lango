## Context

Currently, Lango's initialization flow in `internal/app/app.go` checks if the session is interactive and immediately prompts for a passphrase if a local salt is found. This blocks background/automated startups even if the `LANGO_PASSPHRASE` environment variable is set. Additionally, there are no "ready" logs for channels, and component initialization progress is hidden.

## Goals / Non-Goals

**Goals:**
- Resolve the startup block by prioritizing environment variables.
- Add granular "Initializing..." and "Ready" logs for all system components.
- Standardize the bot readiness confirmation log.

**Non-Goals:**
- Implementing a persistent "status" API (UI-facing status is already handled).
- Changing the underlying encryption/decryption logic.

## Decisions

### 1. Passphrase Precedence
In `internal/app/app.go`, the check for `LANGO_PASSPHRASE` will be moved to the top of the `LocalCryptoProvider` initialization block. Interactive prompting will only be triggered as a fallback if the environment variable is empty.

### 2. Component Heartbeat Logs
Each major stage of `New()` and `Start()` in `app.go` will include a log entry:
- `Initializing Supervisor...`
- `Initializing Agent...`
- `Registering Tools...`
- `Initializing Channels...`

### 3. Channel Readiness Protocol
Each channel implementation (e.g., `internal/channels/telegram/telegram.go`) will log a specific "ready" message upon successful authorization and connection, including the bot's username/ID for confirmation.

## Risks / Trade-offs

- **Log Verbosity**: Increased logs might clutter extremely high-throughput environments, but since these are startup-only logs, the impact is negligible.
- **Security Logs**: We must strictly avoid logging the passphrase or its derived key. Existing logging patterns already redact sensitive parameters, but extra care will be taken during implementation.
