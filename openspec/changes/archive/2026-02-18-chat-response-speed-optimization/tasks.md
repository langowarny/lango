## 1. Parallel RAG Collection Search

- [x] 1.1 Convert RAGService.Retrieve() collection loop to errgroup parallel execution with pre-allocated index-based result collection
- [x] 1.2 Add golang.org/x/sync/errgroup import to embedding/rag.go

## 2. Embedding Query Cache

- [x] 2.1 Create internal/embedding/cache.go with embeddingCache struct (sync.RWMutex + map, TTL, maxSize, eviction logic)
- [x] 2.2 Add cache field to RAGService and initialize in NewRAGService (5min TTL, 100 max entries)
- [x] 2.3 Wire cache lookup in Retrieve() before embedding API call

## 3. Memory Context Limits

- [x] 3.1 Add MaxReflectionsInContext and MaxObservationsInContext fields to ObservationalMemoryConfig in config/types.go
- [x] 3.2 Add ListRecentReflections method to memory Store (DESC + Limit + reverse for chronological order)
- [x] 3.3 Add ListRecentObservations method to memory Store (DESC + Limit + reverse for chronological order)
- [x] 3.4 Add ListRecentReflections and ListRecentObservations to MemoryProvider interface in adk/context_model.go
- [x] 3.5 Add maxReflections/maxObservations fields and WithMemoryLimits() builder to ContextAwareModelAdapter
- [x] 3.6 Modify assembleMemorySection() to use ListRecent methods when limits are set
- [x] 3.7 Wire WithMemoryLimits in both ctxAdapter paths in app/wiring.go initAgent()

## 4. Parallel Context Retrieval

- [x] 4.1 Convert GenerateContent() sequential knowledge/RAG/memory retrieval to errgroup parallel execution
- [x] 4.2 Combine results into prompt after g.Wait() with section concatenation

## 5. WebSocket Streaming

- [x] 5.1 Add ChunkCallback type and RunStreaming method to adk/agent.go
- [x] 5.2 Modify gateway handleChatMessage to use RunStreaming with agent.chunk broadcast callback
- [x] 5.3 Verify existing agent.thinking and agent.done events remain in place

## 6. Verification

- [x] 6.1 Run go build ./... and verify clean compilation
- [x] 6.2 Run go test ./internal/embedding/... and verify all pass
- [x] 6.3 Run go test ./internal/memory/... and verify all pass
- [x] 6.4 Run go test ./internal/adk/... and verify all pass
- [x] 6.5 Run go test ./internal/gateway/... and verify all pass
