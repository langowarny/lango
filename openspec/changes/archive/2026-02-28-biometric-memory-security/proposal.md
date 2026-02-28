## Why

The BiometricProvider (Touch ID + Keychain) implementation stores and retrieves plaintext
passphrases through C/CGo interop but does not zero sensitive memory before freeing it.
This leaves plaintext passphrases lingering in freed heap pages, exposable via memory dumps
or core dumps. This is a HIGH severity gap against secure coding best practices.

## What Changes

- Add `secure_free()` C helper that zeroes memory via volatile pointer before calling `free()`,
  preventing compiler optimization from eliding the wipe
- Change `Get()` to copy Keychain data into a Go `[]byte`, call `secure_free` on the C buffer,
  then zero the Go `[]byte` after extracting the string
- Change `Set()` to zero the `C.CString` buffer via `memset` before freeing it
- Update `SourceKeyring` documentation comment to reflect hardware-backed keyring (Touch ID/TPM)
  instead of generic OS keyring references

## Capabilities

### New Capabilities

_(none — this is a hardening change to existing capability)_

### Modified Capabilities

- `passphrase-acquisition`: Update SourceKeyring comment to reflect hardware keyring terminology
- `keyring-security-tiering`: Add memory zeroing requirements for C interop buffers in BiometricProvider

## Impact

- `internal/keyring/biometric_darwin.go` — C block and Go Get/Set methods
- `internal/security/passphrase/acquire.go` — SourceKeyring comment
- No API changes, no breaking changes, no new dependencies
