## Context

The `internal/prompt/defaults.go` file contains four Go `const` strings that serve as default system prompts. These are sparse (1-4 lines each) and don't reflect Lango's actual capabilities. Editing multi-paragraph markdown inside Go string literals is awkward and error-prone. Go 1.16+ provides `embed.FS` for embedding static files at build time.

## Goals / Non-Goals

**Goals:**
- Replace hardcoded prompt constants with `.md` files embedded via `go:embed`
- Write production-quality prompts that accurately describe Lango's 5 tools, security model, knowledge system, and multi-channel support
- Maintain backward compatibility — `DefaultBuilder()` API unchanged
- Ensure zero-config operation — prompts embedded in binary, no external files needed

**Non-Goals:**
- Runtime prompt editing or hot-reload (out of scope)
- Prompt versioning or A/B testing infrastructure
- Changing the `Builder`, `Section`, or `LoadFromDir` APIs
- Internationalization of prompts

## Decisions

**D1: Separate `prompts/` package at project root**
- Rationale: `go:embed` requires the embedded files to be in the same package directory or subdirectories. A dedicated package keeps prompt files at the top level for easy discovery and editing.
- Alternative: Embed in `internal/prompt/` — rejected because it mixes infrastructure code with content files and makes prompts harder to find.

**D2: Fallback strings in `defaults.go`**
- Rationale: `embed.FS.ReadFile()` should never fail for correctly built binaries, but a minimal fallback ensures the system degrades gracefully if something unexpected happens.
- Alternative: `log.Fatal` on embed failure — rejected because a degraded prompt is better than a crash.

**D3: Copy prompts to Docker image at `/usr/share/lango/prompts/`**
- Rationale: Prompts are already embedded in the binary, but having them on the filesystem lets operators inspect and reference them. FHS convention places read-only architecture-independent data under `/usr/share/`.
- Alternative: Don't copy — rejected because operators lose visibility into active prompt content.

## Risks / Trade-offs

- **Binary size increase** (~5-10KB for 4 markdown files) → Negligible impact.
- **Prompt content drift** — prompts may become outdated as features evolve → Mitigated by placing prompts in a visible, top-level directory that's easy to review and update.
- **Circular dependency risk** — `internal/prompt` imports `prompts/` package → No risk since `prompts/` has zero imports (only `embed` stdlib).
