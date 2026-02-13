# Proposal: Test Coverage Improvement (P0 Core + P2 Existing)

## Summary
Improve test coverage across the Lango project from ~25% to ~55% by adding comprehensive tests for untested core packages (knowledge, learning, skill, crypto, secrets) and enhancing existing tests (session, provider, app).

## Problem Statement
The project's test coverage is critically low at approximately 25%. Core business logic packages have 0% coverage:
- `internal/knowledge/` — CRUD store for knowledge, learnings, skills, audit logs, external refs
- `internal/learning/` — Error pattern analysis and learning engine
- `internal/skill/` — Skill building, execution, and registry
- `internal/tools/crypto/` — Cryptographic operations tool
- `internal/tools/secrets/` — Secrets management tool

Existing tests cover only minimal happy paths:
- `internal/session/` — Single monolithic test (~15%)
- `internal/provider/` — Constructor-only tests (~5%)
- `internal/app/` — All tests require credentials (skipped, 0% effective)

## Proposed Solution
Create 9 new test files and modify 4 existing test files across 5 phases:

1. **Phase 1 (Pure Unit)**: `analyzer_test.go`, `builder_test.go` — no DB required
2. **Phase 2 (DB Integration)**: `knowledge/store_test.go`, `knowledge/retriever_test.go` — enttest in-memory SQLite
3. **Phase 3 (Store-dependent)**: `engine_test.go`, `executor_test.go`, `registry_test.go`
4. **Phase 4 (Security)**: `crypto/crypto_test.go`, `secrets/secrets_test.go` — mock + real providers
5. **Phase 5 (Existing)**: Improve `session/store_test.go`, `anthropic_test.go`, `openai_test.go`, `app_test.go`

## Capabilities

### New Capabilities
- `test-coverage`: Comprehensive test suite for core business logic and security tools.

### Modified Capabilities
- `session-store`: Enhanced test coverage with individual test functions.
- `provider-anthropic`, `provider-openai`: ListModels test coverage.
- `application-core`: Non-credential-dependent test cases.

## Impact
- **Codebase**: 13 files touched (9 new, 4 modified), test-only changes
- **Dependencies**: No new external dependencies (uses existing `enttest`, `go.uber.org/zap`)
- **Testing**: Coverage increase from ~25% to ~55% overall
- **Risk**: Zero risk to production — all changes are test files only
