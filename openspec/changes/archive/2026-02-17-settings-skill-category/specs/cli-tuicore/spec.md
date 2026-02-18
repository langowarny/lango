## ADDED Requirements

### Requirement: Skill field mappings in UpdateConfigFromForm
The `UpdateConfigFromForm` method SHALL map the following field keys to config paths:
- `skill_enabled` → `config.Skill.Enabled` (boolean)
- `skill_dir` → `config.Skill.SkillsDir` (string)

#### Scenario: Apply skill form values
- **WHEN** a form containing `skill_enabled` and `skill_dir` fields is processed by `UpdateConfigFromForm`
- **THEN** the values SHALL be written to `config.Skill.Enabled` and `config.Skill.SkillsDir` respectively
