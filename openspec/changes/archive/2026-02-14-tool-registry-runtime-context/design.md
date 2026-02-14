## Context

The Context Retriever implements a 6-layer architecture (plus 2 memory layers handled separately) for retrieval-augmented generation. Currently, 4 layers are implemented (UserKnowledge, SkillPatterns, ExternalKnowledge, AgentLearnings), while Tool Registry and Runtime Context are skipped in the switch statement with `continue // handled elsewhere` — but no code handles them anywhere.

The retriever lives in `internal/knowledge` (no external dependencies), adapters live in `internal/adk` (depends on knowledge + agent), and wiring happens in `internal/app`.

## Goals / Non-Goals

**Goals:**
- Implement Tool Registry layer: allow the LLM to see which tools match the current query
- Implement Runtime Context layer: expose session key, channel type, tool count, and system feature flags
- Maintain backward compatibility: default layer list (nil) stays as the original 4
- Follow existing patterns: provider interfaces in knowledge, adapters in adk, wiring in app

**Non-Goals:**
- Semantic/vector search for tools — substring matching is sufficient given typical tool counts (<50)
- Dynamic tool registration at runtime — tool list is fixed at agent initialization
- Persisting runtime context — it is ephemeral session state

## Decisions

### Interface location: knowledge package
Provider interfaces (`ToolRegistryProvider`, `RuntimeContextProvider`) are defined in the knowledge package alongside the retriever. This avoids circular dependencies since knowledge has no dependency on agent or adk.

Alternative considered: defining interfaces in adk. Rejected because the retriever in knowledge would then need to import adk, creating a circular dependency.

### Builder pattern for optional providers
`WithToolRegistry()` and `WithRuntimeContext()` builder methods on `ContextRetriever` keep the constructor signature stable. Existing callers are unaffected.

Alternative considered: adding providers to the constructor. Rejected because it would break all existing call sites and the providers are truly optional.

### Boundary copy in ToolRegistryAdapter
The adapter copies the input `[]*agent.Tool` slice at construction time to prevent external mutation. This follows the project's Go guidelines on copying slices at boundaries.

### sync.RWMutex for RuntimeContextAdapter
`SetSession()` may be called from the gateway goroutine while `GetRuntimeContext()` is called from the retriever. RWMutex allows concurrent reads with exclusive writes.

### Extended layers requested only from ContextAwareModelAdapter
The default layer list (when `Layers` is nil) remains the original 4. Only `ContextAwareModelAdapter.GenerateContent()` explicitly requests all 6 layers. This ensures backward compatibility for any code using `Retrieve()` directly.

## Risks / Trade-offs

- [Substring matching may return irrelevant tools for short queries] → Acceptable given small tool sets; limit parameter caps results
- [Runtime context adds a fixed item regardless of query relevance] → Single item overhead is negligible; provides consistent session awareness
- [SetSession called per LLM request] → Lightweight operation (string assignment + mutex); no performance concern
