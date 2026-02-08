## 1. Interface & Core Logic

- [x] 1.1 Rename `internal/security/signer.go` to `internal/security/crypto.go`.
- [x] 1.2 Update `Signer` interface to `CryptoProvider` and add `Encrypt`/`Decrypt` methods.
- [x] 1.3 Update `internal/app` and other consumers to reference `CryptoProvider` instead of `Signer`.

## 2. RPC Implementation

- [x] 2.1 Update `internal/security/rpc_signer.go` to `rpc_provider.go` and implement `Encrypt` / `Decrypt`.
- [x] 2.2 Add `EncryptRequest`, `EncryptResponse`, `DecryptRequest`, `DecryptResponse` structs.
- [x] 2.3 Implement async request handling for new methods in `rpc_provider.go`.

## 3. Gateway Integration

- [x] 3.1 Update `internal/gateway/server.go` to register `encrypt.response` and `decrypt.response` handlers.
- [x] 3.2 Update `handleSignResponse` pattern to support generic or specific response handlers for new types.

## 4. Verification

- [x] 4.1 Update `internal/security/rpc_provider_test.go` (if exists) or create new tests mocking the sender.
- [x] 4.2 Verify manual end-to-end flow with Host App (if available) or mock response. (Manual)
