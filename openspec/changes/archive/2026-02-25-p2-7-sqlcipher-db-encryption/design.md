## Context

The application database (`~/.lango/lango.db`) uses `mattn/go-sqlite3` with standard SQLite. P0/P1 security hardening is complete (node key encryption, keyring, session invalidation, subprocess sandbox). The DB stores session history, config profiles, peer reputation, and encryption salt/checksum — all as plaintext on disk.

SQLCipher extends SQLite with AES-256-CBC encryption, activated via `PRAGMA key`. The same `mattn/go-sqlite3` driver supports SQLCipher when linked against `libsqlcipher` at build time, making this a zero-code-change-to-driver approach.

## Goals / Non-Goals

**Goals:**
- Transparent encryption of the entire application DB using SQLCipher `PRAGMA key`
- Reversible migration tools: plaintext→encrypted and encrypted→plaintext
- Backwards-compatible bootstrap: unencrypted DBs continue to work when encryption is disabled
- Detection of encryption status via SQLite header magic bytes
- CLI commands for migration and status inspection

**Non-Goals:**
- Replacing the `mattn/go-sqlite3` driver (kept for sqlite-vec compatibility)
- Key management via external KMS (deferred to P2-9)
- Per-column or per-table encryption (SQLCipher encrypts entire DB)
- Automatic migration on first boot (explicit CLI command required)

## Decisions

1. **Keep `mattn/go-sqlite3` driver** — `mutecomm/go-sqlcipher/v4` bundles its own SQLite amalgamation which conflicts with `sqlite-vec-go-bindings` (also CGO SQLite). Instead, link against system `libsqlcipher` at build time; the same driver transparently supports `PRAGMA key`.

2. **Use raw passphrase as DB key** — SQLCipher's internal PBKDF2-HMAC-SHA512 (256K iterations) derives the actual encryption key. Avoids circular dependency: can't use `CryptoProvider` (needs DB open → needs key → needs provider). The passphrase is acquired before DB open.

3. **Bootstrap restructure: passphrase-first** — Detect encryption via header check → acquire passphrase → open DB with `PRAGMA key`. For new DBs, passphrase is acquired first anyway (same path).

4. **Migration via `sqlcipher_export()`** — Atomic: open source → `ATTACH target KEY` → export → DETACH → swap files. Backup with secure delete (zero-overwrite before removal). Verify target DB before removing backup.

5. **`IsDBEncrypted()` via header check** — Standard SQLite files start with "SQLite format 3\0". Encrypted files have random bytes. Simple, reliable, no SQL required.

## Risks / Trade-offs

- **Build dependency**: Encryption requires system `libsqlcipher-dev`. Without it, `PRAGMA key` is silently ignored. The `IsSQLCipherAvailable()` function checks `PRAGMA cipher_version` at runtime.
- **Migration data loss**: Interrupted migration could leave DB in inconsistent state. Mitigated by atomic rename (original→.bak, temp→original) and backup retention on failure.
- **Performance**: SQLCipher adds ~5-15% overhead for encrypted I/O. Acceptable for this use case.
- **sqlite-vec compatibility**: Verified — both use the same underlying SQLite library via CGO. SQLCipher PRAGMAs don't affect sqlite-vec extension loading.
