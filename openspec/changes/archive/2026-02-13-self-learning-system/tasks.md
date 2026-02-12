## 1. Ent Schemas & Database

- [x] 1.1 Create `Knowledge` Ent schema with fields: key (unique), category (enum), content, tags (JSON), source, relevance_score, use_count, timestamps
- [x] 1.2 Create `Learning` Ent schema with fields: trigger, error_pattern, diagnosis, fix, category (enum), tags (JSON), confidence, occurrence_count, success_count, timestamps
- [x] 1.3 Create `Skill` Ent schema with fields: name (unique), description, skill_type (enum), definition (JSON), parameters (JSON), status (enum), created_by, requires_approval, use/success counts, timestamps
- [x] 1.4 Create `AuditLog` Ent schema with fields: session_key, action (enum), actor, target, details (JSON), created_at
- [x] 1.5 Create `ExternalRef` Ent schema with fields: name (unique), ref_type (enum), location, summary, metadata (JSON), timestamps
- [x] 1.6 Run `go generate ./internal/ent` and verify auto-migration

## 2. Knowledge Store (`internal/knowledge`)

- [x] 2.1 Define domain types in `types.go`: KnowledgeEntry, LearningEntry, SkillEntry, AuditEntry, ExternalRefEntry, ContextLayer, ContextItem, RetrievalRequest/Result
- [x] 2.2 Implement `Store` in `store.go` with CRUD for Knowledge (save/get/search/delete/increment)
- [x] 2.3 Implement Learning CRUD (save/search/searchEntities/boostConfidence)
- [x] 2.4 Implement Skill persistence (save/get/listActive/activate/incrementUsage)
- [x] 2.5 Implement AuditLog and ExternalRef persistence
- [x] 2.6 Implement atomic per-session rate limiting with `reserveKnowledgeSlot`/`reserveLearningSlot`

## 3. Context Retriever (`internal/knowledge`)

- [x] 3.1 Implement `ContextRetriever` in `retriever.go` with `Retrieve()` method searching 4 layers
- [x] 3.2 Implement `AssemblePrompt()` to build augmented system prompt with markdown sections
- [x] 3.3 Implement `extractKeywords()` with stop-word filtering and 2-char minimum length
- [x] 3.4 Implement per-layer retrieval: `retrieveKnowledge`, `retrieveSkills`, `retrieveExternalRefs`, `retrieveLearnings`

## 4. Learning Engine (`internal/learning`)

- [x] 4.1 Implement `Engine` in `engine.go` with `OnToolResult()` observer method
- [x] 4.2 Implement `handleError()` for error pattern extraction, categorization, and learning creation
- [x] 4.3 Implement `handleSuccess()` for confidence boosting of related learnings
- [x] 4.4 Implement `analyzer.go` with `extractErrorPattern()`, `categorizeError()`, and `summarizeParams()`
- [x] 4.5 Fix `isDeadlineExceeded()` to use `errors.Is()` for wrapped error detection

## 5. Skill System (`internal/skill`)

- [x] 5.1 Implement `Registry` in `registry.go` with CreateSkill, ActivateSkill, GetSkill, LoadSkills
- [x] 5.2 Implement skill type validation (composite: steps array, script: script string, template: template string)
- [x] 5.3 Implement `Executor` in `executor.go` with executeComposite, executeScript, executeTemplate
- [x] 5.4 Implement dangerous pattern validation in `ValidateScript()`
- [x] 5.5 Implement `Builder` in `builder.go` with BuildFromSteps and BuildScript
- [x] 5.6 Fix `NewExecutor()` to properly handle `os.UserHomeDir()` and `os.MkdirAll()` errors

## 6. Meta-Tools & Wiring (`internal/app`)

- [x] 6.1 Implement `buildMetaTools()` in `tools.go` with 6 meta-tools: save_knowledge, search_knowledge, save_learning, search_learnings, create_skill, list_skills
- [x] 6.2 Implement `wrapWithLearning()` tool handler wrapper
- [x] 6.3 Add `KnowledgeConfig` to config types with enabled, limits, and auto-approve settings
- [x] 6.4 Implement `initKnowledge()` in `wiring.go` to initialize Store, Engine, Registry
- [x] 6.5 Integrate `ContextAwareModelAdapter` in `initAgent()` when knowledge is enabled
- [x] 6.6 Wire meta-tools and learning wrapper into agent tool registration

## 7. Verification

- [x] 7.1 Run `go build ./...` — clean build with no errors
- [x] 7.2 Run `go vet ./...` — no warnings or issues
- [x] 7.3 Code review and fix critical quality issues (TOCTOU race, errors.Is, error handling)
