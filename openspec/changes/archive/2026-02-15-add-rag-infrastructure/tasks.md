# Tasks: RAG Infrastructure with sqlite-vec

## Phase 1: Core Interfaces
- [x] 1.1 Create `internal/embedding/provider.go` — EmbeddingProvider interface
- [x] 1.2 Create `internal/embedding/store.go` — VectorStore interface + VectorRecord/SearchResult types

## Phase 2: Provider Implementations
- [x] 2.1 Create `internal/embedding/openai.go` — OpenAI embedding provider (text-embedding-3-small)
- [x] 2.2 Create `internal/embedding/google.go` — Google embedding provider (text-embedding-004)
- [x] 2.3 Create `internal/embedding/local.go` — Local/Ollama embedding provider via OpenAI-compat API
- [x] 2.4 Create `internal/embedding/registry.go` — Provider factory with fallback support

## Phase 3: sqlite-vec VectorStore
- [x] 3.1 Add `sqlite-vec-go-bindings/cgo` dependency to go.mod
- [x] 3.2 Create `internal/embedding/sqlite_vec.go` — vec0 virtual table init, Upsert, Search, Delete

## Phase 4: Async Embedding Pipeline
- [x] 4.1 Create `internal/embedding/buffer.go` — EmbeddingBuffer with batched background processing

## Phase 5: RAG Service
- [x] 5.1 Create `internal/embedding/rag.go` — RAGService with multi-collection semantic search
- [x] 5.2 Create `internal/embedding/resolver.go` — ContentResolver for knowledge/memory stores

## Phase 6: Integration
- [x] 6.1 Add `EmbeddingConfig`/`RAGConfig` to `internal/config/types.go`
- [x] 6.2 Add `RawDB` field to `internal/bootstrap/bootstrap.go` Result
- [x] 6.3 Add `EmbeddingBuffer`/`RAGService` fields to `internal/app/types.go`
- [x] 6.4 Create `initEmbedding()` in `internal/app/wiring.go` with store callback wiring
- [x] 6.5 Wire embedding Start/Stop in `internal/app/app.go`
- [x] 6.6 Add `EmbedCallback` + `SetEmbedCallback` to `internal/knowledge/store.go`
- [x] 6.7 Add `EmbedCallback` + `SetEmbedCallback` to `internal/memory/store.go`
- [x] 6.8 Add `GetLearning()` to knowledge.Store, `GetObservation()`/`GetReflection()` to memory.Store
- [x] 6.9 Add `WithRAG()` and `assembleRAGSection()` to `internal/adk/context_model.go`
- [x] 6.10 Wire RAG into `initAgent()` in wiring.go

## Phase 7: Doctor Check
- [x] 7.1 Create `internal/cli/doctor/checks/embedding.go` — EmbeddingCheck
- [x] 7.2 Register EmbeddingCheck in AllChecks()

## Phase 8: Tests
- [x] 8.1 Create `sqlite_vec_test.go` — VectorStore CRUD and similarity search tests
- [x] 8.2 Create `buffer_test.go` — Async processing and graceful shutdown tests
- [x] 8.3 Create `rag_test.go` — RAG retrieval, filtered collection, empty query tests
- [x] 8.4 Create `registry_test.go` — Provider creation, unknown type, fallback tests

## Verification
- [x] `go build ./...` passes
- [x] `go test ./...` passes (all 13 new embedding tests + all existing tests)
