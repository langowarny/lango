## 1. System Prompts & Infrastructure

- [x] 1.1 Add Skill Tool section to `prompts/TOOL_USAGE.md` with import_skill usage instructions
- [x] 1.2 Add Skills tool category to `prompts/AGENTS.md`
- [x] 1.3 Add git and curl packages to Dockerfile runtime image

## 2. Type & Interface Extensions

- [x] 2.1 Add `AllowedTools []string` field to `SkillEntry` in `internal/skill/types.go`
- [x] 2.2 Add `allowed-tools` field to frontmatter struct in `internal/skill/parser.go`
- [x] 2.3 Parse `allowed-tools` in `ParseSkillMD` using `strings.Fields()`
- [x] 2.4 Serialize `AllowedTools` in `RenderSkillMD` using `strings.Join()`
- [x] 2.5 Add `SaveResource` method to `SkillStore` interface in `internal/skill/store.go`
- [x] 2.6 Implement `SaveResource` in `FileSkillStore` in `internal/skill/file_store.go`

## 3. Importer Enhancement

- [x] 3.1 Add `hasGit()` function using `exec.LookPath("git")`
- [x] 3.2 Add `cloneRepo()` function for shallow git clone to temp directory
- [x] 3.3 Implement `importViaGit()` — clone, scan dirs, parse SKILL.md, save + copy resources
- [x] 3.4 Implement `copyResourceDirs()` — copy scripts/references/assets from source to store
- [x] 3.5 Refactor `ImportFromRepo()` to prefer git clone with HTTP fallback
- [x] 3.6 Implement `importViaHTTP()` with resource fetching support
- [x] 3.7 Add `fetchAndSaveResources()` for HTTP-based resource directory fetching
- [x] 3.8 Add `fetchGitHubFileContent()` for single file fetch via GitHub Contents API
- [x] 3.9 Add `ImportSingleWithResources()` with git/HTTP modes

## 4. Exec Guards & Handler

- [x] 4.1 Add git clone skill redirect guard to `blockLangoExec` in `internal/app/tools.go`
- [x] 4.2 Add curl/wget skill redirect guard to `blockLangoExec`
- [x] 4.3 Update `import_skill` handler to use `ImportSingleWithResources` for single skill imports

## 5. Tests

- [x] 5.1 Add `TestBlockLangoExec_SkillGuards` table-driven test in `internal/app/tools_test.go`
- [x] 5.2 Add `TestParseSkillMD_AllowedTools` in `internal/skill/parser_test.go`
- [x] 5.3 Add `TestRenderSkillMD_AllowedTools_Roundtrip` in `internal/skill/parser_test.go`
- [x] 5.4 Add `TestParseSkillMD_NoAllowedTools` in `internal/skill/parser_test.go`
- [x] 5.5 Add `TestFileSkillStore_SaveResource` in `internal/skill/file_store_test.go`
- [x] 5.6 Add `TestFileSkillStore_SaveResource_NestedDir` in `internal/skill/file_store_test.go`
- [x] 5.7 Add `TestHasGit` in `internal/skill/importer_test.go`
- [x] 5.8 Add `TestCopyResourceDirs` in `internal/skill/importer_test.go`
- [x] 5.9 Add `TestImportViaGit_LocalCloneSimulation` in `internal/skill/importer_test.go`

## 6. Verification

- [x] 6.1 Run `go build ./...` — verify build success
- [x] 6.2 Run `go test ./internal/skill/... ./internal/app/...` — verify all tests pass
