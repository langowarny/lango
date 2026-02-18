## Context

The current Makefile has 11 targets covering basic build/test/lint/clean/docker workflows. The project now requires Docker Compose orchestration across 3 profiles, CGO-dependent builds, code quality gates, and coverage reporting. The Makefile is the single entry point for all developer workflows.

## Goals / Non-Goals

**Goals:**
- Explicit `CGO_ENABLED=1` on all build/test commands for sqlite3/sqlite-vec safety
- Docker Compose targets matching all 3 profiles (default, browser, browser-sidecar)
- Local CI pipeline (`make ci`) that mirrors future CI/CD gates
- Self-documenting help via `## target: description` comment convention
- Coverage report generation to `.coverage/` directory

**Non-Goals:**
- File watching / hot reload (requires external tool dependency)
- Release automation / goreleaser (no CI/CD pipeline yet)
- Proto/swagger generation (no such files exist)
- Security scanning / vuln checking (better suited for CI pipeline)

## Decisions

1. **Comment-based help generation** — `grep -E '^## '` parses `## target: description` comments and pipes through `column`. Simple, zero-dependency, self-documenting. Alternative: hardcoded echo statements (current approach) — doesn't scale and drifts from reality.

2. **Auto-install golangci-lint in `lint` target** — Checks `which golangci-lint` and installs via `go install` if missing. Reduces onboarding friction. Alternative: fail with install instructions — adds friction for new contributors.

3. **`REGISTRY ?=` as empty default** — `docker-push` validates non-empty `REGISTRY` at runtime. Prevents accidental pushes to wrong registry. Alternative: hardcode a default registry — risky for open-source project.

4. **Sequential `ci` target** — `fmt-check → vet → lint → test` runs in order using Make prerequisite syntax. Fails fast on cheap checks before expensive test suite.

## Risks / Trade-offs

- [Cross-compilation with CGO] → `build-linux` and `build-darwin` require matching C toolchain. This is existing behavior, not introduced by this change.
- [`go install` for golangci-lint] → May install a version incompatible with project. Acceptable for local dev; CI should pin version.
- [Comment-based help] → Relies on `## target:` convention discipline. Low risk — convention is visible and easy to follow.
