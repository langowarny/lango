# CLI TUI Core Spec

## Goal
The `tuicore` package provides shared TUI form components used by both the onboard wizard and the settings editor.

## Requirements

### Field Types
The package SHALL define the following input types:
- `InputText` — Free-text input
- `InputInt` — Integer input
- `InputBool` — Boolean toggle (spacebar)
- `InputSelect` — Cycle through options (left/right arrows)
- `InputPassword` — Masked text input

### Field Struct
Each field SHALL have:
- Key, Label, Type, Value, Placeholder, Options, Checked, Width, Validate
- Exported `TextInput` field (bubbletea textinput.Model) for cross-package access

### FormModel
The form model SHALL:
- Manage a list of fields with cursor navigation (up/down, tab/shift-tab)
- Support text input, boolean toggle, and select cycling
- Render with title, field labels, and help footer
- Call OnCancel on Esc

### ConfigState
The config state SHALL:
- Hold current `*config.Config` and dirty field tracking
- Provide `UpdateConfigFromForm`, `UpdateProviderFromForm`, `UpdateAuthProviderFromForm` methods
- Map all field keys to their corresponding config paths

### Skill field mappings in UpdateConfigFromForm
The `UpdateConfigFromForm` method SHALL map the following field keys to config paths:
- `skill_enabled` → `config.Skill.Enabled` (boolean)
- `skill_dir` → `config.Skill.SkillsDir` (string)

#### Scenario: Apply skill form values
- **WHEN** a form containing `skill_enabled` and `skill_dir` fields is processed by `UpdateConfigFromForm`
- **THEN** the values SHALL be written to `config.Skill.Enabled` and `config.Skill.SkillsDir` respectively

### Cron field mappings in UpdateConfigFromForm
The `UpdateConfigFromForm` method SHALL map the following field keys to config paths:
- `cron_enabled` → `config.Cron.Enabled` (boolean)
- `cron_timezone` → `config.Cron.Timezone` (string)
- `cron_max_jobs` → `config.Cron.MaxConcurrentJobs` (integer)
- `cron_session_mode` → `config.Cron.DefaultSessionMode` (string)
- `cron_history_retention` → `config.Cron.HistoryRetention` (string)

#### Scenario: Apply cron form values
- **WHEN** a form containing cron fields is processed by `UpdateConfigFromForm`
- **THEN** the values SHALL be written to the corresponding `config.Cron` fields

### Background field mappings in UpdateConfigFromForm
The `UpdateConfigFromForm` method SHALL map the following field keys to config paths:
- `bg_enabled` → `config.Background.Enabled` (boolean)
- `bg_yield_ms` → `config.Background.YieldMs` (integer)
- `bg_max_tasks` → `config.Background.MaxConcurrentTasks` (integer)

#### Scenario: Apply background form values
- **WHEN** a form containing background fields is processed by `UpdateConfigFromForm`
- **THEN** the values SHALL be written to the corresponding `config.Background` fields

### Workflow field mappings in UpdateConfigFromForm
The `UpdateConfigFromForm` method SHALL map the following field keys to config paths:
- `wf_enabled` → `config.Workflow.Enabled` (boolean)
- `wf_max_steps` → `config.Workflow.MaxConcurrentSteps` (integer)
- `wf_timeout` → `config.Workflow.DefaultTimeout` (duration parsed from string)
- `wf_state_dir` → `config.Workflow.StateDir` (string)

#### Scenario: Apply workflow form values
- **WHEN** a form containing workflow fields is processed by `UpdateConfigFromForm`
- **THEN** the values SHALL be written to the corresponding `config.Workflow` fields
