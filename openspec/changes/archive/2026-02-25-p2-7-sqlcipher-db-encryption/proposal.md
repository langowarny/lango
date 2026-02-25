## Why

The SQLite database (`~/.lango/lango.db`) stores session history, configuration profiles, peer reputation scores, and encryption metadata as plaintext. While `SecretsStore` encrypts individual secrets at the application level, the bulk of persisted data remains unencrypted and exposed to any process or user with file-system read access. P0/P1 security phases are complete, making transparent DB encryption the next priority.

## What Changes

- Add `DBEncryptionConfig` to `SecurityConfig` with `enabled` and `cipherPageSize` fields
- Restructure the bootstrap sequence to detect DB encryption, acquire passphrase before opening DB, and pass `PRAGMA key` / `PRAGMA cipher_page_size` to enable SQLCipher transparent encryption
- Export `IsDBEncrypted()` helper for header-based encryption detection (checks SQLite magic bytes)
- Create `internal/dbmigrate` package with `MigrateToEncrypted` and `DecryptToPlaintext` functions using `ATTACH DATABASE ... KEY` + `sqlcipher_export()` workflow
- Add CLI commands: `lango security db-migrate` (plaintext→encrypted) and `lango security db-decrypt` (encrypted→plaintext) with `--force` flag for non-interactive use
- Update `lango security status` to display DB encryption state: "encrypted (active)" / "enabled (pending migration)" / "disabled (plaintext)"
- Update doctor security check to warn when encryption is enabled but DB is still plaintext
- Add `Confirm()` helper to `internal/cli/prompt` for interactive yes/no prompts

## Capabilities

### New Capabilities
- `db-encryption`: Transparent SQLite database encryption via SQLCipher PRAGMA key, including migration tools and bootstrap integration

### Modified Capabilities
- `security-config`: Add `DBEncryption` sub-config with `enabled` and `cipherPageSize` fields
- `bootstrap`: Restructure DB open sequence to support encrypted databases (detect → passphrase → open with key)

## Impact

- **Config**: `internal/config/types.go` (SecurityConfig), `internal/config/loader.go` (defaults)
- **Bootstrap**: `internal/bootstrap/bootstrap.go` — new `openDatabase` signature, `IsDBEncrypted` export, restructured `Run()`
- **New package**: `internal/dbmigrate/` — migration and decryption tools with secure file deletion
- **CLI**: `internal/cli/security/` — db-migrate, db-decrypt commands; status output updated
- **CLI prompt**: `internal/cli/prompt/prompt.go` — added `Confirm()` function
- **Doctor**: `internal/cli/doctor/checks/security.go` — DB encryption status warning
- **Build**: Requires system libsqlcipher for encryption functionality; standard SQLite build operates without encryption (PRAGMA key is no-op)
