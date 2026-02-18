## 1. SkillEntry Extension

- [x] 1.1 Add `Source` string field to `SkillEntry` in `internal/skill/types.go`
- [x] 1.2 Update type comment to include `instruction`

## 2. Instruction Skill Type — Parser

- [x] 2.1 Add `Source` field to `frontmatter` struct in `internal/skill/parser.go`
- [x] 2.2 Change default type from `script` to `instruction` in `ParseSkillMD`
- [x] 2.3 Add `instruction` case to `parseBody` (store body as `definition["content"]`)
- [x] 2.4 Add `instruction` case to `RenderSkillMD` (output raw content)
- [x] 2.5 Wire `Source` field in both parse and render paths for roundtrip

## 3. Instruction Skill Type — Executor

- [x] 3.1 Add `instruction` case to `Execute` in `internal/skill/executor.go` (return content directly)

## 4. Instruction Skill Type — Registry

- [x] 4.1 Update `CreateSkill` validation to accept `instruction` type
- [x] 4.2 Allow instruction type with empty Definition
- [x] 4.3 Add instruction-specific branch to `skillToTool` with content-returning handler
- [x] 4.4 Add `Store()` accessor method to Registry

## 5. Builder

- [x] 5.1 Add `BuildInstructionSkill` function in `internal/skill/builder.go`

## 6. Importer

- [x] 6.1 Create `internal/skill/importer.go` with `Importer` struct, `GitHubRef`, and `ImportResult` types
- [x] 6.2 Implement `ParseGitHubURL` and `IsGitHubURL` functions
- [x] 6.3 Implement `DiscoverSkills` (GitHub Contents API directory listing)
- [x] 6.4 Implement `FetchSkillMD` (GitHub Contents API file fetch with base64 decode)
- [x] 6.5 Implement `FetchFromURL` (arbitrary URL HTTP GET)
- [x] 6.6 Implement `ImportFromRepo` (discover → fetch → parse → save loop with skip/error tracking)
- [x] 6.7 Implement `ImportSingle` (parse raw content → save with source)

## 7. Agent Tool

- [x] 7.1 Add `import_skill` tool to `buildMetaTools` in `internal/app/tools.go`
- [x] 7.2 Implement handler: GitHub bulk import, GitHub single import, direct URL import
- [x] 7.3 Call `registry.LoadSkills(ctx)` after import to reload tool list
- [x] 7.4 Record audit log entries for import actions

## 8. Config

- [x] 8.1 Add `AllowImport bool` to `SkillConfig` in `internal/config/types.go`
- [x] 8.2 Set default `AllowImport: true` in `internal/config/loader.go`

## 9. Tests

- [x] 9.1 Create `internal/skill/importer_test.go` with ParseGitHubURL, IsGitHubURL, FetchFromURL, ImportFromRepo tests
- [x] 9.2 Extend `internal/skill/parser_test.go` with instruction type and source roundtrip tests
- [x] 9.3 Extend `internal/skill/registry_test.go` with instruction tool conversion, content return, description tests
- [x] 9.4 Update existing registry validation test for new error message
- [x] 9.5 Verify `go build ./...` passes
- [x] 9.6 Verify `go test ./internal/skill/...` passes
- [x] 9.7 Verify `go test ./...` passes
