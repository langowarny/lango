## Context

The project had no `.golangci.yml` configuration, causing golangci-lint v2.4.0 to run with default settings. This included linting ent auto-generated code and flagging standard patterns (defer Close, fmt.Fprint* return values) as errors. 90 issues blocked CI.

## Goals / Non-Goals

**Goals:**
- Zero golangci-lint issues in CI
- Establish `.golangci.yml` v2 config as project standard
- Fix all legitimate code quality issues (unchecked errors, unused code, dead assignments)

**Non-Goals:**
- Changing any runtime behavior or public APIs
- Adding new linters beyond the `standard` default set
- Refactoring code beyond what's needed to fix lint issues

## Decisions

1. **golangci-lint v2 format** — Use `version: "2"` config format with `default: standard` linter set. This matches the CI runner version and provides a good baseline without being overly strict.

2. **`generated: strict` exclusion** — Exclude all files with `// Code generated` headers (ent). This eliminates ~30 false positives from auto-generated code without maintaining manual exclusion lists.

3. **`std-error-handling` preset** — Suppress errcheck for standard patterns (`defer .Close()`, `fmt.Fprint*`, `io.Writer.Write`). These are universally accepted patterns where error handling adds noise without value.

4. **Test file errcheck exclusion** — Disable errcheck in `_test.go` files. Test helpers commonly ignore errors for brevity, and enforcing this in tests adds noise.

5. **`writeJSON` helper in p2p_routes.go** — Rather than adding `_ =` to 12 `json.Encode` calls, extract a helper that properly handles the error. This is the one case where a helper reduces repetition meaningfully.

6. **`_ =` for intentionally ignored errors** — Use explicit `_ =` assignment for errors that are intentionally ignored (rollback in error paths, send-error helpers). This documents intent clearly.

## Risks / Trade-offs

- [Risk] Future ent schema changes might generate code that triggers new lint rules → Mitigation: `generated: strict` handles this automatically via the `// Code generated` header
- [Trade-off] `std-error-handling` preset may suppress some legitimate error checks → Acceptable: the suppressed patterns (defer Close, fmt.Fprint) have extremely low error probability in practice
