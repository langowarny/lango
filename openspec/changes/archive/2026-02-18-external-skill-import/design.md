## Context

Lango's skill system (`internal/skill/`) supports file-based skills with three types: `script`, `template`, and `composite`. Skills are stored as SKILL.md files with YAML frontmatter, managed by a `FileSkillStore`, and converted to agent tools by the `Registry`. External skill repositories (e.g., `github.com/kepano/obsidian-skills`) publish SKILL.md files as agent reference documents — not executable code, but structured knowledge the agent loads on demand.

## Goals / Non-Goals

**Goals:**
- Import SKILL.md files from GitHub repositories and arbitrary URLs
- Support bulk import (all skills from a repo) and single skill import
- Add `instruction` skill type for non-executable reference documents
- Register instruction skills as tools so the agent autonomously decides when to load them
- Track import origin via a `Source` field

**Non-Goals:**
- Version management or auto-update of imported skills
- Authentication for private GitHub repos (public only)
- Execution of instruction skills (they are reference-only)
- System prompt injection of instruction skills (always tool-based)

## Decisions

### 1. Instruction skills as tools, not system prompt
**Decision**: Instruction skills are registered as tools with a content-returning handler, never injected into the system prompt.
**Rationale**: System prompt space is limited. Tool registration lets the agent autonomously reason about "when do I need this reference?" based on the tool description. This follows the same pattern as other skill types.
**Alternative**: Inject into system prompt — rejected because it wastes context window for rarely-needed references.

### 2. GitHub Contents API for repo discovery
**Decision**: Use GitHub Contents API (`/repos/:owner/:repo/contents/:path`) for directory listing and file fetching.
**Rationale**: No git clone needed, works with public repos without authentication, returns base64-encoded file content directly.
**Alternative**: `git clone` — rejected due to dependency overhead and unnecessary full-repo download.

### 3. Default type changed to "instruction"
**Decision**: When a SKILL.md has no explicit `type` in frontmatter, default to `instruction` instead of `script`.
**Rationale**: External SKILL.md files (e.g., Obsidian skills) don't have a `type` field. They are reference documents by nature. All existing internal skills have explicit types, so this change is backward-compatible.

### 4. Importer as standalone component
**Decision**: `Importer` struct in `internal/skill/importer.go` with injected HTTP client.
**Rationale**: Separates fetch/parse logic from registry/store concerns. HTTP client injection enables testing with `httptest.NewServer`.

### 5. Store() accessor on Registry
**Decision**: Expose `Store()` method on Registry for the import tool handler to save skills directly.
**Rationale**: The import tool in `tools.go` needs store access. Exposing the store interface (not internals) maintains encapsulation while enabling the import flow.

## Risks / Trade-offs

- **[GitHub API rate limits]** → Public API allows 60 requests/hour without auth. Bulk import of large repos could hit this. Mitigation: Import is user-initiated, not automatic.
- **[Malicious SKILL.md content]** → Instruction skills contain markdown text only, not executable code. The content is returned to the agent as context, not executed. Risk is minimal.
- **[Default type change]** → Changing default from `script` to `instruction` could affect SKILL.md files without explicit type that were previously parsed as script. Mitigation: All existing internal skills have explicit types.
- **[No duplicate detection across sources]** → Two repos could define skills with the same name. Mitigation: `ImportFromRepo` checks if skill already exists and skips duplicates.
