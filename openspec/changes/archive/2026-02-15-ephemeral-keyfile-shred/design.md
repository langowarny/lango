## Context

During Docker-based deployments, a plaintext passphrase is written to `~/.lango/keyfile` from Docker secrets. The bootstrap process reads this file, derives a PBKDF2 key, and initializes the crypto provider. After this point the keyfile serves no purpose, but it remains on disk — accessible via the exec tool's `cat` command, which bypasses the filesystem tool's `~/.lango/` access restriction.

## Goals / Non-Goals

**Goals:**
- Eliminate the passphrase keyfile from disk immediately after successful crypto initialization
- Use secure shred (zero-overwrite + sync + remove) rather than simple deletion
- Maintain backward compatibility — no code changes needed for existing callers
- Provide opt-out (`KeepKeyfile`) for debugging or special deployment scenarios

**Non-Goals:**
- Multi-pass overwrite (DoD 5220.22-M style) — single zero-pass is sufficient for application-level security
- Protecting against memory forensics (passphrase string remains in Go heap until GC)
- Modifying `docker-entrypoint.sh` — it already recreates the keyfile on each container start

## Decisions

### Decision 1: Zero-overwrite before removal
**Choice**: Overwrite file content with zero bytes → `Sync()` → `Remove()`
**Rationale**: Simple `os.Remove()` only unlinks the inode; data remains recoverable from disk sectors. Writing zeros ensures the passphrase is not recoverable via simple file recovery tools. Single-pass is sufficient since multi-pass is only relevant for magnetic media with older recording techniques.
**Alternatives considered**: `shred` CLI command (not portable, not available in all containers), crypto-random overwrite (unnecessary overhead for zeroing purpose).

### Decision 2: Shred after checksum verification, not before
**Choice**: Place the shred call after both crypto initialization and checksum verification succeed.
**Rationale**: If the keyfile is shredded before verification and the passphrase is wrong, the user loses access permanently (keyfile is destroyed, passphrase doesn't match). By shredding only after verification, we ensure the passphrase was correct.

### Decision 3: Non-fatal shred failure
**Choice**: Shred failure logs a warning to stderr; bootstrap continues.
**Rationale**: At the point of shredding, the crypto provider is already fully initialized. A shred failure (e.g., permission issue) should not prevent the application from running. The warning alerts operators to investigate.

### Decision 4: Secure-by-default via Go zero value
**Choice**: `KeepKeyfile bool` defaults to `false` (shred).
**Rationale**: Go zero value for bool is `false`, so all existing callers automatically get the secure behavior without code changes. Only callers that explicitly opt out need to change.

## Risks / Trade-offs

- **[Risk]** Keyfile shredded on first correct passphrase entry, but user may want to reuse it for debugging → **Mitigation**: `KeepKeyfile: true` opt-out available
- **[Risk]** Docker secret keyfile not available after shred → **Mitigation**: `docker-entrypoint.sh` recreates keyfile from Docker secret on every container start, so next restart is unaffected
- **[Risk]** Single zero-pass insufficient on some storage media → **Mitigation**: Acceptable for application-level security; full disk encryption should be used for defense-in-depth
