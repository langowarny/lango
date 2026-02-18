## Why

When the agent tries to import external skills, it attempts `exec` with `git clone` or `curl` directly — leading to incorrect storage paths and failures when tools are unavailable. The `import_skill` tool needs to handle git clone internally (with HTTP API fallback) and support resource directories (scripts, references, assets) alongside SKILL.md files.

## What Changes

- `import_skill` tool now prefers `git clone` (shallow, depth=1) when git is available, falls back to GitHub HTTP API when it is not.
- Resource directories (`scripts/`, `references/`, `assets/`) are automatically imported alongside SKILL.md files.
- `blockLangoExec` redirects skill-related `git clone`/`curl`/`wget` commands to the `import_skill` tool with guidance messages.
- `SkillEntry` gains an `AllowedTools` field, parsed from the `allowed-tools` YAML frontmatter in SKILL.md.
- `SkillStore` interface gains a `SaveResource` method for persisting resource files.
- System prompts (`TOOL_USAGE.md`, `AGENTS.md`) now include Skill Tool usage instructions.
- Dockerfile runtime image includes `git` and `curl` packages.

## Capabilities

### New Capabilities
- `skill-resource-dirs`: Support for importing and persisting skill resource directories (scripts, references, assets) alongside SKILL.md files.
- `skill-allowed-tools`: Parse and serialize `allowed-tools` frontmatter field in SKILL.md for pre-approved tool lists.

### Modified Capabilities
- `skill-import`: Enhanced with git clone priority + HTTP fallback, resource directory support, and exec guard redirects.
- `skill-system`: Added `AllowedTools` field to `SkillEntry` and `SaveResource` to `SkillStore` interface.

## Impact

- **Code**: `internal/skill/` (importer.go, types.go, parser.go, store.go, file_store.go), `internal/app/tools.go`
- **Prompts**: `prompts/TOOL_USAGE.md`, `prompts/AGENTS.md`
- **Infrastructure**: `Dockerfile` (added git, curl to runtime image)
- **Tests**: New tests in `internal/app/tools_test.go`, updated tests in `internal/skill/parser_test.go`, `internal/skill/file_store_test.go`, `internal/skill/importer_test.go`
