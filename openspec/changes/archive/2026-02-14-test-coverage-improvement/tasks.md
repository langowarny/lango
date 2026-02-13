## 1. Phase 1: Pure Unit Tests

- [x] 1.1 Create `internal/learning/analyzer_test.go` with TestExtractErrorPattern, TestCategorizeError, TestIsDeadlineExceeded, TestSummarizeParams.
- [x] 1.2 Create `internal/skill/builder_test.go` with TestBuildCompositeSkill, TestBuildScriptSkill, TestBuildTemplateSkill.

## 2. Phase 2: Knowledge Store + Retriever Tests

- [x] 2.1 Create `internal/knowledge/store_test.go` with 18 test functions covering CRUD, search, rate limiting for all entity types.
- [x] 2.2 Create `internal/knowledge/retriever_test.go` with TestExtractKeywords (pure unit) + TestContextRetriever_Retrieve, TestAssemblePrompt (DB integration).

## 3. Phase 3: Learning Engine + Skill Tests

- [x] 3.1 Create `internal/learning/engine_test.go` with TestEngine_OnToolResult_Success/Error, TestEngine_GetFixForError, TestEngine_RecordUserCorrection.
- [x] 3.2 Create `internal/skill/executor_test.go` with TestValidateScript, TestExecute_Composite/Template/Script/UnknownType.
- [x] 3.3 Create `internal/skill/registry_test.go` with TestRegistry_CreateSkill_Validation, TestRegistry_LoadSkills_AllTools, TestRegistry_ActivateSkill, TestRegistry_GetSkillTool.

## 4. Phase 4: Security Tool Tests

- [x] 4.1 Create `internal/tools/crypto/crypto_test.go` with TestCryptoTool_Hash/Encrypt/Decrypt/Sign/Keys and TestMapToStruct using mock CryptoProvider.
- [x] 4.2 Create `internal/tools/secrets/secrets_test.go` with TestSecretsTool_Store/Get/List/Delete/UpdateExisting using real LocalCryptoProvider.

## 5. Phase 5: Existing Test Improvements

- [x] 5.1 Rewrite `internal/session/store_test.go` from 1 monolithic test to 13 individual test functions (CreateAndGet, Get_NotFound, Update, Delete_Idempotent, AppendMessage_WithToolCalls, MaxHistoryTurns, TTL, Salt/Checksum management).
- [x] 5.2 Add TestAnthropicProvider_ListModels to `internal/provider/anthropic/anthropic_test.go`.
- [x] 5.3 Add TestOpenAIProvider_ListModels to `internal/provider/openai/openai_test.go`.
- [x] 5.4 Add TestNew_NoProviders and TestNew_InvalidProviderType to `internal/app/app_test.go`.

## 6. Verification

- [x] 6.1 Run `go test -v -race -cover` on all affected packages — all tests pass with race detector.
- [x] 6.2 Verify coverage targets: knowledge 85.1%, learning 86.3%, skill 85.4%, crypto 86.2%, secrets 83.3%, session 64.1%.
