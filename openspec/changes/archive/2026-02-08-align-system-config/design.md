# Align System Config Design

## Goal
Align the `doctor` command, configuration examples, and documentation with the new TUI capabilities (Security & Providers).

## Architecture

### `doctor` Command Fix
-   **Problem**: `checks/security.go` tries to open the session DB directly without the passphrase, failing on encrypted databases.
-   **Solution**: Update `SecurityCheck.Run`:
    -   Detect if encryption is enabled (check `security.passphrase` or env var).
    -   If enabled, allow checking for DB existence/permissions *without* full decryption if possible, OR prompt/require the passphrase to fully validate.
    -   Alternatively, catch the specific "out of memory" / encryption error and report it as "Database locked (passphrase required)" instead of failing with a scary error.
    -   Given `doctor` is non-interactive usually, we should prioritize graceful failure or clearer error messages over forcing interaction.

### Documentation Updates
-   **`lango.example.json`**:
    -   Add `providers` map with examples for OpenAI/Anthropic.
    -   Expand `security` block to show `interceptor` and `signer` defaults.
-   **`README.md`**:
    -   Add section on "Providers" configuration.
    -   Update "Security" section to reflect recent changes (Rpc/Local/Passphrase).
    -   Mention `lango onboard` capabilities.

## Technical Implementation
-   `internal/cli/doctor/checks/security.go`: Refactor `Run` method.
-   `lango.example.json`: Edit JSON structure.
-   `README.md`: Edit Markdown sections.
