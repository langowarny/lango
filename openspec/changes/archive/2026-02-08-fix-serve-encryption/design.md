# Fix Serve Encryption Design

## Architecture

### Session Store (`internal/session`)
-   **`EntStore`**: Will be updated to support Functional Options pattern (`StoreOption`).
-   **`NewEntStore`**: Signature updated to `NewEntStore(dbPath string, opts ...StoreOption)`.
-   **`WithPassphrase`**: New option function to set the passphrase field in `EntStore`.
-   **Connection Logic**:
    -   If a passphrase is provided, append `_pragma_key` to the DSN (if supported) AND execute `PRAGMA key = '...'` after connection for broader driver compatibility.

### Application Wiring (`internal/app`)
-   **`App` Initialization**:
    -   Read `LANGO_PASSPHRASE` from environment or `config.Security.Passphrase`.
    -   Pass `session.WithPassphrase(...)` to `NewEntStore` during initialization.

## Technical Details
-   **Driver Compatibility**: Switched to `mattn/go-sqlite3` to support SQLCipher via `PRAGMA key`. This requires CGO enabled and `sqlcipher` installed on the system.
-   **Path Expansion**: Added support for `~` expansion in database paths to resolve "no such file or directory" errors when using home-relative paths in configuration.

## Risks
-   **Dependency**: Users must have `sqlcipher` installed for encrypted database support. Pure Go build without tags will fail to open encrypted databases.
