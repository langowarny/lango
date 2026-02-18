## 1. Security Initialization Fix

- [x] 1.1 In `internal/app/app.go`, refactor the `LocalCryptoProvider` initialization to check `LANGO_PASSPHRASE` before attempting an interactive prompt.
- [x] 1.2 Ensure the system gracefully handles the absence of a passphrase in non-interactive environments by exiting with a clear error.

## 2. Enhanced System Logging

- [x] 2.1 Add `logger.Info` or `logger.Infow` heartbeat logs for the following lifecycle events in `internal/app/app.go`:
    - Starting Supervisor
    - Initializing Session Store
    - Initializing Agent Runtime
    - Registering Tools
    - Initializing Channels
- [x] 2.2 Standardize the Gateway startup log to include the full listening address.
- [x] 2.3 Ensure the Telegram channel (and others) logs a "Success" message upon successful authorization/connection.

## 3. Verification

- [x] 3.1 Manually verify that `lango serve` starts without user intervention when `LANGO_PASSPHRASE` is set.
- [x] 3.2 Verify that the terminal output shows a clear, sequential log of component initialization.
- [x] 3.3 Verify that the Telegram bot confirms its readiness in the console.
