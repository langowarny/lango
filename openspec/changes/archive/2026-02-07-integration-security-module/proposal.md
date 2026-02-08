## Why

Lango has a security module (`internal/security`) with `CryptoProvider` interface and `RPCProvider` implementation for Zero Trust cryptographic operations. However, it's not integrated into the runtime - there's no way for AI agents to use secure secrets storage, and the RPC provider isn't connected to any transport. This change integrates security throughout Lango, enabling hardware-backed encryption via iOS/macOS companion app while maintaining local fallback.

## What Changes

- Add `secrets` and `crypto` tools for AI agent use
- Create Key and Secret entities in entgo.io for encrypted data storage
- Integrate RPCProvider with Gateway WebSocket transport
- Add local encryption fallback when companion app unavailable
- Add Bonjour service discovery for companion app
- Configure ApprovalMiddleware to require user approval for sensitive operations (secrets.get)

## Capabilities

### New Capabilities
- `tool-secrets`: Secure secrets management tool (store, get, list, delete) with encrypted storage
- `tool-crypto`: Cryptographic operations tool (encrypt, decrypt, sign, hash) exposed to AI
- `key-registry`: Key metadata storage using entgo.io for managing encryption keys
- `companion-discovery`: Bonjour/mDNS service discovery for iOS/macOS companion app

### Modified Capabilities
- `secure-signer`: Extend to support local fallback encryption when companion unavailable
- `agent-runtime`: Register new security tools and configure approval for sensitive tools
- `gateway-server`: Wire RPCProvider sender to WebSocket connection for companion communication

## Impact

- **Code**: `internal/security/`, `internal/tools/`, `internal/ent/schema/`, `internal/agent/`, `internal/gateway/`
- **Dependencies**: May need `grandcat/zeroconf` for Bonjour discovery
- **Config**: New `security` section in lango.json for fallback mode and companion settings
- **UX**: AI can securely store/retrieve secrets; companion app receives approval requests
