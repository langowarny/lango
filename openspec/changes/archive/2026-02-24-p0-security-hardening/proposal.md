## Why

The P2P node key (`~/.lango/p2p/node.key`) is stored as plaintext binary protected only by file permissions (`0600`), while wallet keys are properly encrypted in `SecretsStore` (AES-256-GCM). This architectural inconsistency creates a critical security gap. Additionally, the handshake signature verification stub accepts any non-empty signature, and the `KeyDir` config field unnecessarily exposes the key storage path.

## What Changes

- Migrate P2P node key storage from plaintext file to `SecretsStore` (encrypted), with auto-migration of legacy files and fallback for backward compatibility
- Complete ECDSA signature verification in handshake using `go-ethereum/crypto` public key recovery, replacing the stub that accepted any non-empty signature
- Apply `zeroBytes()` memory cleanup pattern to P2P node key material (matching wallet key handling)
- Deprecate `KeyDir` config field and replace identity CLI output with `keyStorage` info
- Add constant-time nonce comparison (`hmac.Equal`) to prevent timing attacks

## Capabilities

### New Capabilities
- `p2p-node-key-encryption`: Encrypted storage of P2P node keys in SecretsStore with auto-migration from legacy plaintext files

### Modified Capabilities
- `p2p-handshake`: Complete ECDSA signature verification with secp256k1 public key recovery, constant-time nonce comparison, and 65-byte signature length validation
- `p2p-identity`: Replace `keyDir` output with `keyStorage` info reflecting encrypted vs file-based storage
- `p2p-networking`: Accept `*security.SecretsStore` parameter for encrypted node key management

## Impact

- **Code**: `internal/p2p/node.go`, `internal/p2p/handshake/handshake.go`, `internal/app/wiring.go`, `internal/app/app.go`, `internal/cli/p2p/p2p.go`, `internal/cli/p2p/identity.go`, `internal/config/types.go`, `internal/config/loader.go`
- **APIs**: `NewNode()` signature gains `*security.SecretsStore` parameter; `initP2P()` gains `*security.SecretsStore` parameter
- **Dependencies**: No new dependencies (reuses existing `go-ethereum/crypto`)
- **Migration**: Automatic on startup (legacy `node.key` detected, stored in SecretsStore, plaintext deleted); failure is non-fatal (warn log, retry on next restart)
