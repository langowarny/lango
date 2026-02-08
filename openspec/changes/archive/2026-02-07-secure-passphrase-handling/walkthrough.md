# Secure Passphrase Handling

This change enhances the security of the Local Crypto Provider by implementing interactive passphrase entry, checksum verification, and a migration path for rotating keys.

## Changes Implemented

### 1. Interactive Passphrase Prompt
- Removed deprecated `security.passphrase` config usage.
- Implemented `prompt.Passphrase` using `golang.org/x/term` for secure input.
- Added `LANGO_PASSPHRASE` environment variable fallback for non-interactive environments.

### 2. Checksum Verification
- Added `checksum` column to `security_config` table.
- Implemented `CalculateChecksum` (SHA256 of passphrase + salt) to verify passphrase correctness before initialization.
- Added checksum validation logic in `app.go`.

### 3. Migration Command
- Added `lango security migrate-passphrase` command.
- Securely rotates keys by re-encrypting all secrets with a new passphrase and salt.
- Updates salt and checksum atomically within a single transaction.

### 4. Doctor Checks
- Implemented `SecurityCheck` in `lango doctor`.
- Warns about insecure `local` provider usage in production.
- Detects deprecated config usage.
- Validates checksum existence.

## Verification

### Manual Testing

1. **Interactive Mode**:
   - Run `lango serve` without `LANGO_PASSPHRASE`.
   - Verify prompt appears and hides input.
   - Verify incorrect passphrase fails initialization.

2. **Environment Variable Mode**:
   - Run `LANGO_PASSPHRASE=correct lango serve`.
   - Verify it starts without prompt.
   - Run `LANGO_PASSPHRASE=wrong lango serve`.
   - Verify it fails with checksum error.

3. **Migration**:
   - Run `lango security migrate-passphrase`.
   - Follow prompts to rotate key.
   - Verify secrets are still accessible after restart with new passphrase.

4. **Doctor**:
   - Run `lango doctor`.
   - Verify Security Configuration section appears with appropriate status (Pass/Warn).

## Future Work
- Implement remote RPC signer integration.
- Add hardware security module support via companion app.
