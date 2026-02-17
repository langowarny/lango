## Tasks

### Phase 1: P0 Bug Fixes

- [x] Fix BoostLearningConfidence signature in `internal/knowledge/store.go` — add `confidenceBoost float64` parameter
- [x] Fix confidence propagation in `internal/learning/graph_engine.go:140-148` — apply `0.1 * propagation` (0.03)
- [x] Update base engine call site in `internal/learning/engine.go:135` — pass `0.0` for confidenceBoost
- [x] Update `internal/learning/engine_test.go` — reflect new signature
- [x] Add warn logging + atomic drop counter to `internal/embedding/buffer.go`
- [x] Add warn logging + atomic drop counter to `internal/graph/buffer.go`

### Phase 2: P1 Conversation Analyzer

- [x] Create `internal/learning/conversation_analyzer.go` — LLM-based knowledge extraction
- [x] Create `internal/learning/session_learner.go` — session-end high-confidence learning
- [x] Create `internal/learning/analysis_buffer.go` — async analysis buffer
- [x] Create `internal/learning/parse.go` — LLM JSON response parsing + category mapping
- [x] Create `internal/learning/token.go` — token estimation (mirror memory/token.go)
- [x] Add config fields to `internal/config/types.go` — AnalysisTurnThreshold, AnalysisTokenThreshold
- [x] Wire conversation analysis in `internal/app/wiring.go` — initConversationAnalysis()
- [x] Add analysis buffer lifecycle to `internal/app/app.go` — Start/Stop
- [x] Add fields to `internal/app/types.go` — AnalysisBuffer
- [x] Create `internal/learning/conversation_analyzer_test.go`
- [x] Create `internal/learning/session_learner_test.go`
- [x] Create `internal/learning/analysis_buffer_test.go`

### Phase 3: P2 RAG Quality

- [x] Add MaxDistance to `internal/config/types.go` RAGConfig
- [x] Add SearchOptions and update VectorStore interface in `internal/embedding/store.go`
- [x] Implement metadata post-filter in `internal/embedding/sqlite_vec.go`
- [x] Add MaxDistance filtering + session key pass-through in `internal/embedding/rag.go`
- [x] Add MaxDistance to VectorRetrieveOptions in `internal/graph/rag.go`
- [x] Inject session key in `internal/adk/context_model.go` assembleRAGSection/assembleGraphRAGSection
- [x] Pass MaxDistance in ragOpts in `internal/app/wiring.go`

### Verification

- [x] `go build ./...`
- [x] `go test ./...`
