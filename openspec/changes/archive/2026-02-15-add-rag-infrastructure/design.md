# Design: RAG Infrastructure with sqlite-vec

## Architecture

```
┌──────────────────────────────────────────────────────────┐
│ Agent (ContextAwareModelAdapter)                         │
│  ├── Knowledge Retriever (keyword, existing)             │
│  ├── RAG Service (semantic, NEW)                         │
│  └── Memory Provider (session, existing)                 │
│           ↓                                              │
│  System Prompt Assembly → LLM                            │
└──────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│ Async Embedding Pipeline                                 │
│  Store.Save() → EmbedCallback → EmbeddingBuffer.Enqueue()│
│                                       ↓ (batch)         │
│                            EmbeddingProvider.Embed()     │
│                                       ↓                  │
│                            VectorStore.Upsert()          │
│                            (sqlite-vec vec0 table)       │
└──────────────────────────────────────────────────────────┘
```

## Key Design Decisions

### 1. Callback-based Store Integration (no circular imports)
- `embedding` package imports `knowledge` and `memory` (via resolver)
- Stores define their own `EmbedCallback` func type and accept it via `SetEmbedCallback`
- Wiring connects them at app init time in `initEmbedding()`

### 2. Shared SQLite Database
- vec0 virtual table lives in the same DB file as ent-managed tables
- `bootstrap.Result.RawDB` exposes the `*sql.DB` handle for direct SQL
- sqlite-vec extension auto-loaded via `sqlite_vec.Auto()` in init()

### 3. Provider Registry with Fallback
- Single `Registry` manages primary + fallback provider
- Factory creates providers from `ProviderConfig` structs
- Local provider reuses `sashabaranov/go-openai` SDK with custom BaseURL

### 4. Batched Async Processing
- `EmbeddingBuffer` follows `memory.Buffer` lifecycle: Start/Stop/WaitGroup
- Queue size 256, batch size 32, flush every 2s
- Non-blocking enqueue with drop-on-full semantics

### 5. RAG Integration Point
- `ContextAwareModelAdapter.WithRAG()` adds RAG to the existing adapter chain
- RAG section injected after knowledge retrieval, before memory section
- Disabled by default; requires `embedding.rag.enabled: true` in config

## Package Structure

```
internal/embedding/
├── provider.go      # EmbeddingProvider interface
├── store.go         # VectorStore interface + types
├── openai.go        # OpenAI provider
├── google.go        # Google provider
├── local.go         # Local/Ollama provider
├── registry.go      # Provider factory + registry
├── sqlite_vec.go    # sqlite-vec VectorStore implementation
├── buffer.go        # Async embedding buffer
├── rag.go           # RAG retrieval service
├── resolver.go      # ContentResolver (knowledge + memory stores)
└── *_test.go        # Tests
```

## Modified Packages

- `config` — `EmbeddingConfig`, `LocalEmbeddingConfig`, `RAGConfig`
- `bootstrap` — `Result.RawDB` field
- `app` — `initEmbedding()` wiring, Start/Stop integration
- `knowledge` — `EmbedCallback`, `SetEmbedCallback`, `GetLearning`
- `memory` — `EmbedCallback`, `SetEmbedCallback`, `GetObservation`, `GetReflection`
- `adk` — `WithRAG()`, `assembleRAGSection()`
- `cli/doctor/checks` — `EmbeddingCheck`
