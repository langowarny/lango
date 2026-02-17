## Why

The Skill system was migrated to file-based storage with `SkillConfig` (`skill.enabled`, `skill.skillsDir`), but the `lango settings` editor has no UI to configure these fields. Users must manually edit config JSON to manage skill settings, while all other config sections have dedicated TUI forms.

## What Changes

- Add a new "Skill" menu category to the `lango settings` editor with Enabled and Skills Directory fields
- Add `NewSkillForm()` function following the existing form pattern
- Add `skill_enabled` and `skill_dir` field mappings to `UpdateConfigFromForm`
- Update Knowledge menu description from "Learning, Skills, Context limits" to "Learning, Context limits" (skills are now independent)
- Update `lango settings` command Long description to include Skill entry
- Update README to reflect both skill-file-migration and onboard-settings-split changes

## Capabilities

### New Capabilities

### Modified Capabilities
- `cli-settings`: Add Skill category to the settings menu and form routing
- `cli-tuicore`: Add skill field mappings to UpdateConfigFromForm

## Impact

- `internal/cli/settings/menu.go` — new category entry, Knowledge description update
- `internal/cli/settings/editor.go` — new case in handleMenuSelection
- `internal/cli/settings/forms_impl.go` — new NewSkillForm function
- `internal/cli/settings/settings.go` — Long description update
- `internal/cli/tuicore/state_update.go` — new skill field mappings
- `README.md` — updated for skill-file-migration and onboard-settings-split changes
