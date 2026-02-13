## Architecture Decisions

### 1. Functional Options Pattern for ProviderProxy

**Decision**: Use interface-based functional options (ProxyOption) instead of expanding the constructor signature.

**Rationale**: The proxy already has 2 required params (providerID, defaultModel). Adding 4 more optional params (temperature, maxTokens, fallbackProvider, fallbackModel) would make the constructor unwieldy. The functional options pattern aligns with Go best practices and the existing StoreOption pattern in session package.

```go
type ProxyOption interface {
    apply(*proxyOptions)
}

proxy := NewProviderProxy(sv, providerID, model,
    WithTemperature(0.7),
    WithMaxTokens(4096),
    WithFallback("openai", "gpt-4"),
)
```

### 2. Fail-Open Approval for Sensitive Tools

**Decision**: When `ApprovalRequired=true` but no companion is connected, log a warning and proceed (fail-open) rather than blocking.

**Rationale**: Blocking would break all tool execution in standalone mode. Security-conscious deployments connect a companion; others get warnings in logs. This matches the existing companion-optional architecture.

```
companion connected → RequestApproval() → wait for response
companion absent    → log warning → proceed
```

### 3. Deferred Agent Wiring via SetAgent()

**Decision**: Create Gateway before Agent, then wire Agent via `SetAgent()`.

**Rationale**: The tool approval wrapper needs a reference to Gateway (for `HasCompanions()` and `RequestApproval()`). Since tools are created before the agent, and tools need Gateway, Gateway must be created first. But Gateway also needs Agent for chat handling. Solution: create Gateway with nil agent, wrap tools with approval, create Agent, then call `SetAgent()`.

```
Gateway(nil agent) → wrapWithApproval(tools, gateway) → Agent(tools) → gateway.SetAgent(agent)
```

### 4. Session TTL Check in Get() Only

**Decision**: Check TTL expiration only in `Get()`, not in `AppendMessage()` or `Update()`.

**Rationale**: TTL is a read-side concern. Writing to an expired session should still work (the session gets refreshed via `UpdatedAt`). Checking on read prevents stale sessions from being served while allowing in-flight writes to complete naturally.

### 5. Lazy Config Loading for Security CLI

**Decision**: Changed `NewSecurityCmd(cfg *config.Config)` to accept `func() (*config.Config, error)` closure.

**Rationale**: The security command is registered at startup before config is loaded. Other commands (serve, config validate) load config in their own RunE. The closure pattern defers config loading to command execution time while keeping the registration clean.

### 6. Companion Discovery Deletion (Not Deprecation)

**Decision**: Delete `internal/companion/` entirely rather than deprecating.

**Rationale**: The package has no imports anywhere in the codebase, contains multiple TODO stubs, and the Gateway already handles companion WebSocket connections directly. The mDNS discovery was a prototype that was superseded by the direct WebSocket approach. Keeping dead code increases maintenance burden and confuses audits.

## Data Flow Changes

### Before (Phantom)
```
Config → [temperature=0.7] → ProviderProxy → Generate(params{Temperature:0}) → LLM
         ^^ ignored                           ^^ always zero
```

### After (Wired)
```
Config → [temperature=0.7] → ProviderProxy(WithTemperature(0.7))
                                → Generate(params) → if params.Temperature==0, use 0.7 → LLM
```

### Tool Registration Before/After
```
Before:
  exec + filesystem + browser → Agent

After:
  exec + filesystem + browser + crypto(5) + secrets(4) → [approval wrap] → Agent
                                                            ↓
                                                     if ApprovalRequired &&
                                                     tool in SensitiveTools
```
