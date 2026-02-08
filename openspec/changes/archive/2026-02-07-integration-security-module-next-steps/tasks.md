## 1. LocalCryptoProvider Integration

- [x] 1.1 Add salt storage method to session store interface
- [x] 1.2 Implement salt persistence in EntStore
- [x] 1.3 Update app.go to initialize LocalCryptoProvider when provider is "local"
- [x] 1.4 Add passphrase prompt during initialization
- [x] 1.5 Add unit tests for salt persistence

## 2. Security Tools Registration

- [x] 2.1 Create secrets tool wrapper with proper parameters
- [x] 2.2 Create crypto tool wrapper with proper parameters
- [x] 2.3 Register secrets tool in app.go
- [x] 2.4 Register crypto tool in app.go
- [x] 2.5 Add integration test for tool registration

## 3. Approval Middleware Update

- [x] 3.1 Add "secrets.get" to default SensitiveTools list
- [x] 3.2 Add config option for additional sensitive tools
- [x] 3.3 Add unit test for approval requirement

## 4. Companion Protocol

- [x] 4.1 Add /companion WebSocket endpoint to gateway
- [x] 4.2 Implement companion.hello handler
- [x] 4.3 Implement approval.request/response handlers
- [x] 4.4 Wire companion connection to CompositeCryptoProvider
- [x] 4.5 Add integration test for companion handshake

## 5. Doctor Check

- [x] 5.1 Create security check in internal/cli/doctor/checks/security.go
- [x] 5.2 Check crypto provider initialization status
- [x] 5.3 Check companion connection status
- [x] 5.4 Register check in checks.go
