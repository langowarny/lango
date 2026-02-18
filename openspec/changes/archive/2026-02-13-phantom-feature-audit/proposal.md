## Why

A post-security-audit analysis revealed **15 phantom features** across the codebase: config fields that users can set but have zero runtime effect because the implementation code was never wired into the application lifecycle. This creates a dangerous false sense of security and functionality.

Additionally, **3 dead code packages** were identified that exist in the repo but are never imported.

## What Changes

### Config-to-Runtime Wiring (13 fixes)

- **SystemPromptPath**: Load custom system prompt from file instead of using hardcoded default
- **Temperature/MaxTokens**: Pass agent config values through ProviderProxy to LLM requests
- **FallbackProvider**: Implement single-retry failover when primary provider fails
- **AuthManager**: Wire OIDC auth into Gateway when providers are configured
- **Crypto/Secrets Tools**: Register crypto and secrets tool implementations as agent tools
- **MaxHistoryTurns**: Trim oldest messages when session history exceeds configured limit
- **Session TTL**: Return expired error when session exceeds TTL on Get()
- **HTTPEnabled**: Make /health and /status endpoints conditional on the flag
- **MaxSkillsPerDay**: Enforce daily rate limit on skill creation
- **Security CLI**: Register `lango security migrate-passphrase` command in main.go
- **Learning sessionKey**: Extract session key from context instead of passing empty string
- **Tool Approval**: Implement fail-open approval flow via companion WebSocket

### Dead Code Removal (3 fixes)

- **Companion Discovery**: Delete unused mDNS discovery package (~254 LOC)
- **PairingEnabled**: Remove phantom config field from Telegram channel
- **zeroconf dependency**: Remove via `go mod tidy`

## Capabilities

### New Capabilities
- `phantom-feature-wiring`: All 15 config-exposed features now have runtime effect

### Modified Capabilities
- `agent-runtime`: SystemPromptPath, Temperature, MaxTokens, FallbackProvider now respected
- `session-store`: TTL expiration and MaxHistoryTurns truncation
- `gateway-server`: HTTPEnabled conditional routing, HasCompanions() method, SetAgent() deferred wiring
- `security-tools`: Crypto and secrets tools registered as agent tools
- `config-system`: PairingEnabled removed from TelegramConfig

## Impact

| File | Changes |
|------|---------|
| `internal/supervisor/proxy.go` | ProxyOption pattern: WithTemperature, WithMaxTokens, WithFallback |
| `internal/session/ent_store.go` | WithMaxHistoryTurns, WithTTL store options; truncation + expiry logic |
| `internal/knowledge/store.go` | maxSkillsPerDay parameter; daily rate limiting on SaveSkill |
| `internal/gateway/server.go` | Conditional HTTP routes; HasCompanions(); SetAgent() |
| `internal/app/tools.go` | Context-based sessionKey; buildCryptoTools; buildSecretsTools; wrapWithApproval |
| `internal/app/wiring.go` | loadSystemPrompt; initAuth; proxy options; session store options |
| `internal/app/app.go` | Auth wiring; crypto/secrets tool registration; approval wrapping |
| `internal/app/channels.go` | PairingEnabled reference removed |
| `internal/config/types.go` | PairingEnabled field removed |
| `internal/channels/telegram/telegram.go` | PairingEnabled field removed |
| `internal/cli/security/migrate.go` | Lazy config loading via closure |
| `cmd/lango/main.go` | Security CLI command registered |
| `internal/companion/discovery.go` | **Deleted** |
