## Phase 1: Quick Wins — Config Wiring

- [x] 1.1 SystemPromptPath loading (M1)
  - Added `loadSystemPrompt()` in wiring.go
  - Reads file via `os.ReadFile()`, falls back to default on error
  - Passed to ContextAwareModelAdapter and NewAgent

- [x] 1.2 Temperature/MaxTokens passing (M2)
  - Added `WithTemperature()`, `WithMaxTokens()` ProxyOptions to ProviderProxy
  - Apply defaults in `Generate()` when request params are zero
  - Wired in `initAgent()` from `cfg.Agent.Temperature/MaxTokens`

- [x] 1.3 HTTPEnabled flag (M5)
  - Wrapped `/health` and `/status` registration in `if s.config.HTTPEnabled`
  - WebSocket routes remain separately controlled by `WebSocketEnabled`

- [x] 1.4 Security CLI registration (M7)
  - Changed `NewSecurityCmd()` to accept `func() (*config.Config, error)` closure
  - Registered in main.go with lazy config loading

- [x] 1.5 Learning sessionKey fix (M9)
  - Defined `sessionKeyCtxKey` type and `WithSessionKey()`/`SessionKeyFromContext()` helpers
  - `wrapWithLearning()` extracts session key from context

## Phase 2: Core Feature Wiring

- [x] 2.1 AuthManager wiring (C1)
  - Added `initAuth()` function: creates AuthManager when `cfg.Auth.Providers` is non-empty
  - Updated `initGateway()` signature to accept `*gateway.AuthManager`
  - Gateway already handles `auth != nil` for route registration

- [x] 2.2 Crypto/Secrets Tools registration (C2)
  - Added `buildCryptoTools()`: 5 tools (encrypt, decrypt, sign, hash, keys)
  - Added `buildSecretsTools()`: 4 tools (store, get, list, delete)
  - Registered in app.go when `app.Crypto != nil` / `app.Secrets != nil`

- [x] 2.3 Fallback Provider (C5)
  - Added `WithFallback()` ProxyOption to ProviderProxy
  - `Generate()` retries with fallback on primary failure (single retry)
  - Wired in `initAgent()` from `cfg.Agent.FallbackProvider/FallbackModel`

- [x] 2.4 MaxHistoryTurns (M3)
  - Added `WithMaxHistoryTurns()` StoreOption to EntStore
  - After `AppendMessage()`, counts messages and deletes oldest excess
  - Wired in `initSessionStore()` from `cfg.Session.MaxHistoryTurns`

- [x] 2.5 Session TTL (M4)
  - Added `WithTTL()` StoreOption to EntStore
  - `Get()` returns "session expired" when `time.Since(UpdatedAt) > ttl`
  - Wired in `initSessionStore()` from `cfg.Session.TTL`

- [x] 2.6 MaxSkillsPerDay (M6)
  - Added `maxSkillsPerDay` field and daily counter to knowledge Store
  - `SaveSkill()` calls `reserveSkillSlot()` to enforce daily limit
  - Updated `NewStore()` signature to accept `maxSkillsPerDay`
  - Wired in `initKnowledge()` from `cfg.Knowledge.MaxSkillsPerDay`

## Phase 3: Dead Code Cleanup

- [x] 3.1 Companion package deletion (C3)
  - Deleted `internal/companion/discovery.go` (254 LOC)
  - `go mod tidy` removed `zeroconf` dependency
  - Gateway companion WebSocket handler retained (direct connection mode)

- [x] 3.2 PairingEnabled removal (M8)
  - Removed `PairingEnabled` from `config.TelegramConfig`
  - Removed `PairingEnabled` from `telegram.Config`
  - Removed `PairingEnabled` assignment from `app/channels.go`

- [x] 3.3 Tool Approval fail-open (C4)
  - Added `wrapWithApproval()` in tools.go
  - Added `HasCompanions()` method to Gateway Server
  - Added `SetAgent()` for deferred agent wiring
  - Wired in app.go: Gateway created before Agent, tools wrapped, Agent set after

## Verification

- [x] `go build ./...` — compilation successful
- [x] `go vet ./...` — no warnings
- [x] `go test ./...` — all tests pass
- [x] `go mod tidy` — zeroconf dependency removed
