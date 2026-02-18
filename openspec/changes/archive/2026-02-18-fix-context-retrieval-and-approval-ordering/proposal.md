## Why

Cron job execution fails with "LIKE or GLOB pattern too complex" SQLite error when `extractKeywords()` generates excessively long search queries. Additionally, tool approval "Approved" messages appear after the agent's final response due to synchronous Telegram API calls blocking the agent pipeline.

## What Changes

- Limit `extractKeywords()` to max 5 keywords, each max 50 characters, with non-alphanumeric character sanitization
- Change knowledge store search methods (`SearchKnowledge`, `SearchLearnings`, `SearchExternalRefs`, `SearchLearningEntities`) from single long LIKE pattern to per-keyword OR predicates
- Reorder `HandleCallback()` in Telegram approval to send channel result before editing the approval message

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `context-retriever`: Keyword extraction now limits count and sanitizes characters to prevent SQLite LIKE pattern complexity overflow
- `knowledge-store`: Search queries split into per-keyword OR predicates instead of single concatenated LIKE pattern
- `channel-approval`: Approval callback unblocks agent before Telegram API edit call to fix message ordering

## Impact

- `internal/knowledge/retriever.go` — `extractKeywords()` with limits and `sanitizeKeyword()` helper
- `internal/knowledge/store.go` — Per-keyword OR predicates in `SearchKnowledge`, `SearchLearnings`, `SearchLearningEntities`, `SearchExternalRefs`
- `internal/channels/telegram/approval.go` — `HandleCallback()` channel send before edit
- No API changes, no breaking changes, no new dependencies
