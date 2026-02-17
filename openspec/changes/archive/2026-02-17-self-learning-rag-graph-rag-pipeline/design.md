## Context

Lango's self-learning pipeline has structural components in place (knowledge store, learning engine, graph engine, embedding buffer, graph buffer, RAG service, Graph RAG service) but critical bugs and missing analysis logic prevent the system from actually learning from conversations. The confidence propagation bug (`int(0.3) = 0`) makes all boosts fixed at +1, buffer drops are silent, and there is no mechanism to extract knowledge from conversation content — only tool errors are tracked.

## Goals / Non-Goals

**Goals:**
- Fix confidence propagation to use proper float64 math
- Make buffer drops visible via warn-level logging with counters
- Enable LLM-based knowledge extraction from conversations (facts, patterns, corrections, preferences)
- Enable session-end learning for high-confidence cross-session knowledge
- Improve RAG quality with distance thresholds and session-scoped filtering
- Wire everything into the application lifecycle

**Non-Goals:**
- Changing config defaults (they remain false)
- Phase 4 (Graph RAG implicit feedback, A2A graph integration) — deferred
- Changing the embedding provider or vector store implementation
- Adding new CLI commands for analysis management

## Decisions

### 1. BoostLearningConfidence signature: add `confidenceBoost float64`
**Choice**: Add a `confidenceBoost float64` parameter; when > 0, add it directly to confidence and clamp [0.1, 1.0]. When 0, use existing success/occurrence ratio.
**Rationale**: The graph engine needs fractional confidence propagation (0.03 = 0.1 * 0.3). The base engine can pass 0.0 to keep existing behavior. This is minimally invasive.

### 2. ConversationAnalyzer uses its own TextGenerator interface
**Choice**: Declare `TextGenerator` in the learning package (same signature as `memory.TextGenerator`) rather than importing from memory.
**Rationale**: Avoids import cycles between learning and memory packages. The interface is trivial (one method).

### 3. Analysis triggering: turn count OR token threshold
**Choice**: AnalysisBuffer fires when either `turnThreshold` (default 10) or `tokenThreshold` (default 2000) is exceeded since last analysis for that session.
**Rationale**: Short sessions with dense content should still trigger analysis. Token threshold catches this.

### 4. Session learner: sampling for long sessions
**Choice**: For sessions > 20 turns, sample first 3 + every 5th + last 5 messages instead of full history.
**Rationale**: LLM context limits and cost. The sampling strategy captures session start/evolution/end patterns.

### 5. RAG metadata filtering: post-filter approach
**Choice**: Over-fetch by 3x from sqlite-vec, then post-filter by metadata in Go.
**Rationale**: sqlite-vec virtual tables don't support WHERE clauses on metadata columns in the MATCH query. Post-filtering is simple and correct. The 3x multiplier ensures enough results survive filtering.

### 6. SearchOptions as optional parameter
**Choice**: Add `*SearchOptions` parameter to `VectorStore.Search()`. Nil means no filtering (backward compatible).
**Rationale**: All existing callers pass nil. Only RAG service passes session-scoped options. No breaking change for tests.

## Risks / Trade-offs

- **LLM cost for conversation analysis**: Each analysis call consumes tokens. Mitigated by turn/token thresholds and session-end sampling.
- **BoostLearningConfidence signature change**: Breaking change for callers. Mitigated by updating all 2 call sites simultaneously.
- **Post-filter may reduce result count**: If most results don't match metadata filter, fewer results returned. Mitigated by 3x over-fetch.
- **Analysis buffer backpressure**: Queue full drops analysis requests. Acceptable — analysis is best-effort, not critical path.
