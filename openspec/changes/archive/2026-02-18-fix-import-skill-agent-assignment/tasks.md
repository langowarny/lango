## 1. Prefix Routing Fix

- [x] 1.1 Add `"import_skill"` to librarian's `Prefixes` slice in `agentSpecs` (`internal/orchestration/tools.go`)
- [x] 1.2 Add `"import_skill"` entry to `capabilityMap` (`internal/orchestration/tools.go`)

## 2. Verification

- [x] 2.1 Verify `go build ./...` passes
- [x] 2.2 Verify `go test ./internal/orchestration/...` passes
- [x] 2.3 Verify `import_skill` is partitioned to librarian (not unmatched) in existing tests
