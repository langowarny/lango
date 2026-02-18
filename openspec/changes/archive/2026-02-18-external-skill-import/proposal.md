## Why

Lango's skill system only supports locally created skills (script/template/composite). There is no way to import skills from external sources like GitHub repositories. External "instruction" skills (YAML frontmatter + markdown body) serve as agent reference documents that are loaded on-demand via tool invocation, enabling the agent to autonomously decide when contextual knowledge is needed.

## What Changes

- Add `instruction` skill type: agent reference documents stored as markdown content, registered as tools for on-demand context loading
- Add `Source` field to `SkillEntry` for tracking import origin URLs
- Add `Importer` component: fetches SKILL.md files from GitHub repositories (via Contents API) or arbitrary URLs
- Add `import_skill` agent tool: supports bulk repo import and single skill import
- Add `AllowImport` config flag to `SkillConfig`
- Parser/Renderer support for instruction type and Source field roundtrip
- Executor handles instruction type (returns content directly)
- Registry converts instruction skills to tools with content-returning handlers

## Capabilities

### New Capabilities
- `skill-import`: External skill import from GitHub repositories and URLs, including GitHub Contents API integration, URL fetching, bulk/single import, and SKILL.md parsing

### Modified Capabilities
- `skill-system`: Add `instruction` skill type, `Source` field on SkillEntry, instruction-specific tool conversion in registry, and `AllowImport` config flag

## Impact

- `internal/skill/types.go`: New `Source` field on `SkillEntry`
- `internal/skill/parser.go`: Instruction type parsing/rendering, Source frontmatter roundtrip
- `internal/skill/executor.go`: Instruction type execution case
- `internal/skill/registry.go`: Instruction tool conversion, `Store()` accessor, validation update
- `internal/skill/builder.go`: `BuildInstructionSkill` function
- `internal/skill/importer.go`: New file — GitHub API client, URL fetcher, import orchestration
- `internal/app/tools.go`: `import_skill` agent tool
- `internal/config/types.go`: `AllowImport` flag on `SkillConfig`
- `internal/config/loader.go`: Default `AllowImport: true`
