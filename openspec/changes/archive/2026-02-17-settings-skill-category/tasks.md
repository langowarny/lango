## 1. Settings Menu & Form

- [x] 1.1 Add "Skill" category to menu.go (ID: "skill", Title: "Skill", Desc: "File-based skill system")
- [x] 1.2 Update Knowledge category description from "Learning, Skills, Context limits" to "Learning, Context limits"
- [x] 1.3 Add NewSkillForm() to forms_impl.go with Enabled (bool) and Skills Directory (text) fields
- [x] 1.4 Add "skill" case to handleMenuSelection in editor.go

## 2. Config State Mapping

- [x] 2.1 Add skill_enabled and skill_dir field mappings to UpdateConfigFromForm in state_update.go

## 3. Command Description

- [x] 3.1 Update lango settings command Long description to include Skill entry in settings.go

## 4. README Updates

- [x] 4.1 Update README for skill-file-migration: file-based skill description, skill config keys, architecture tree
- [x] 4.2 Update README for onboard-settings-split: lango settings command, 5-step wizard, architecture tree

## 5. Verification

- [x] 5.1 Run go build ./... to verify compilation
- [x] 5.2 Run go test ./internal/cli/settings/... ./internal/cli/tuicore/... to verify tests pass
