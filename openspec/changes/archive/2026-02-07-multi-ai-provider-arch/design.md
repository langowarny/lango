## Context

Lango currently uses Google's ADK SDK directly for Gemini integration. The `Config` struct in `internal/agent/runtime.go` has single-provider fields (`Provider`, `Model`, `APIKey`). The unused `ProviderAdapter` interface and `providers` map indicate planned extensibility.

OpenClaw analysis revealed patterns for multi-provider support: auth profile rotation, model fallback, and OpenAI-compatible abstraction that works with 10+ services.

## Goals / Non-Goals

**Goals:**
- Unified `Provider` interface supporting streaming, tools, and model listing
- OpenAI-compatible provider for OpenAI, Ollama, Groq, Together AI, etc.
- Anthropic provider for Claude models
- Extended `lango.json` with `providers` configuration section
- Model fallback for resilience

**Non-Goals:**
- Auth profile rotation (OpenClaw complexity, keep simple initially)
- AWS Bedrock (high complexity, defer to future)
- CLI providers (Claude CLI, Codex CLI)
- Model discovery/catalog UI

## Decisions

### 1. Provider Interface Design
**Decision**: Use streaming iterator pattern with `iter.Seq2[StreamEvent, error]`

**Rationale**: Go 1.23+ iterators provide clean streaming abstraction. Matches existing `StreamEvent` type in runtime.go. Allows lazy evaluation and early termination.

**Alternatives considered**:
- Callback-based streaming: More complex error handling
- Channel-based: Harder to compose and integrate

### 2. Implementation Priority
**Decision**: OpenAI-Compatible first, then Anthropic, then native Ollama

**Rationale**: OpenAI-Compatible covers 10+ services with one implementation. Ollama already works via OpenAI-compatible endpoint. Anthropic needs native SDK for Extended Thinking.

### 3. Package Structure
**Decision**: `internal/provider/` package with sub-packages per provider

```
internal/provider/
├── provider.go          # Interface, Registry
├── openai/openai.go     # OpenAI-compatible
├── anthropic/anthropic.go
└── gemini/gemini.go     # Refactored from current code
```

**Rationale**: Clean separation, testable, follows existing Lango patterns.

### 4. Configuration Schema
**Decision**: Add `providers` map with per-provider settings

```json
{
  "providers": {
    "openai": { "apiKey": "${OPENAI_API_KEY}" },
    "ollama": { "baseUrl": "http://localhost:11434" }
  },
  "agent": {
    "provider": "openai",
    "model": "gpt-4o",
    "fallbacks": [{"provider": "gemini", "model": "gemini-2.0-flash"}]
  }
}
```

**Rationale**: Backwards compatible, environment variable substitution already exists in config-system.

### 5. Gemini Preservation
**Decision**: Keep current Gemini code, wrap in Provider interface

**Rationale**: Minimal disruption, maintains existing functionality, allows gradual migration.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| SDK version conflicts | Pin specific versions in go.mod |
| OpenAI-compatible endpoint variations | Graceful fallback, document tested services |
| Streaming API differences | Normalize in provider implementations |
| Breaking existing configs | Full backwards compatibility, existing `agent.provider: gemini` works unchanged |

## Open Questions

1. Should we support multiple API keys per provider for rotation? (Defer to future)
2. How to handle provider-specific parameters (e.g., Anthropic Extended Thinking)? (Pass via metadata map)
