## 1. Onboard Forms — Field Additions

- [x] 1.1 Add `browser_enabled` (InputBool) field before `browser_headless` in `NewToolsForm`
- [x] 1.2 Add `browser_session_timeout` (InputText, duration) field after `browser_headless` in `NewToolsForm`
- [x] 1.3 Add `max_history_turns` (InputInt) field after `ttl` in `NewSecurityForm`
- [x] 1.4 Add `system_prompt_path` (InputText), `fallback_provider` (InputSelect), `fallback_model` (InputText) to `NewAgentForm`

## 2. Onboard Forms — Knowledge Form

- [x] 2.1 Create `NewKnowledgeForm` function with 6 fields: enabled, max_learnings, max_knowledge, max_context, auto_approve, max_skills_day
- [x] 2.2 Add `{"knowledge", "🧠 Knowledge", "Learning, Skills, Context limits"}` category to `NewMenuModel` in `menu.go`
- [x] 2.3 Add `case "knowledge":` routing to `handleMenuSelection` in `wizard.go`

## 3. Config State Mapping

- [x] 3.1 Add `browser_enabled`, `browser_session_timeout` cases to `UpdateConfigFromForm` in `state_update.go`
- [x] 3.2 Add `max_history_turns` case to `UpdateConfigFromForm`
- [x] 3.3 Add `system_prompt_path`, `fallback_provider`, `fallback_model` cases to `UpdateConfigFromForm`
- [x] 3.4 Add all 6 Knowledge field cases (`knowledge_enabled`, `knowledge_max_learnings`, `knowledge_max_knowledge`, `knowledge_max_context`, `knowledge_auto_approve`, `knowledge_max_skills_day`) to `UpdateConfigFromForm`

## 4. Bug Fixes and Description Updates

- [x] 4.1 Remove duplicate `w.activeForm.Focus = true` in `wizard.go` security case
- [x] 4.2 Update `onboard.go` Long description to list all 7 configurable sections
- [x] 4.3 Replace if/else-if chain with switch in `SecurityCheck.Run` and add `enclave` case in `security.go`
- [x] 4.4 Update `doctor.go` Long description to list all 7 checks

## 5. Tests

- [x] 5.1 Write `TestNewAgentForm_AllFields` — verify 7 fields exist with correct values
- [x] 5.2 Write `TestNewToolsForm_AllFields` — verify 6 fields exist with correct defaults
- [x] 5.3 Write `TestNewSecurityForm_AllFields` — verify 10 fields exist
- [x] 5.4 Write `TestNewKnowledgeForm_AllFields` — verify 6 fields exist with correct defaults
- [x] 5.5 Write `TestUpdateConfigFromForm_*` — verify Agent, Browser, MaxHistoryTurns, Knowledge config mapping
- [x] 5.6 Write `TestNewMenuModel_HasKnowledgeCategory` — verify menu includes knowledge
- [x] 5.7 Write `TestSecurityCheck_Run_EnclaveProvider` — enclave does not return Fail
- [x] 5.8 Write `TestSecurityCheck_Run_UnknownProvider` — unknown provider returns Fail
