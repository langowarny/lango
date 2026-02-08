## Why

The security module infrastructure was implemented (entities, providers, tools) but is not yet wired into the application. AI agents cannot use security features until app.go connects LocalCryptoProvider, secrets.get requires approval, and the companion protocol is defined.

## What Changes

- Wire `LocalCryptoProvider` into app.go initialization (with passphrase prompt)
- Register `secrets` and `crypto` tools in agent runtime
- Add `secrets.get` to ApprovalMiddleware SensitiveTools list
- Define WebSocket protocol for companion app communication
- Add security doctor check to verify provider status

## Capabilities

### New Capabilities
- `companion-protocol`: WebSocket message types and handshake for companion app integration
- `security-tools`: Registration and configuration of secrets/crypto tools in agent runtime

### Modified Capabilities
- `session-store`: Add method to retrieve encryption salt for LocalCryptoProvider persistence

## Impact

- `internal/app/app.go` - LocalCryptoProvider initialization and tool registration
- `internal/security/` - Salt persistence interface
- `internal/agent/approval_middleware.go` - SensitiveTools list update
- `internal/gateway/` - Companion WebSocket handler
- `internal/cli/doctor/checks/` - New security check
