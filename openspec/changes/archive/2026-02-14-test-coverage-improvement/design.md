# Design: Test Coverage Improvement

## Context
Lango's core business logic (knowledge, learning, skill) and security tools (crypto, secrets) have zero test coverage. Existing tests for session, provider, and app packages cover only minimal scenarios. This creates a reliability gap where regressions can go undetected.

## Goals / Non-Goals

**Goals:**
- Achieve 80%+ coverage for `knowledge`, `learning`, `skill`, `tools/crypto`, `tools/secrets`.
- Achieve 60%+ coverage for `session`.
- Add non-credential-dependent tests for `provider` and `app`.
- Follow existing codebase testing patterns (stdlib `testing`, table-driven, `enttest`).
- Use parallel team agents for efficient implementation.

**Non-Goals:**
- Achieving 100% coverage (not practical for streaming/integration code).
- Adding new testing frameworks (no testify, gomock).
- Changing production code to improve testability.
- Testing external API integrations (provider Generate methods).

## Decisions

### 1. Testing Patterns
Follow the existing codebase conventions confirmed across multiple files:
- **stdlib `testing`** only — no testify or third-party assertion libraries.
- **Table-driven tests**: `tests := []struct{ give, want }` with `t.Run()`.
- **Manual assertions**: `t.Fatal`, `t.Errorf` with descriptive messages.
- **DB tests**: `enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")` for in-memory SQLite.
- **Isolation**: `t.Setenv("HOME", t.TempDir())` for filesystem-dependent tests.

### 2. Test Categorization
Tests are categorized into three types:
- **Pure unit tests** (no DB): `analyzer_test.go`, `builder_test.go`, `extractKeywords` in `retriever_test.go`.
- **DB integration tests** (enttest): `store_test.go`, `retriever_test.go` (DB portion), `engine_test.go`, `registry_test.go`, `secrets_test.go`.
- **Mock-based tests**: `crypto_test.go` uses `mockCryptoProvider` for the `CryptoProvider` interface.

### 3. Parallelized Implementation
Work is divided into 4 independent streams executed by parallel agents:
- **Agent A**: Phase 1 pure unit tests (learning/analyzer, skill/builder)
- **Agent B**: Phase 2 knowledge store + retriever tests
- **Agent C**: Phase 3 learning engine + skill executor/registry
- **Agent D**: Phase 4 security + Phase 5 existing test improvements

### 4. Session Store Test Refactoring
The existing monolithic `TestSQLiteStore` is replaced with 13 individual test functions covering specific scenarios (CRUD, TTL, max history turns, salt/checksum management).

## Risks / Trade-offs

- **enttest in-memory DB**: Slightly different from production SQLite (no file I/O, no encryption). Acceptable for unit-level validation.
- **Mock CryptoProvider**: The mock in `crypto_test.go` doesn't test real encryption. Mitigated by `secrets_test.go` using real `LocalCryptoProvider`.
- **OpenAI ListModels test**: Tests error path only (no real API server). Verifies the error handling code path.
