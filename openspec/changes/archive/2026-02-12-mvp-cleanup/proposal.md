## Why

Lango has accumulated significant complexity through 25+ OpenSpec changes, but still cannot run as a complete service. The root causes are: a dual runtime path (legacy `agent/runtime.go` Run() coexists with ADK), a 720-line God Object (`app.go`), blocking security initialization, 4 incomplete changes leaving code in an inconsistent state, and premature features (browser tool, companion protocol, OIDC, crypto tools) that obscure the core message flow. This cleanup removes everything that doesn't answer "yes" to: **"Does `lango serve` fail without this?"**

## What Changes

- **BREAKING**: Remove legacy agent `Run()` execution loop from `internal/agent/runtime.go`; keep only type definitions (`Tool`, `ToolHandler`, `StreamEvent`). ADK becomes the single runtime path.
- **BREAKING**: Remove browser tool (`internal/tools/browser/`), crypto tool (`internal/tools/crypto/`), and secrets tool (`internal/tools/secrets/`) from MVP build. Remove corresponding registrations from `app.go`.
- **BREAKING**: Remove companion protocol support (`internal/companion/`).
- **BREAKING**: Remove OIDC auth (`internal/cli/auth/`, gateway auth manager). Gateway still serves HTTP/WS but without OIDC.
- **BREAKING**: Remove RPC crypto provider path from `app.go`. Security signer defaults to disabled (no passphrase required to start).
- Decompose `app.go` (720 lines) into focused files: `app.go` (lifecycle), `tools.go` (exec + filesystem registration), `wiring.go` (component assembly).
- Make security initialization non-blocking: if no passphrase configured, skip security tools and log a warning instead of failing startup.
- Simplify configuration: add sensible defaults so only `agent.provider`, `providers.<name>.apiKey`, and one channel token are required to start.
- Merge pending config-json-refactoring changes: remove all `agent.apiKey` references, enforce `providers` map as single source of credentials.
- Remove `go-rod` dependency (browser tool removal).

## Capabilities

### New Capabilities
_None. This change removes complexity, it does not add features._

### Modified Capabilities
- `application-core`: Decompose God Object into focused files; make security optional; remove browser/crypto/secrets tool registration from startup
- `agent-runtime`: Remove `Run()` execution loop; retain only type definitions; ADK is sole runtime
- `config-system`: Add sensible defaults; remove `agent.apiKey` field; simplify minimum viable configuration
- `supervisor-architecture`: Remove RPC crypto provider path; security signer becomes optional
- `tool-browser`: **REMOVED** from MVP build (code retained but not registered)
- `tool-crypto`: **REMOVED** from MVP build (code retained but not registered)
- `tool-secrets`: **REMOVED** from MVP build (code retained but not registered)
- `passphrase-management`: Make non-blocking; skip if unconfigured instead of failing
- `gateway-server`: Remove OIDC auth manager dependency; simplify server initialization

## Impact

- **Code**: `internal/app/app.go` restructured into 3 files; `internal/agent/runtime.go` reduced to types only; tool packages remain but are not wired
- **Dependencies**: `go-rod` removed from `go.mod`; reduces binary size and build complexity
- **Configuration**: Minimum config reduced from ~50 fields to ~5 fields; `lango.example.json` simplified
- **Security**: Crypto/secrets tools unavailable until Phase 2 re-enablement; local passphrase still supported but optional
- **Channels**: Telegram primary; Discord/Slack code retained but focus is on Telegram working end-to-end
- **Tests**: Tests for removed tool registrations need updating; core agent/channel tests must pass
