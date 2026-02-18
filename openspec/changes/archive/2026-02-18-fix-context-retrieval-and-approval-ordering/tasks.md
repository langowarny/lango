## 1. Keyword Extraction Fix

- [x] 1.1 Add constants `_maxSearchKeywords` (5) and `_maxKeywordLength` (50) to `internal/knowledge/retriever.go`
- [x] 1.2 Add `sanitizeKeyword()` function that strips non-alphanumeric characters (except hyphens/underscores) and truncates to max length
- [x] 1.3 Update `extractKeywords()` to apply sanitization and enforce keyword count limit

## 2. Store Per-Keyword OR Predicates

- [x] 2.1 Add `knowledgeKeywordPredicates()` helper to `internal/knowledge/store.go` that splits query into per-keyword `ContentContains`/`KeyContains` OR predicates
- [x] 2.2 Update `SearchKnowledge()` to use `knowledgeKeywordPredicates()`
- [x] 2.3 Add `learningKeywordPredicates()` helper that splits query into per-keyword `ErrorPatternContains`/`TriggerContains` OR predicates
- [x] 2.4 Update `SearchLearnings()` and `SearchLearningEntities()` to use `learningKeywordPredicates()`
- [x] 2.5 Add `externalRefKeywordPredicates()` helper that splits query into per-keyword `NameContains`/`SummaryContains` OR predicates
- [x] 2.6 Update `SearchExternalRefs()` to use `externalRefKeywordPredicates()`

## 3. Approval Message Ordering

- [x] 3.1 In `HandleCallback()` (`internal/channels/telegram/approval.go`), move `pending.ch <- approved` channel send before `editApprovalMessage()` call

## 4. Verification

- [x] 4.1 Run `go build ./...` and verify no compilation errors
- [x] 4.2 Run `go test ./internal/knowledge/...` and verify all tests pass
- [x] 4.3 Run `go test ./internal/channels/telegram/...` and verify all tests pass
- [x] 4.4 Run `go test ./...` and verify full test suite passes
