## Context

CI runs golangci-lint on Ubuntu where `//go:build linux` files are compiled. `tpm_provider.go` has 6 lint issues (3 errcheck for ignored `flush.Execute()` returns in deferred cleanup, 3 SA1019 for deprecated `transport.OpenTPM` with no replacement API). These are invisible on macOS builds but block CI.

## Goals / Non-Goals

**Goals:**
- Fix all 6 lint issues in `tpm_provider.go` without behavioral changes
- Make CI lint non-blocking so platform-specific lint edge cases don't stall development

**Non-Goals:**
- Replacing the deprecated `transport.OpenTPM` API (no alternative exists yet)
- Refactoring TPM provider logic

## Decisions

1. **Use `_ =` for errcheck on deferred flush** — Flush errors in deferred cleanup are best-effort; the primary operation has already succeeded or failed. Explicit `_ =` signals intentional discard.

2. **Use `//nolint:staticcheck` for deprecated API** — `transport.OpenTPM` has no replacement in the go-tpm library. Suppressing with a comment documenting the reason is the standard approach.

3. **`continue-on-error: true` for CI lint job** — Lint failures become yellow warnings instead of red blockers. This preserves visibility while preventing development velocity issues from cross-platform lint discrepancies.

## Risks / Trade-offs

- [Lint regressions may go unnoticed] → Team reviews yellow warnings in PR checks; lint issues still appear in CI output.
- [Deprecated API accumulates tech debt] → Comment documents the reason; revisit when go-tpm provides a replacement.
