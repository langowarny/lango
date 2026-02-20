## Context

PII Redaction Enhancement code implementation is complete (13 builtin patterns, PIIDetector interface, Presidio integration, Settings TUI, Doctor checks). The documentation layer (README, prompts, example config) has not been updated to reflect these capabilities. This is a docs-only change — no code modifications.

## Goals / Non-Goals

**Goals:**
- Update README.md configuration reference table with 6 new PII/Presidio fields
- Expand README.md AI Privacy Interceptor section with detailed pattern coverage
- Update SAFETY.md prompt so the agent knows its PII protection scope
- Update config.json example with new interceptor fields for Docker headless users

**Non-Goals:**
- No code changes to any Go source files
- No new features or behavioral changes
- No changes to Settings TUI, Doctor, or Onboard (already updated)

## Decisions

1. **README config table placement**: New rows added immediately after existing `piiRegexPatterns` row to maintain logical grouping. Presidio fields use dot-notation (`presidio.enabled`, `presidio.url`, etc.) consistent with existing nested config patterns.

2. **AI Privacy Interceptor section scope**: Expanded with pattern categories (Contact, Identity, Financial, Network) rather than listing all 13 individually. This provides useful detail without overwhelming the reader.

3. **SAFETY.md wording**: Updated to enumerate specific PII categories (national IDs, financial account numbers) and mention 13 builtin patterns and Presidio. This ensures the agent can accurately inform users about protection coverage.

4. **config.json structure**: Added `piiDisabledPatterns` as empty array, `piiCustomPatterns` as empty object, and `presidio` as nested object with defaults. Matches the Go struct defaults in `internal/config/types.go`.

## Risks / Trade-offs

- [Prompt drift] SAFETY.md hardcodes "13 builtin patterns" — if patterns are added/removed later, this number becomes stale → Mitigation: pattern count is stable post-release; future pattern additions should update this prompt
- [Config example drift] config.json example may diverge from actual defaults → Mitigation: `go test ./internal/config/...` validates config loading
