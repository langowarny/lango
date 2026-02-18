## MODIFIED Requirements

### Requirement: SkillEntry domain type
`SkillEntry` SHALL include all fields necessary for skill persistence and metadata, including `AllowedTools []string` for pre-approved tool lists.

#### Scenario: SkillEntry with AllowedTools
- **WHEN** a `SkillEntry` is created with `AllowedTools: ["exec", "fs_read"]`
- **THEN** the field SHALL be persisted through Save and restored through Get

### Requirement: SkillStore persistence interface
`SkillStore` SHALL provide a `SaveResource(ctx, skillName, relPath, data)` method for persisting resource files alongside skill SKILL.md files.

#### Scenario: SaveResource writes file to correct path
- **WHEN** `SaveResource` is called with skillName="my-skill" and relPath="scripts/run.sh"
- **THEN** the file SHALL be written to `<store-dir>/my-skill/scripts/run.sh`
