## 1. P2P Node Key Encrypted Storage

- [x] 1.1 Add `security` import and `nodeKeySecret` constant to `internal/p2p/node.go`
- [x] 1.2 Change `NewNode()` signature to accept `*security.SecretsStore` parameter
- [x] 1.3 Refactor `loadOrGenerateKey()` with SecretsStore priority: secrets → legacy file → generate
- [x] 1.4 Implement `migrateKeyToSecrets()` for auto-migration of legacy plaintext keys
- [x] 1.5 Add `zeroBytes()` function and apply `defer zeroBytes(data)` to all key material paths

## 2. Handshake Signature Verification

- [x] 2.1 Add `ethcrypto`, `bytes`, `hmac` imports to `internal/p2p/handshake/handshake.go`
- [x] 2.2 Replace nonce comparison with `hmac.Equal()` for constant-time comparison
- [x] 2.3 Implement ECDSA recovery verification: Keccak256 hash → SigToPub → CompressPubkey comparison
- [x] 2.4 Add 65-byte signature length validation

## 3. Wiring and CLI Updates

- [x] 3.1 Add `*security.SecretsStore` parameter to `initP2P()` in `internal/app/wiring.go`
- [x] 3.2 Update `initP2P()` call in `internal/app/app.go` to pass `app.Secrets`
- [x] 3.3 Build `SecretsStore` from bootstrap result in `initP2PDeps()` (`internal/cli/p2p/p2p.go`)
- [x] 3.4 Replace `keyDir` output with `keyStorage` in `internal/cli/p2p/identity.go`

## 4. Config Cleanup

- [x] 4.1 Add `omitempty` to `KeyDir` json tag and deprecated comment in `internal/config/types.go`
- [x] 4.2 Add `nodeKeyName` default in `internal/config/loader.go`

## 5. Tests

- [x] 5.1 Write `TestVerifyResponse_ValidSignature` — valid sig accepted
- [x] 5.2 Write `TestVerifyResponse_InvalidSignature` — pubkey mismatch rejected
- [x] 5.3 Write `TestVerifyResponse_WrongSignatureLength` — non-65-byte rejected
- [x] 5.4 Write `TestVerifyResponse_NonceMismatch` — constant-time nonce rejection
- [x] 5.5 Write `TestVerifyResponse_NoProofOrSignature` — empty response rejected
- [x] 5.6 Write `TestVerifyResponse_CorruptedSignature` — corrupted sig rejected
- [x] 5.7 Write `TestLoadOrGenerateKey_NewKeyWithoutSecrets` — file fallback
- [x] 5.8 Write `TestLoadOrGenerateKey_LegacyFileLoaded` — legacy file loading
- [x] 5.9 Write `TestZeroBytes` — memory zeroing verification

## 6. Verification

- [x] 6.1 `go build ./...` passes with no errors
- [x] 6.2 `go test ./internal/p2p/...` all tests pass
- [x] 6.3 `go test ./internal/p2p/handshake/...` all tests pass
