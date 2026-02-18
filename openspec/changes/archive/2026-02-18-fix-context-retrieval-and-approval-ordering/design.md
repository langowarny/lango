## Context

The context retrieval system (`internal/knowledge/`) extracts keywords from user queries and searches the knowledge store via Ent-generated SQLite queries. When cron jobs process long or complex messages, `extractKeywords()` generates unbounded keyword lists that get concatenated into a single `LIKE '%very long query%'` pattern, exceeding SQLite's pattern complexity limit.

Separately, the Telegram approval flow (`internal/channels/telegram/approval.go`) performs a synchronous `editMessage` API call before unblocking the waiting agent goroutine via channel send. This causes cumulative latency when multiple tools require consecutive approval, and Telegram's non-guaranteed message delivery order can cause the approval status edit to appear after the agent's final response.

## Goals / Non-Goals

**Goals:**
- Eliminate "LIKE or GLOB pattern too complex" errors during context retrieval
- Improve search relevance by matching individual keywords independently (OR semantics)
- Reduce agent pipeline latency during tool approval by removing synchronous Telegram API dependency
- Fix approval message ordering so status edits appear before agent responses

**Non-Goals:**
- Full-text search engine (e.g., FTS5) integration
- Approval message ordering guarantees at the Telegram API level (best-effort)
- Changes to keyword extraction algorithm quality (e.g., stemming, NLP)

## Decisions

1. **Keyword count limit (5) and length limit (50 chars)** — SQLite LIKE patterns have O(n*m) complexity. Capping at 5 short keywords ensures the total pattern count stays well within limits. 5 keywords provide sufficient recall for context retrieval without over-matching. Alternative: dynamic limit based on query length — rejected for unnecessary complexity.

2. **Per-keyword OR predicates instead of single concatenated pattern** — Splitting `LIKE '%kw1 kw2 kw3%'` into `(content LIKE '%kw1%' OR key LIKE '%kw1%') OR (content LIKE '%kw2%' OR key LIKE '%kw2%')` produces shorter individual patterns and enables order-independent matching, which is semantically more correct. Alternative: use SQLite FTS5 — rejected as over-engineering for current scale.

3. **Sanitize keywords to alphanumeric/hyphen/underscore only** — Prevents special characters from causing unexpected LIKE behavior even though Ent escapes `_` and `%`. Provides defense-in-depth.

4. **Channel send before edit in HandleCallback** — The agent goroutine only needs the boolean approval result. Moving channel send before the Telegram edit API call eliminates blocking on network I/O. The edit is cosmetic (keyboard removal + status text) and can complete asynchronously relative to agent execution. Alternative: goroutine for edit — rejected as channel reorder is simpler and sufficient.

## Risks / Trade-offs

- [Per-keyword OR may return more results than exact phrase match] → Acceptable; results are already ranked by relevance score and limited by `maxPerLayer`
- [Keyword sanitization may remove meaningful characters in non-Latin scripts] → Current user base is primarily Korean/English; regex `[^a-zA-Z0-9\-_]` covers these adequately. Korean characters pass through the stop-word filter and punctuation trim but would be stripped by sanitize — however, the retriever already operates on extracted English keywords for SQLite LIKE matching, and Korean content is handled by the embedding/RAG path
- [Agent may produce response before approval edit completes] → This is the intended behavior; the edit is cosmetic and the previous ordering was causing worse UX
