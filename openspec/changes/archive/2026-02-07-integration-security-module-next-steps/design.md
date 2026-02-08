## Context

Security module infrastructure is complete (entities, providers, tools) but disconnected from the application. The security package (`internal/security/`) contains:
- LocalCryptoProvider (AES-256-GCM + PBKDF2)
- CompositeCryptoProvider (primary → fallback)
- KeyRegistry and SecretsStore (entgo.io)
- RPCProvider (companion communication)

app.go has a TODO at L185-188 for LocalProvider initialization.

## Goals / Non-Goals

**Goals:**
- Wire LocalCryptoProvider with passphrase-based initialization
- Register secrets/crypto tools in agent runtime
- Add secrets.get to approval-required list
- Define companion WebSocket protocol
- Add security doctor check

**Non-Goals:**
- Implement iOS/macOS companion app (separate project)
- Hardware security module integration beyond protocol
- Multi-user key management

## Decisions

### 1. Passphrase Storage
**Decision**: Prompt for passphrase on first run, store salt in database, derive key at startup.

**Alternatives considered:**
- Environment variable: Less secure, visible in process list
- Keychain integration: Platform-specific, complex

**Rationale**: PBKDF2 with stored salt balances security and simplicity.

### 2. Tool Registration
**Decision**: Register tools in app.go alongside exec/browser/filesystem tools.

**Rationale**: Consistent with existing tool registration pattern.

### 3. Companion Protocol
**Decision**: JSON-RPC over WebSocket with method prefixes (sign.request, encrypt.request, approval.request).

**Alternatives considered:**
- gRPC: More complex, overkill for mobile client
- REST: No bidirectional push support

**Rationale**: WebSocket enables bidirectional, real-time communication for approval flows.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Passphrase forgotten → data loss | Add passphrase recovery flow in onboard |
| Companion offline during operation | CompositeCryptoProvider falls back to local |
| Salt stored in database | Salt alone is useless without passphrase |
