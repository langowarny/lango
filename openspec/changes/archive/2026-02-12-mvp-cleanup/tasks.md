## 1. Remove Legacy Agent Runtime

- [x] 1.1 Delete `Run()` method (lines 154-368) from `internal/agent/runtime.go`
- [x] 1.2 Retain type definitions: `Runtime`, `Config`, `Tool`, `ToolHandler`, `StreamEvent`, `ParameterDef`, `AdkToolAdapter`
- [x] 1.3 Remove unused imports (`encoding/json`, `time`, `provider`, `session`) from `runtime.go`
- [x] 1.4 Verify no other files call `agent.Runtime.Run()` — grep and confirm zero callers

## 2. Decompose `app.go` into Focused Files

- [x] 2.1 Create `internal/app/wiring.go` — extract `initSupervisor()`, `initSessionStore()`, `initAgent()`, `initGateway()` from `New()`
- [x] 2.2 Create `internal/app/tools.go` — extract `buildTools()` returning `[]*agent.Tool` for exec (exec, exec_bg, exec_status, exec_stop) and filesystem (fs_read, fs_list, fs_write, fs_edit, fs_mkdir, fs_delete) only
- [x] 2.3 Rewrite `internal/app/app.go` — slim `New()` that calls wiring + tools functions, plus `Start()` and `Stop()`
- [x] 2.4 Remove `BrowserSessionID` field from `internal/app/types.go`
- [x] 2.5 Verify no single file in `internal/app/` exceeds 200 lines

## 3. Make Security Non-blocking

- [x] 3.1 In `wiring.go`: when `cfg.Security.Signer.Provider` is empty, skip all security initialization and log info message
- [x] 3.2 When `security.signer.provider` is `"local"` but passphrase cannot be obtained, log warning and continue without security tools (do not return error)
- [x] 3.3 Log deprecation warning when `security.passphrase` is set in config
- [x] 3.4 Remove RPC crypto provider path from application initialization (no `security.NewRPCProvider()` in MVP)

## 4. Remove Non-MVP Tool Registrations

- [x] 4.1 Remove all browser tool imports and registrations (browser_navigate, browser_read, browser_screenshot) from app package
- [x] 4.2 Remove all crypto tool imports and registration from app package
- [x] 4.3 Remove all secrets tool imports and registration from app package
- [x] 4.4 Remove `internal/tools/browser` import from `internal/app/` (but keep the package source)
- [x] 4.5 Remove `internal/tools/crypto` import from `internal/app/`
- [x] 4.6 Remove `internal/tools/secrets` import from `internal/app/`

## 5. Simplify Gateway Initialization

- [x] 5.1 Update `gateway.New()` to accept `nil` for `rpcProvider` and `authManager` parameters without panicking
- [x] 5.2 Remove OIDC auth manager creation from `app.go` (remove `gateway.NewAuthManager` call)
- [x] 5.3 Skip companion WebSocket endpoint registration when `rpcProvider` is nil
- [x] 5.4 Verify gateway starts and serves `chat.message` RPC without auth

## 6. Config Defaults and Simplification

- [x] 6.1 Add default values in config loader for all fields listed in config-system spec (server, session, logging, agent, tools)
- [x] 6.2 Remove `agent.apiKey` field references from config loader if any remain
- [x] 6.3 Update `lango.example.json` to show minimal viable configuration (provider + one channel)
- [x] 6.4 Validate that `agent.provider` references an existing key in `providers` map during config validation

## 7. Dependency Cleanup

- [x] 7.1 Remove `go-rod` related imports from all files
- [x] 7.2 Run `go mod tidy` to remove unused dependencies
- [x] 7.3 Verify `go build ./cmd/lango` compiles successfully

## 8. Test Updates

- [x] 8.1 Update `internal/app/app_test.go` to reflect removed tool registrations (no browser/crypto/secrets)
- [x] 8.2 Update `internal/app/supervisor_test.go` if it references removed security paths
- [x] 8.3 Update `internal/gateway/gateway_test.go` to work without auth manager
- [x] 8.4 Run `go test ./internal/...` and fix any failures
- [x] 8.5 Verify `go vet ./...` passes with no warnings
