## Why

The default prompts in `internal/prompt/defaults.go` are hardcoded as Go string constants — a few generic lines that don't reflect Lango's actual capabilities (5 tool categories, 6-layer knowledge system, security architecture, multi-channel support). Managing multi-paragraph prompts as Go constants is error-prone and hard to review. Moving them to `.md` files with `go:embed` makes prompts editable, reviewable, and production-quality while remaining embedded in the binary.

## What Changes

- Create `prompts/` package at project root with 4 production-quality `.md` prompt files
- Add `prompts/embed.go` with `go:embed *.md` to expose an `embed.FS`
- Replace `const` strings in `internal/prompt/defaults.go` with `prompts.FS.ReadFile()` calls
- Add fallback strings for resilience if embedded reads fail
- Add `prompts/embed_test.go` to verify all files are correctly embedded
- Update `Dockerfile` to copy prompt files for user reference (already embedded in binary)

## Capabilities

### New Capabilities
- `embedded-prompt-files`: Production-quality default prompts stored as `.md` files and embedded into the binary via `go:embed`, covering agent identity, safety rules, conversation rules, and tool usage guidelines.

### Modified Capabilities
- `structured-prompt-builder`: `DefaultBuilder()` now reads content from embedded filesystem instead of hardcoded constants. Public API unchanged.

## Impact

- **New package**: `prompts/` (4 `.md` files, `embed.go`, `embed_test.go`)
- **Modified**: `internal/prompt/defaults.go` (imports `prompts` package, removes `const` declarations)
- **Modified**: `internal/prompt/defaults_test.go`, `internal/prompt/loader_test.go` (assertions updated for new content)
- **Modified**: `Dockerfile` (adds `COPY` for reference prompt files)
- **No breaking changes**: `DefaultBuilder()` signature and behavior unchanged
