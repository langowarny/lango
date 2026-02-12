## Context

Lango's `internal/app/app.go` is a 720-line function (`New()`) that handles tool definition, security initialization, crypto provider setup, ADK wiring, channel initialization, and gateway creation — all in a single call chain. Meanwhile, `internal/agent/runtime.go` contains a complete `Run()` execution loop (lines 154-368) that duplicates what `internal/adk/agent.go` already does via Google ADK. The codebase also registers browser, crypto, and secrets tools that add complexity and dependencies (go-rod) without being essential for the core message flow: user → channel → agent → provider → response.

Four incomplete OpenSpec changes (`config-json-refactoring`, `fix-log-and-bot-responses`, `adk-go-migration`, `add-oauth-login`) leave the code in an inconsistent intermediate state.

## Goals / Non-Goals

**Goals:**
- `lango serve` starts successfully with minimal configuration (provider API key + channel token)
- Single runtime path: all messages flow through ADK (`internal/adk/agent.go`), no legacy runtime
- `app.go` decomposed into files under 150 lines each with clear responsibilities
- Security is optional: no passphrase → no security tools, but agent still works
- `go build ./cmd/lango` compiles without errors after all changes
- All existing tests pass or are updated to reflect removals

**Non-Goals:**
- Adding new features (this is purely a cleanup)
- Removing tool package source code (packages stay, just not wired in MVP)
- Changing the ADK integration or provider interface
- Modifying channel implementations (Telegram/Discord/Slack code stays as-is)
- Addressing OAuth login (deferred to Phase 2)

## Decisions

### Decision 1: Decompose `app.go` into 3 files

**Choice**: Split `app.go` into `app.go` (lifecycle), `tools.go` (tool registration), `wiring.go` (component assembly).

**Rationale**: The current 720-line `New()` mixes concerns. Splitting by responsibility makes each file independently readable and testable. The `wiring.go` file handles supervisor + provider + ADK creation. The `tools.go` file handles exec + filesystem tool definitions. The `app.go` retains only `App` struct, `Start()`, `Stop()`.

**Alternative considered**: Keep single file but extract helper functions → Still hard to navigate; doesn't solve the concern-mixing problem.

**File mapping**:
| New File | Responsibility | Source lines from current `app.go` |
|----------|---------------|-------------------------------------|
| `app.go` | `New()` orchestrator (calls wiring + tools), `Start()`, `Stop()` | Lines 34-37, 632-658, 662-719 |
| `wiring.go` | `initSupervisor()`, `initSessionStore()`, `initAgent()`, `initGateway()` | Lines 39-63, 620-657 |
| `tools.go` | `buildTools()` returning `[]*agent.Tool` for exec + filesystem only | Lines 66-410 (reduced to exec + fs only) |
| `types.go` | Already exists, remove `BrowserSessionID` field | Existing |
| `channels.go` | Already exists, no changes needed | Existing |

### Decision 2: Remove legacy `Run()` from `agent/runtime.go`

**Choice**: Delete `Run()` method (lines 154-368) and all streaming/provider integration code. Keep type definitions: `Runtime`, `Config`, `Tool`, `ToolHandler`, `StreamEvent`, `ParameterDef`, `AdkToolAdapter`.

**Rationale**: `Run()` duplicates ADK's execution loop. The `adk.Agent.Run()` in `internal/adk/agent.go` is the actual runtime used by `channels.go:runAgent()`. The legacy `Run()` is dead code — nothing calls it. Type definitions are still referenced by `app.go` for tool building.

**Alternative considered**: Keep `Run()` as fallback → Creates confusion about which path is active; violates "single path" principle.

### Decision 3: Remove browser/crypto/secrets tools from wiring (not from source)

**Choice**: Remove tool registrations from `app.go` but keep `internal/tools/browser/`, `internal/tools/crypto/`, `internal/tools/secrets/` packages intact.

**Rationale**: Removing source would make Phase 2 re-enablement harder. By only removing the wiring, we eliminate the go-rod import and all security-blocking initialization from the startup path while preserving the code for later. The `go.mod` can drop `go-rod` since nothing imports it.

### Decision 4: Make security non-blocking

**Choice**: In `wiring.go`, if `cfg.Security.Signer.Provider` is empty or unset, skip all security initialization (no passphrase prompt, no crypto provider, no secrets/crypto tools). Log an info-level message: "security disabled, set security.signer.provider to enable".

**Rationale**: Current `app.go` lines 416-512 block on passphrase input or fail with error if no passphrase in non-interactive mode. This prevents the bot from starting in simple deployment scenarios. Security should be opt-in, not opt-out.

**Flow**:
```
cfg.Security.Signer.Provider == "" → skip security entirely → agent starts
cfg.Security.Signer.Provider == "local" → existing passphrase flow (kept but not blocking agent start on failure — warn and continue)
```

### Decision 5: Simplify config defaults

**Choice**: In config loading, apply these defaults when fields are zero-valued:
- `server.host` → `"localhost"`
- `server.port` → `18789`
- `server.httpEnabled` → `true`
- `server.wsEnabled` → `true`
- `session.databasePath` → `"~/.lango/data.db"`
- `session.maxHistoryTurns` → `100`
- `logging.level` → `"info"`
- `logging.format` → `"console"`
- `agent.maxTokens` → `4096`
- `agent.temperature` → `0.7`

**Rationale**: Currently all these must be explicitly set. With defaults, a 10-line config file is enough to start.

### Decision 6: Remove OIDC and companion from gateway

**Choice**: Remove `authManager` parameter from `gateway.New()`. Remove `internal/cli/auth/` from build. Remove `internal/companion/`.

**Rationale**: OIDC requires external provider configuration (issuerURL, clientID, clientSecret) — unnecessary for MVP. Companion protocol requires native app pairing — Phase 2 feature.

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| Removing legacy `Run()` breaks code that imports it | Grep for all callers first; current analysis shows zero callers outside of runtime.go itself |
| Removing browser tool breaks tests | Update `internal/app/app_test.go` to not expect browser tools; browser package tests stay |
| Users with existing `security.signer.provider: "local"` config get different behavior | Warn in logs when security is configured but initialization fails; don't hard-fail |
| Removing go-rod from go.mod may cascade | Run `go mod tidy` after removing imports; verify build |
| `types.go` still imports `security.CryptoProvider` | Keep the field but make it truly optional (nil is valid) |

## Migration Plan

**Step-by-step execution order** (dependencies flow downward):

1. Remove legacy `Run()` from `agent/runtime.go` (no callers, safe first step)
2. Create `internal/app/wiring.go` — extract supervisor/store/agent/gateway init from `app.go`
3. Create `internal/app/tools.go` — extract exec + filesystem tool definitions only
4. Rewrite `internal/app/app.go` — slim orchestrator calling wiring + tools functions
5. Make security optional in `wiring.go` — skip when unconfigured
6. Remove browser/crypto/secrets tool imports and registrations
7. Remove OIDC auth manager from gateway initialization
8. Remove `BrowserSessionID` from `types.go`
9. Apply config defaults in config loader
10. Update `lango.example.json` to minimal config
11. Run `go mod tidy` to remove unused dependencies (go-rod, etc.)
12. Update tests in `app_test.go` and `supervisor_test.go`
13. Verify: `go build ./cmd/lango && go test ./internal/...`

**Rollback**: All changes are in `internal/` with no external API surface. Git revert restores previous state.

## Open Questions

_None. All decisions are informed by the codebase analysis. Proceed to specs and tasks._
