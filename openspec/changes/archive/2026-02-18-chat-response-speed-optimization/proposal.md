## Why

Chat response latency is significantly impacted by sequential context retrieval before LLM calls, lack of response streaming, and unbounded memory context that grows with session length. Parallelizing retrieval, caching embeddings, limiting memory context, and enabling WebSocket streaming can reduce overall response delay by 40-60% and dramatically improve perceived latency.

## What Changes

- Parallel RAG collection search using errgroup (4 collections searched concurrently)
- TTL-based in-memory embedding cache to skip redundant embedding API calls
- Configurable limits on reflections/observations injected into LLM context
- Parallel context retrieval (knowledge, RAG, memory run concurrently via errgroup)
- WebSocket streaming of LLM tokens via `agent.chunk` events for real-time UI updates

## Capabilities

### New Capabilities
- `embedding-query-cache`: TTL-based in-memory cache for query embedding vectors to avoid redundant embedding API calls
- `websocket-streaming`: Real-time LLM token streaming through WebSocket `agent.chunk` events

### Modified Capabilities
- `embedding-rag`: Parallel collection search and embedding cache integration
- `observational-memory`: Configurable limits on reflections/observations in LLM context, new ListRecent store methods
- `gateway-server`: Streaming agent responses via RunStreaming instead of RunAndCollect
- `context-retriever`: Parallel knowledge/RAG/memory retrieval using errgroup

## Impact

- **Code**: `internal/embedding/rag.go`, `internal/embedding/cache.go` (new), `internal/adk/context_model.go`, `internal/adk/agent.go`, `internal/gateway/server.go`, `internal/memory/store.go`, `internal/config/types.go`, `internal/app/wiring.go`
- **Dependencies**: `golang.org/x/sync` promoted from indirect to direct
- **APIs**: New WebSocket event `agent.chunk` (backward-compatible, existing events preserved)
- **Config**: New fields `observationalMemory.maxReflectionsInContext` and `observationalMemory.maxObservationsInContext`
