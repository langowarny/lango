## Context

The current `agent-runtime` implementation in `internal/agent/runtime.go` is a 284-line placeholder that defines interfaces and structures but lacks actual LLM integration. The `generate()` function returns hardcoded placeholder text instead of calling real LLM providers. The spec requires ADK framework usage, but no ADK integration exists in the codebase.

**Current state:**
- Placeholder runtime with defined interfaces (`Runtime`, `Tool`, `ProviderAdapter`)
- Tool registration system without execution integration
- Session management that works but can't generate real responses
- Provider abstraction without concrete implementations

**Constraints:**
- Must maintain existing external API (channels, gateway continue to work)
- Must integrate with existing session store (SQLite via Ent)
- Must support existing tool infrastructure
- Go 1.22+, prefer standard library where possible

## Goals / Non-Goals

**Goals:**
- Integrate ADK (`google.golang.org/adk`) framework for LLM agent runtime
- Use ADK's `llmagent` for LLM-driven agent creation
- Support Gemini models natively (via `google.golang.org/genai`)
- Enable real tool calling with existing tool handlers via ADK's `tool.Tool` interface
- Implement streaming responses using ADK's execution flow
- Maintain session context across conversation turns

**Non-Goals:**
- Multi-agent orchestration (Sequential, Parallel, Loop agents) - just single LlmAgent for now
- Advanced agentic workflows - focus on basic chat + tools
- Embeddings, RAG, or vector search
- Custom model providers beyond Gemini (extensible for future)
- Web UI or ADK launcher integration (use existing channels)

## Decisions

### D1: ADK Integration Approach

**Decision:** Use ADK's `llmagent.New()` to create an LLM agent and wrap it in our existing `Runtime` interface.

**Rationale:**
- ADK provides `llmagent.Config` with Model, Instructions, and Tools
- Clean separation: ADK handles LLM interaction, we handle session/channels
- Allows keeping our `Runtime` struct as the primary interface
- ADK's tool system maps well to our `Tool` type

**Implementation:**
```go
// Wrap ADK llmagent in our Runtime
type Runtime struct {
    config       Config
    adkAgent     agent.Agent  // ADK llmagent instance
    tools        map[string]*Tool
    sessionStore session.Store
}
```

**Alternatives considered:**
- Direct Gemini SDK: Miss out on ADK's agent orchestration features
- Custom LLM wrapper: Reinventing what ADK already provides

### D2: Tool Integration Strategy

**Decision:** Convert our `Tool` structs to ADK `tool.Tool` interface implementations.

**Rationale:**
- ADK expects `tool.Tool` interface with `Name()`, `Description()`, `Run(ctx, input)`
- Our tools have `Tool{Name, Description, Parameters, Handler}`
- Create adapter type that implements ADK's interface

**Mapping:**
```go
type AdkToolAdapter struct {
    tool *Tool
}

func (a *AdkToolAdapter) Name() string { return a.tool.Name }
func (a *AdkToolAdapter) Description() string { return a.tool.Description }
func (a *AdkToolAdapter) Run(ctx context.Context, input any) (any, error) {
    params, _ := input.(map[string]interface{})
    return a.tool.Handler(ctx, params)
}
```

### D3: Model Selection and Configuration

**Decision:** Start with Gemini support, use ADK's model interface for extensibility.

**Models:**
- Gemini: `gemini-2.0-flash-exp` (fast, reliable, built-in ADK support)
- Future: Claude/OpenAI via custom model adapters implementing ADK's `model.Model` interface

**Configuration:**
```go
type Config struct {
    Provider    string  // "gemini" (anthropic, openai for future)
    Model       string  // "gemini-2.0-flash-exp"
    APIKey      string  // GOOGLE_API_KEY
    MaxTokens   int
    Temperature float64
}
```

**How ADK models work:**
```go
model, err := gemini.NewModel(ctx, "gemini-2.0-flash-exp", &genai.ClientConfig{
    APIKey: apiKey,
})
```

### D4: Streaming Implementation

**Decision:** Use ADK agent execution and stream responses via our existing `StreamEvent` channel.

**Event flow:**
```go
// ADK agent execution
resp, err := adkAgent.Execute(ctx, session)

// Convert ADK response to our StreamEvents
for _, part := range resp.Parts {
    switch part.Type {
    case "text":
        events <- StreamEvent{Type: "text_delta", Text: part.Text}
    case "tool_call":
        events <- StreamEvent{Type: "tool_start", ToolCall: convertToolCall(part)}
    }
}
```

**Note:** ADK's `Execute()` returns complete response. For true streaming, may need to explore ADK's event system or iterate on parts as they arrive.

### D5: Session Context Management

**Decision:** Load history from session store, convert to ADK session format, pass to agent.

**Strategy:**
```go
func (r *Runtime) buildAdkSession(sess *session.Session) *agent.Session {
    adkSession := agent.NewSession()
    
    // Add conversation history
    for _, msg := range sess.History {
        adkSession.AddMessage(agent.Message{
            Role: msg.Role,
            Content: []agent.Part{{Type: "text", Text: msg.Content}},
        })
    }
    
    return adkSession
}
```

**Truncation:**
- Keep message history within reasonable bounds (last 20 turns)
- ADK handles token limits internally per model
- Log warning when truncating

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| ADK is v0.4.0 (early stage), APIs may change | Pin specific version, prepare for migration path |
| Limited to Gemini initially | Design model interface for future extensibility |
| Streaming may not be true token-by-token | Acceptable for v1, can enhance with ADK events later |
| ADK session format differs from ours | Conversion layer keeps concerns separated |
| Tool execution errors break agent | Wrap tool calls in error handling, emit tool_end with errors |

## Migration Plan

**Phase 1: Add dependency**
1. ADK already added: `go get google.golang.org/adk@v0.4.0` ✓
2. Gemini SDK added automatically as ADK dependency

**Phase 2: Refactor runtime.go**
1. Keep existing `Runtime` struct and `New()` constructor
2. Add ADK agent creation in `New()`
3. Implement tool adapters (Tool → ADK tool.Tool)
4. Replace `generate()` with ADK agent execution
5. Convert ADK responses to`StreamEvent`s

**Phase 3: Testing**
1. Unit tests: tool conversion, session building
2. Integration test: real Gemini API call (env var gated)
3. Manual testing: Telegram channel + real conversations

**Rollback:**
- Placeholder code in git history
- Revert ADK dependency if needed
- No data migration (session store unchanged)

## Open Questions

- [x] Which ADK package to use?
  → **Resolved:** `google.golang.org/adk` (not firebase/genkit)

- [ ] Does ADK support true streaming or only batch responses?
  → **Need to check**: ADK agent.Execute() API docs

- [ ] Token limit configuration - does ADK handle this automatically?
  → **Assumption:** Yes, based on model config

- [ ] Max tool calls per turn?
  → **Decision:** 10 tool calls max (ADK may have built-in limits)
