## 1. Interactive Passphrase Prompt

- [x] 1.1 Add golang.org/x/term dependency
- [x] 1.2 Implement promptPassphrase() function with hidden input
- [x] 1.3 Implement promptPassphraseConfirm() for new passphrase setup
- [x] 1.4 Add TTY detection (isInteractive check)
- [x] 1.5 Update app.go to use interactive prompt instead of config

## 2. Checksum Storage

- [x] 2.1 Add checksum column to security_config table in EntStore
- [x] 2.2 Implement SetChecksum() method in EntStore
- [x] 2.3 Implement GetChecksum() method in EntStore
- [x] 2.4 Update LocalCryptoProvider to verify checksum on initialization
- [x] 2.5 Store checksum on first-time passphrase setup

## 3. Migration Command

- [x] 3.1 Create internal/cli/security package and migrate command structure
- [x] 3.2 Implement logic to re-encrypt secrets with new passphrase
- [x] 3.3 Implement transactional update of salt and checksum
- [x] 3.4 Add rollback mechanism on failure
- [ ] 3.5 Implement rollback on failure
- [ ] 3.6 Update salt and checksum after successful migration

## 4. Config Deprecation & Cleanup

- [x] 4.1 Remove passphrase reading from config in app.go
- [x] 4.2 Add deprecation warning when config contains passphrase field
- [ ] 4.3 Update config documentation

## 5. Doctor Checks

- [ ] 5.1 Add Security Provider Mode check in security.go
- [x] 5.1 Add check for insecure config (passphrase in config)
- [x] 5.2 Add check for production usage of local provider
- [x] 5.3 Verify checksum existence in doctor check
- [ ] 5.4 Add tests for new doctor checks

## 6. Documentation & Final Review
- [x] 6.1 Update README Security section with two modes
- [x] 6.2 Document passphrase migration process
- [x] 6.3 Add warning about passphrase loss = data loss
```
