# Fix Serve Encryption Proposal

## Summary
Update the `lango serve` command to correctly handle encrypted session databases by passing the passphrase to the storage layer.

## Motivation
Currently, `lango serve` crashes with an "out of memory" error (SQLite error 14) when attempting to open an encrypted session database because it initializes `NewEntStore` without providing the encryption key. This prevents users who have enabled encryption from running the server.

## Proposed Solution
1.  **Session Store Update**: Modify `NewEntStore` to accept an optional passphrase and configure the underlying SQLite connection to use it (via PRAGMA key or URI parameters).
2.  **Application Wiring**: Update `internal/app/app.go` to retrieve the passphrase from configuration or environment variables (`LANGO_PASSPHRASE`) and pass it when initializing the store.

## Impact
-   **Reliability**: `lango serve` will start successfully with encrypted databases.
-   **Security**: Maintains encryption at rest while allowing the application to function.
