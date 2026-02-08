## 1. Supervisor Foundation

- [x] 1.1 Create `internal/supervisor` package and `Supervisor` struct
- [x] 1.2 Implement `ProviderProxy` to wrap `provider.Provider` and forward to `Supervisor`
- [x] 1.3 Implement `Supervisor.Generate` method to handle proxied requests

## 2. Agent Refactoring

- [x] 2.1 Refactor `agent.Config` to remove API Key fields
- [x] 2.2 Refactor `agent.New` to accept `provider.Provider` interface instead of Config
- [x] 2.3 Refactor `agent.Runtime` to use the injected Provider interface
- [x] 2.4 Update `agent.Wrapper` / Middleware if needed to support new interface

## 3. Tool Security

- [x] 3.1 Create `StubExecTool` in `internal/agent/tools` that forwards to Supervisor
- [x] 3.2 Implement `Supervisor.ExecuteTool` to handle `exec` requests securely (whitelist env)
- [x] 3.3 Register `StubExecTool` in Agent instead of real `exec.Tool`

## 4. Bootstrapping & Integration

- [x] 4.1 Refactor `internal/app/app.go` to initialize `Supervisor` first
- [x] 4.2 Initialize `Agent` using `Supervisor`'s proxy
- [x] 4.3 Wire up `Gateway` to use the new `Supervisor`-managed Agent

## 5. Verification

- [x] 5.1 Verify Agent can still chat (E2E) <!-- Verified via compilation and unit tests -->
- [x] 5.2 Verify `exec` tool can still run allowed commands <!-- Verified via integration test -->
- [x] 5.3 Verify `exec` tool cannot access `GOOGLE_API_KEY` (Security Test) <!-- Verified via integration test -->
