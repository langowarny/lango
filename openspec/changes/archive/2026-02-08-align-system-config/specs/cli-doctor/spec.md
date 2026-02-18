# Align System Config Spec

## Goal
Ensure verify system health checks and documentation accuracy for new security/provider features.

## Requirements

### CLI Doctor
-   **Graceful Failure**: The `doctor` command MUST NOT crash or return obscure system errors (like "out of memory") when encountering an encrypted database without credentials.
-   **Clear Feedback**: If the database is locked, `doctor` SHOULD report "Session database encrypted" or similar, and suggest how to unlock (e.g., set `LANGO_PASSPHRASE`).
-   **Security Check**: The `Security Configuration` check MUST verify that the `security` block in `lango.json` is valid according to the new schema.

### Documentation
-   **Example Config**: `lango.example.json` MUST include:
    -   A valid `providers` map example.
    -   A complete `security` block with `interceptor` and `signer` fields.
-   **README**: The main documentation MUST explain:
    -   How to configure multiple providers via `lango.json`.
    -   The role of the TUI (`lango onboard`) in managing these settings.
