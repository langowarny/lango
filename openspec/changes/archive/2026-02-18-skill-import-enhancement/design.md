## Context

The skill import system currently uses only the GitHub Contents API (HTTP) to fetch SKILL.md files from repositories. This has two limitations: (1) the agent bypasses `import_skill` and directly runs `git clone` via `exec`, storing skills in arbitrary paths instead of `~/.lango/skills/`, and (2) resource files (scripts, references, assets) accompanying skills are not imported. The Dockerfile runtime image also lacks `git` and `curl`, making shell-based fallbacks impossible in Docker environments.

## Goals / Non-Goals

**Goals:**
- `import_skill` tool SHALL prefer `git clone` (shallow, depth=1) when git is available for faster, full-directory imports
- `import_skill` SHALL fall back to GitHub HTTP API when git is unavailable
- Resource directories (`scripts/`, `references/`, `assets/`) SHALL be imported alongside SKILL.md
- `blockLangoExec` SHALL redirect skill-related `git clone`/`curl`/`wget` to `import_skill`
- `SkillEntry` SHALL support `AllowedTools` field parsed from YAML frontmatter
- `SkillStore` SHALL support `SaveResource` for persisting resource files
- Docker runtime image SHALL include `git` and `curl`

**Non-Goals:**
- Private repository authentication (GitHub token) — future enhancement
- Recursive resource directory traversal (only top-level files in resource dirs)
- Skill dependency resolution between skills
- Skill versioning or update detection

## Decisions

### Decision 1: Git clone as primary import method
**Choice**: Use `exec.LookPath("git")` to detect git availability, then `git clone --depth=1 --branch <branch>` to a temp directory.
**Rationale**: Git clone is significantly faster than individual GitHub API calls for repos with many skills, and automatically fetches all resource files. Shallow clone (depth=1) minimizes bandwidth.
**Alternative**: Always use HTTP API — rejected because it requires N+1 API calls (1 directory listing + N file fetches) and GitHub rate limits apply.

### Decision 2: Automatic fallback to HTTP API
**Choice**: If `hasGit()` returns false or `cloneRepo()` fails, transparently fall back to the existing HTTP API approach.
**Rationale**: Ensures import works in environments without git (minimal Docker images, restricted environments).

### Decision 3: Recognized resource directories are fixed
**Choice**: Only `scripts/`, `references/`, and `assets/` are recognized as resource directories.
**Rationale**: Matches the Anthropic skills repository convention. Keeps the implementation simple and predictable.

### Decision 4: Exec guard uses keyword matching
**Choice**: Check `strings.Contains(lower, "skill")` for git clone/curl/wget commands.
**Rationale**: Simple heuristic that catches the common case (`git clone ... skills ...`) without false positives on unrelated repositories.

## Risks / Trade-offs

- [Git clone timeout] → Context timeout propagated to `exec.CommandContext`; overall import timeout applies.
- [Temp directory cleanup] → `defer os.RemoveAll(cloneDir)` ensures cleanup even on error.
- [Git clone failure on network issues] → Automatic fallback to HTTP API.
- [Keyword-based exec guard false negatives] → Agent could use unrelated commands to clone skills; mitigated by system prompt instructions.
