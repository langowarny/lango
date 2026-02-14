## Why

Recent backend commits (browser tool recover, security refactoring, phantom feature audit, knowledge tests) introduced significant new config fields—Browser `enabled`/`sessionTimeout`, the entire Knowledge system, Session `maxHistoryTurns`, and Agent fallback/system-prompt settings—but the Onboard TUI wizard and Doctor CLI were never updated. Users cannot configure these features through the TUI, and Doctor fails to recognize the `enclave` signer provider.

## What Changes

- **Onboard TUI — Browser fields**: Add `browser_enabled` and `browser_session_timeout` fields to the Tools form
- **Onboard TUI — Knowledge form**: Create a new `NewKnowledgeForm` with 6 fields (enabled, maxLearnings, maxKnowledge, maxContextPerLayer, autoApproveSkills, maxSkillsPerDay) and wire it into menu/wizard
- **Onboard TUI — Session field**: Add `max_history_turns` to the Security form
- **Onboard TUI — Agent fields**: Add `system_prompt_path`, `fallback_provider`, `fallback_model` to the Agent form
- **Onboard TUI — descriptions**: Update `onboard` command Long description to reflect actual menu-based flow
- **Doctor — enclave provider**: SecurityCheck now recognizes `enclave` as a valid signer provider (switch instead of if/else)
- **Doctor — description**: Update `doctor` command Long description to list all 7 checks
- **Wizard — duplicate Focus**: Remove duplicate `w.activeForm.Focus = true` assignment in security case

## Capabilities

### New Capabilities

_(none — all changes modify existing capabilities)_

### Modified Capabilities

- `cli-onboard`: Add Knowledge menu category, Browser/Session/Agent field coverage, updated command description
- `cli-doctor`: Recognize `enclave` signer provider in SecurityCheck, updated command description
- `tool-browser`: TUI now exposes `enabled` and `sessionTimeout` config fields
- `knowledge-store`: TUI now provides a dedicated Knowledge configuration form

## Impact

- **Files modified**: `forms_impl.go`, `state_update.go`, `menu.go`, `wizard.go`, `onboard.go`, `doctor/checks/security.go`, `doctor/doctor.go`
- **New files**: `forms_impl_test.go` (onboard form tests)
- **Dependencies**: None added
- **Breaking changes**: None
