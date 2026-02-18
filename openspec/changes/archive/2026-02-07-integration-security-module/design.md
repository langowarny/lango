## Context

Lango has `internal/security` with a `CryptoProvider` interface and `RPCProvider` implementation for external cryptographic operations. Currently:
- RPCProvider exists but has no sender configured (not connected to any transport)
- No tools expose cryptographic operations to AI agents
- No entity schema for storing encrypted secrets
- ApprovalMiddleware exists but isn't configured for security-sensitive operations

The companion app (iOS/macOS with Secure Enclave) is a separate project. This change focuses on the Lango server side.

## Goals / Non-Goals

**Goals:**
- AI agents can securely store and retrieve secrets via `secrets` tool
- AI agents can perform cryptographic operations via `crypto` tool  
- Key metadata stored in SQLite via entgo.io
- RPCProvider connected to Gateway WebSocket for companion communication
- Local fallback encryption when companion unavailable
- Bonjour discovery to find companion app on local network
- User approval required for `secrets.get` operations

**Non-Goals:**
- Companion app implementation (separate project)
- HSM/cloud KMS integration (future enhancement)
- Secret rotation automation
- Multi-user/team secret sharing

## Decisions

### 1. Key and Secret Storage: entgo.io Entities

**Decision**: Add `Key` and `Secret` entities to existing ent schema.

**Alternatives considered**:
- Separate SQLite database: More isolation but duplicates connection management
- JSON file: Simple but no query capability, race conditions

**Rationale**: Reuses existing ent infrastructure. Key stores metadata (name, remoteKeyId, type). Secret stores encrypted blob with FK to Key.

### 2. Local Fallback: AES-256-GCM with PBKDF2

**Decision**: When companion unavailable, use local encryption with key derived from user-provided passphrase.

**Alternatives considered**:
- No fallback (companion required): Poor UX for initial setup
- OS keychain: Platform-specific, CGO dependency
- Random file-based key: No user authentication

**Rationale**: PBKDF2 with passphrase provides auth without platform dependencies. Encrypted secrets can later be re-encrypted with companion when available.

### 3. Companion Discovery: Bonjour/mDNS + WebSocket

**Decision**: Use `grandcat/zeroconf` for mDNS discovery. Companion advertises `_lango-companion._tcp`. Connect via WebSocket with mTLS.

**Alternatives considered**:
- Manual IP configuration: Poor UX
- Cloud relay: Adds latency, privacy concerns
- Bluetooth: Limited range, complex pairing

**Rationale**: Bonjour is native to Apple ecosystem, works over local network, zero configuration.

### 4. Tool Registration: Separate secrets and crypto Tools

**Decision**: Two distinct tools with clear responsibilities.

```
secrets tool:
  - secrets.store(name, value) → stores encrypted
  - secrets.get(name) → returns decrypted (REQUIRES APPROVAL)
  - secrets.list() → returns names only
  - secrets.delete(name)

crypto tool:
  - crypto.encrypt(data, keyId) → returns ciphertext
  - crypto.decrypt(data, keyId) → returns plaintext
  - crypto.sign(data, keyId) → returns signature
  - crypto.hash(data, algorithm) → returns hash
```

**Rationale**: Separation of concerns. Secrets tool is higher-level (automatic key selection). Crypto tool exposes raw operations for advanced use.

### 5. Approval Integration: Extend ApprovalMiddleware Config

**Decision**: Add `secrets.get` to SensitiveTools list. Approval request sent via WebSocket to connected clients (companion app, web UI).

**Alternatives considered**:
- Per-secret approval settings: Over-complicated
- No approval: Security risk

**Rationale**: Uses existing ApprovalMiddleware. Simple configuration.

### 6. CryptoProvider Selection: Strategy Pattern

**Decision**: Create `CompositeCryptoProvider` that tries companion first, falls back to local.

```go
type CompositeCryptoProvider struct {
    companion CryptoProvider  // RPCProvider
    local     CryptoProvider  // LocalProvider (new)
}

func (c *CompositeCryptoProvider) Encrypt(ctx, keyId, data) {
    if c.companion.IsConnected() {
        return c.companion.Encrypt(...)
    }
    return c.local.Encrypt(...)
}
```

**Rationale**: Transparent fallback. Caller doesn't need to know which provider is used.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Local passphrase weak | Enforce minimum length, warn user |
| Companion connection drops mid-operation | Timeout and retry, or fail explicitly |
| mDNS blocked on network | Allow manual IP fallback in config |
| Secrets accessible to anyone with DB file | DB file protected by OS permissions; values encrypted |
| Key rotation complexity | Out of scope; manual re-encryption if needed |

## Open Questions

1. **Passphrase storage**: Should we remember the local passphrase in session or require every time?
   - Proposal: Cache in memory for session duration, clear on shutdown

2. **Companion authentication**: How does companion prove identity?
   - Proposal: mTLS with pre-shared certificate generated during pairing

3. **Multi-key UX**: How does AI know which key to use?
   - Proposal: Default key for secrets tool, explicit keyId for crypto tool
