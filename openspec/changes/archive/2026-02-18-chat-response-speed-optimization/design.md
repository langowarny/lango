## Context

The Lango chat pipeline executes context retrieval (knowledge, RAG, memory) sequentially before every LLM call, adding 650-2700ms of pre-LLM latency. The RAG service also searches 4 collections sequentially (120-600ms) and calls the embedding API on every request (500-2000ms). Additionally, memory context grows unbounded as sessions lengthen, inflating prompt size and LLM processing time. The gateway waits for the complete response before sending anything to the UI.

## Goals / Non-Goals

**Goals:**
- Reduce pre-LLM context retrieval latency by 40-60% via parallelization
- Eliminate redundant embedding API calls via TTL-based caching
- Bound memory context size to prevent prompt inflation in long sessions
- Stream LLM tokens to UI in real-time for improved perceived latency

**Non-Goals:**
- LLM model optimization or provider switching
- Streaming support for channel adapters (Telegram, Discord, Slack) — they use separate `App.runAgent` path
- Persistent embedding cache (disk/Redis) — in-memory is sufficient at current scale
- Hallucination retry in streaming mode — kept simple for initial implementation

## Decisions

### 1. errgroup for parallelization
**Decision**: Use `golang.org/x/sync/errgroup` for both RAG collection search and context retrieval.
**Rationale**: Already an indirect dependency, well-tested, provides context cancellation propagation. Alternative (raw goroutines + WaitGroup) requires more boilerplate and manual error handling.

### 2. Pre-allocated index-based result collection
**Decision**: For parallel RAG search, use `perColResults[i]` indexed by collection position.
**Rationale**: Each goroutine writes to its own index, eliminating mutex contention. Alternative (shared slice with mutex) adds unnecessary synchronization overhead.

### 3. In-memory TTL cache for embeddings
**Decision**: Simple `sync.RWMutex` + `map[string]embeddingCacheEntry` with 5-minute TTL and 100-entry max.
**Rationale**: Query embedding vectors are ~1.5KB each, so 100 entries ≈ 150KB memory. TTL ensures freshness. Alternative (LRU with external library) adds dependency for minimal benefit at this scale.

### 4. Non-fatal error pattern in parallel retrieval
**Decision**: Each parallel goroutine logs errors internally and returns nil to errgroup.
**Rationale**: Preserves existing degradation pattern — a failing knowledge retrieval should not block RAG or memory. Context sections are best-effort.

### 5. RunStreaming as additive API
**Decision**: Add `RunStreaming()` alongside existing `RunAndCollect()` rather than replacing it.
**Rationale**: Channel adapters and automation systems rely on `RunAndCollect`. `RunStreaming` is specifically for the gateway WebSocket path. Hallucination retry remains in `RunAndCollect` only.

### 6. Memory limits with sensible defaults
**Decision**: Default 5 reflections, 20 observations in context. Zero means unlimited (backward compatible).
**Rationale**: Reflections are summaries (few needed), observations are granular (more needed for recent context). Configurable via `observationalMemory` config section.

## Risks / Trade-offs

- **[Parallel search correctness]** → Each goroutine writes to pre-allocated index; no shared mutable state. Go 1.25 loop variable capture eliminates closure bugs.
- **[Cache staleness]** → 5-minute TTL is conservative. Embedding models don't change frequently. Cache miss penalty is the same as current behavior.
- **[Memory limit information loss]** → Oldest reflections/observations are dropped. Mitigated by retrieving most recent entries (DESC + reverse for chronological display).
- **[Streaming partial state]** → If connection drops mid-stream, client receives partial response. Mitigated by RPC response still containing full text.
