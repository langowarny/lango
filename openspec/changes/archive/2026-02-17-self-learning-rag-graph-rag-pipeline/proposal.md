## Why

Conversations are stored in the DB but self-learning, RAG, and Graph RAG pipelines are not functioning effectively. The confidence propagation bug (`int(0.3) = 0`) makes all boosts +1, buffer drops are silent, there is no conversation-level knowledge extraction, RAG returns irrelevant results without distance filtering, and session context is not scoped. Config defaults (false) remain unchanged, but when enabled these features must work end-to-end.

## What Changes

- **Fix confidence propagation bug**: `BoostLearningConfidence` signature change to accept `float64` confidence boost directly, fixing `int(0.3)=0` truncation
- **Fix buffer silent drops**: Change `Debugw` to `Warnw` logging with atomic drop counters on `EmbeddingBuffer` and `GraphBuffer`
- **New Conversation Analyzer**: LLM-based extraction of facts, patterns, corrections, and preferences from conversation turns (every N turns or token threshold)
- **New Session Learner**: End-of-session analysis producing high-confidence knowledge entries with cross-reference graph triples
- **New Analysis Buffer**: Async processing buffer for conversation analysis with Start/Trigger/Stop lifecycle
- **RAG MaxDistance threshold**: Filter out low-relevance vector search results by cosine distance
- **RAG session-scoped filtering**: Post-filter vector search results by metadata (session_key) for session-aware retrieval
- **Context model session key injection**: Pass session key through RAG and Graph RAG retrieval for scoped results

## Capabilities

### New Capabilities
- `conversation-analysis`: LLM-based knowledge extraction from conversation turns and session-end analysis

### Modified Capabilities
- `learning-engine`: BoostLearningConfidence signature change for proper float64 confidence propagation
- `embedding-rag`: Add MaxDistance threshold filtering and session-scoped metadata filtering
- `graph-rag`: Extend VectorRetrieveOptions with MaxDistance support
- `observational-memory`: Buffer drop warning logging with counters (EmbeddingBuffer, GraphBuffer)

## Impact

- **Core**: `internal/knowledge/store.go` (API signature change), `internal/learning/` (bug fix + new files), `internal/embedding/` (store interface + RAG), `internal/graph/` (RAG options)
- **Application**: `internal/app/wiring.go` (new init function), `internal/app/app.go` (buffer lifecycle), `internal/adk/context_model.go` (session key injection)
- **Config**: `internal/config/types.go` (new fields in KnowledgeConfig and RAGConfig)
- **Tests**: New test files for conversation analyzer, session learner, analysis buffer; updates to existing engine tests
