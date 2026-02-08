## 1. Entity Schema (entgo.io)

- [x] 1.1 Create Key entity schema in `internal/ent/schema/key.go`
- [x] 1.2 Create Secret entity schema in `internal/ent/schema/secret.go`
- [x] 1.3 Run `go generate ./internal/ent` to generate ent code
- [x] 1.4 Add migration for new entities

## 2. Local Fallback Provider

- [x] 2.1 Create LocalCryptoProvider in `internal/security/local_provider.go`
- [x] 2.2 Implement PBKDF2 key derivation from passphrase
- [x] 2.3 Implement AES-256-GCM encrypt/decrypt
- [x] 2.4 Add unit tests for local provider

## 3. Composite Provider

- [x] 3.1 Create CompositeCryptoProvider in `internal/security/composite_provider.go`
- [x] 3.2 Implement fallback logic (companion → local)
- [x] 3.3 Add connection status check interface
- [x] 3.4 Add unit tests for composite provider

## 4. Key Registry Service

- [x] 4.1 Create KeyRegistry service in `internal/security/key_registry.go`
- [x] 4.2 Implement key registration (RegisterKey)
- [x] 4.3 Implement key lookup (GetKey, GetDefaultKey, ListKeys)
- [x] 4.4 Add unit tests for key registry

## 5. Secrets Store

- [x] 5.1 Create SecretsStore in `internal/security/secrets_store.go`
- [x] 5.2 Implement Store operation (encrypt + save)
- [x] 5.3 Implement Get operation (load + decrypt)
- [x] 5.4 Implement List and Delete operations
- [x] 5.5 Add unit tests for secrets store

## 6. Secrets Tool

- [x] 6.1 Create secrets tool in `internal/tools/secrets/secrets.go`
- [x] 6.2 Implement secrets.store handler
- [x] 6.3 Implement secrets.get handler
- [x] 6.4 Implement secrets.list handler
- [x] 6.5 Implement secrets.delete handler
- [x] 6.6 Add unit tests for secrets tool

## 7. Crypto Tool

- [x] 7.1 Create crypto tool in `internal/tools/crypto/crypto.go`
- [x] 7.2 Implement crypto.encrypt handler
- [x] 7.3 Implement crypto.decrypt handler
- [x] 7.4 Implement crypto.sign handler
- [x] 7.5 Implement crypto.hash handler
- [x] 7.6 Implement crypto.keys handler
- [x] 7.7 Add unit tests for crypto tool

## 8. Companion Discovery

- [x] 8.1 Add `grandcat/zeroconf` dependency
- [x] 8.2 Create companion discovery in `internal/companion/discovery.go`
- [x] 8.3 Implement mDNS browser for `_lango-companion._tcp`
- [x] 8.4 Implement service instance handling
- [x] 8.5 Add manual address fallback from config

## 9. Gateway Integration

- [x] 9.1 Add /companion WebSocket endpoint to gateway
- [x] 9.2 Wire RPCProvider.SetSender to companion connection
- [x] 9.3 Route sign/encrypt/decrypt responses to RPCProvider
- [x] 9.4 Implement approval.request broadcast to companions
- [x] 9.5 Handle approval.response from companions

## 10. Agent Runtime Integration

- [x] 10.1 Register secrets tool in agent initialization
- [x] 10.2 Register crypto tool in agent initialization
- [x] 10.3 Add secrets.get to ApprovalMiddleware SensitiveTools
- [x] 10.4 Add integration tests for tool registration

## 11. Configuration

- [x] 11.1 Add security section to config schema
- [x] 11.2 Add companion.address and fallback.enabled options
- [x] 11.3 Update doctor check for security configuration
- [x] 11.4 Add onboard step for security setup (optional)

## 12. Testing & Documentation

- [x] 12.1 Add integration tests for secrets flow
- [x] 12.2 Add integration tests for crypto flow
- [x] 12.3 Update README with security features documentation
