## MODIFIED Requirements

### Requirement: Embedded Default Skills
The system SHALL embed default skill files via `//go:embed **/SKILL.md`. When no real skill SKILL.md files are present, a `.placeholder/SKILL.md` file SHALL exist to satisfy the embed glob pattern. The placeholder SHALL NOT contain valid YAML frontmatter and SHALL NOT be deployed as a usable skill. `EnsureDefaults()` SHALL skip any embedded path whose directory name starts with `.` (hidden directories).

#### Scenario: Build with no real default skills
- **WHEN** `go build` is run with only `.placeholder/SKILL.md` in the skills directory
- **THEN** the build SHALL succeed without errors

#### Scenario: Placeholder not deployed as skill
- **WHEN** `EnsureDefaults()` iterates over the embedded filesystem
- **THEN** entries whose directory name starts with `.` SHALL be skipped entirely
- **AND** no files from `.placeholder/` SHALL be written to the user's skills directory

#### Scenario: Future skill addition
- **WHEN** a new `skills/<name>/SKILL.md` file with valid frontmatter is added
- **THEN** it SHALL be automatically included in the embedded filesystem and deployed via `EnsureDefaults()`

#### Scenario: Existing skills preserved
- **WHEN** `EnsureDefaults()` is called and a skill directory already exists
- **THEN** that skill SHALL NOT be overwritten

### Requirement: File-Based Skill Storage
The system SHALL store skills as `<dir>/<name>/SKILL.md` files with YAML frontmatter containing name, description, type, status, and optional parameters. `ListActive()` SHALL skip hidden directories (names starting with `.`) when scanning.

#### Scenario: Save a new skill
- **WHEN** a skill entry is saved via `FileSkillStore.Save()`
- **THEN** the system SHALL create `<skillsDir>/<name>/SKILL.md` with YAML frontmatter and markdown body

#### Scenario: Load active skills
- **WHEN** `FileSkillStore.ListActive()` is called
- **THEN** all skills with `status: active` in their frontmatter SHALL be returned
- **AND** directories whose name starts with `.` SHALL be skipped without logging a warning

#### Scenario: Hidden directory ignored
- **WHEN** `FileSkillStore.ListActive()` encounters a directory starting with `.`
- **THEN** it SHALL skip the directory silently without attempting to parse its contents

#### Scenario: Delete a skill
- **WHEN** `FileSkillStore.Delete()` is called with a skill name
- **THEN** the entire `<skillsDir>/<name>/` directory SHALL be removed

#### Scenario: SaveResource writes file to correct path
- **WHEN** `SaveResource` is called with skillName="my-skill" and relPath="scripts/run.sh"
- **THEN** the file SHALL be written to `<store-dir>/my-skill/scripts/run.sh`
