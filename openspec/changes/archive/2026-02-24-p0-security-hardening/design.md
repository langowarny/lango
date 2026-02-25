## Context

The P2P networking layer stores its Ed25519 node key as plaintext binary at `~/.lango/p2p/node.key` with `0600` file permissions. Meanwhile, wallet private keys are encrypted via `SecretsStore` (AES-256-GCM backed by Ent/SQLite) and zeroed from memory after use. This inconsistency means a filesystem compromise exposes P2P identity while wallet identity remains protected.

The handshake protocol (`internal/p2p/handshake/handshake.go`) has a stub at `verifyResponse()` that accepts any non-empty signature as valid, bypassing actual cryptographic verification.

Existing infrastructure: `SecretsStore.Store()/Get()`, `go-ethereum/crypto` (Keccak256, SigToPub, CompressPubkey), `zeroBytes()` pattern in wallet.

## Goals / Non-Goals

**Goals:**
- Encrypt P2P node keys at rest using existing `SecretsStore`
- Auto-migrate legacy plaintext `node.key` files to encrypted storage
- Complete ECDSA secp256k1 signature verification in handshake
- Apply constant-time comparison to prevent timing attacks on nonces
- Zero key material from memory after use
- Deprecate `KeyDir` config field in favor of encrypted storage

**Non-Goals:**
- OS Keyring integration (P1-4, separate change)
- Process isolation for tool execution (P1-5, separate change)
- Session invalidation (P1-6, separate change)
- HSM/Cloud KMS integration (P2-9, separate change)

## Decisions

### D1: SecretsStore for node key storage (not keyring, not separate encryption)

**Choice**: Reuse existing `SecretsStore` (AES-256-GCM + Ent persistence) for P2P node keys.

**Alternatives considered**:
- OS Keyring: Cross-platform complexity, not available in all environments (CI, containers)
- Separate encryption file: Would duplicate crypto infrastructure already in SecretsStore
- Keep file-based with better permissions: Still vulnerable to filesystem compromise

**Rationale**: SecretsStore is battle-tested for wallet keys, available in all environments, and requires zero new dependencies.

### D2: Graceful fallback when SecretsStore unavailable

**Choice**: When `secrets == nil`, retain file-based storage for backward compatibility.

**Rationale**: CLI commands (`lango p2p identity`) may bootstrap without full security initialization. Forcing SecretsStore would break standalone CLI usage.

### D3: Non-blocking migration

**Choice**: Migration failure logs a warning but does not block startup. Retry occurs on next restart.

**Rationale**: A failed migration (e.g., DB locked) should not prevent the node from starting. The plaintext key still works and migration retries automatically.

### D4: ECDSA recovery verification (not signature verification)

**Choice**: Use `ethcrypto.SigToPub()` to recover the public key from the signature, then compare with the claimed `resp.PublicKey`.

**Alternatives considered**:
- `ethcrypto.VerifySignature()`: Does not recover the key, only verifies. Recovery provides stronger proof of identity.

**Rationale**: Recovery-based verification is the standard pattern in Ethereum (used in `ecrecover`), and matches the wallet's `SignMessage` which produces recoverable signatures (65 bytes with V field).

### D5: Constant-time nonce comparison

**Choice**: Replace byte-by-byte nonce comparison with `hmac.Equal()`.

**Rationale**: Prevents timing side-channel attacks where an attacker could determine nonce match length by measuring response time.

## Risks / Trade-offs

- **[Risk] SecretsStore unavailable at CLI time** → Mitigation: `secrets == nil` fallback to file-based storage
- **[Risk] Migration interrupted (crash during store + delete)** → Mitigation: Store first, delete second. On restart, SecretsStore takes priority, so double-storage is harmless
- **[Risk] Legacy plaintext file left behind after failed delete** → Mitigation: Warn log; next restart retries. Key is already safely in SecretsStore
- **[Trade-off] `NewNode()` API change** → All callers must pass `*security.SecretsStore` (can be nil). Breaking change for external consumers, but this is internal API only
