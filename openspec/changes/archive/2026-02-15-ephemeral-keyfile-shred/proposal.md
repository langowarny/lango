## Why

The plaintext passphrase keyfile (`~/.lango/keyfile`) persists on disk after bootstrap completes. While the filesystem tool blocks `~/.lango/` access, the exec tool can still read it via `cat ~/.lango/keyfile`, allowing an agent to extract the passphrase. Since the keyfile is only needed to derive the PBKDF2 key during bootstrap, it should be securely shredded immediately after successful crypto initialization.

## What Changes

- Add `ShredKeyfile()` function that overwrites file content with zeros, syncs to disk, and removes the file (idempotent — returns nil if file doesn't exist)
- Add `KeepKeyfile` option to bootstrap `Options` (default `false` = secure by default)
- Bootstrap `Run()` shreds the keyfile after successful crypto initialization and checksum verification when the passphrase source is keyfile
- Shred failure emits a stderr warning but does not block bootstrap (crypto is already initialized)

## Capabilities

### New Capabilities
- `keyfile-shred`: Secure zero-overwrite and deletion of passphrase keyfile after bootstrap crypto initialization

### Modified Capabilities
- `bootstrap-lifecycle`: Bootstrap now shreds the keyfile after crypto init when source is keyfile, with opt-out via `KeepKeyfile`
- `passphrase-acquisition`: Adds `ShredKeyfile()` to the keyfile management surface

## Impact

- `internal/passphrase/keyfile.go` — new `ShredKeyfile()` function
- `internal/bootstrap/bootstrap.go` — `KeepKeyfile` field in `Options`, shred call in `Run()`
- Docker workflow unaffected — `docker-entrypoint.sh` recreates keyfile from Docker secrets on each container start
- No breaking changes — existing callers get secure-by-default behavior via Go zero value
