## Context

Lango is a Go-based AI assistant (ADK v0.4.0) with knowledge, memory, embedding/RAG, and security subsystems. Prior to this change, the architecture was a single monolithic agent with flat keyword-based knowledge retrieval and linear memory. The system lacked relationship-aware reasoning, multi-agent specialization, and external agent communication.

Key constraints:
- Go 1.25.4 (Cayley incompatible — direct BoltDB implementation required)
- ADK v0.4.0 sub-agent mechanism (llmagent, remoteagent packages)
- Import cycle avoidance via callback types and mirror types
- Existing async buffer pattern (Start/Enqueue/Stop lifecycle)

## Goals / Non-Goals

**Goals:**
- Add a graph store (BoltDB-backed triple store) for relationship tracking across knowledge, memory, and learning
- Implement 2-phase Graph RAG: vector search (sqlite-vec) → graph expansion (BFS traversal)
- Enable multi-agent orchestration: 4 specialized sub-agents (Executor, Researcher, Planner, MemoryManager)
- Expose agent as A2A-compatible endpoint (Agent Card at `/.well-known/agent.json`)
- Connect remote A2A agents as sub-agents in the orchestrator
- Wire observational memory graph hooks for temporal/session triples
- Wire self-learning graph engine with confidence propagation

**Non-Goals:**
- Full RDF/SPARQL query support (BoltDB triple store only)
- Agent-to-agent task delegation with streaming (future)
- Graph visualization UI
- Distributed graph store (single-node BoltDB only)

## Decisions

### D1: BoltDB instead of Cayley for graph store
**Decision**: Direct BoltDB implementation with SPO/POS/OSP index buckets.
**Rationale**: Cayley requires Go < 1.22. Direct BoltDB gives full control, zero external dependencies beyond bbolt, and sufficient performance for single-node graph operations.
**Alternatives**: (a) Cayley — incompatible with Go 1.25; (b) RocksDB — heavier, overkill for embedded use; (c) BadgerDB — more complex API for same use case.

### D2: ToolResultObserver interface for learning polymorphism
**Decision**: Extract `ToolResultObserver` interface from `learning.Engine`, implement on both `Engine` and `GraphEngine`.
**Rationale**: `wrapWithLearning` in tools.go previously took `*learning.Engine` directly, preventing `GraphEngine.OnToolResult` override from being called. An interface enables polymorphic dispatch without import cycles.
**Alternatives**: (a) Type-switch in wrapWithLearning — fragile, violates OCP; (b) Callback function — loses method grouping.

### D3: Callback pattern for cross-package graph wiring
**Decision**: Use `TripleCallback func([]graph.Triple)` in memory/learning packages, with buffer.Enqueue as the concrete callback.
**Rationale**: Avoids import cycles (memory→graph→memory). Consistent with existing `EmbedCallback` pattern.

### D4: ADK native sub-agents via llmagent
**Decision**: Use ADK's `llmagent.New()` with `SubAgents` field for orchestration hierarchy.
**Rationale**: Native ADK mechanism handles tool routing, session delegation, and agent tree traversal. No custom orchestration protocol needed.

### D5: Graph initialization before Knowledge
**Decision**: Reorder init: `gc := initGraphStore()` → `kc := initKnowledge(cfg, store, tools, gc)`.
**Rationale**: `GraphEngine` needs `gc.store` at creation time. Previous order (kc → gc) made it impossible to wire GraphEngine during knowledge initialization.

### D6: Entity extraction as async goroutine in graph callback
**Decision**: Launch goroutine in `wireGraphCallbacks` graphCB for LLM-based entity extraction.
**Rationale**: Entity extraction is expensive (LLM call). Blocking the save path would degrade user experience. The existing `GraphBuffer` already handles async batching.

## Risks / Trade-offs

**[Risk] Entity extraction goroutine leak** → Mitigated by GraphBuffer.Stop() draining queue before shutdown. Extraction goroutines are fire-and-forget but bounded by buffer queue capacity (256).

**[Risk] Memory GraphHooks ListObservations in SaveReflection** → On each reflection save, all session observations are queried for graph linking. For sessions with many observations, this could be slow. → Acceptable for current scale; can add caching or limit if needed.

**[Risk] Remote A2A agent load failure** → Non-fatal: logged as warning, orchestrator continues with local sub-agents only. No degradation of core functionality.

**[Trade-off] Graph tools only available in multi-agent mode** → Graph tools (graph_traverse, graph_query) are always added when graph is enabled, regardless of multi-agent mode. This allows single-agent mode to also use graph tools.

**[Trade-off] lastObsIDs map grows per session** → In-memory map tracking last observation per session key. Unbounded but each entry is just two strings. Acceptable for expected session counts.
