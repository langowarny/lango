## ADDED Requirements

### Requirement: Resource directory import via git clone
The system SHALL copy files from recognized resource directories (`scripts/`, `references/`, `assets/`) when importing skills via git clone.

#### Scenario: Skill with scripts directory
- **WHEN** a skill directory contains a `scripts/` subdirectory with files
- **THEN** all files in `scripts/` SHALL be copied to `~/.lango/skills/<name>/scripts/`

#### Scenario: Skill with multiple resource directories
- **WHEN** a skill directory contains `scripts/`, `references/`, and `assets/` subdirectories
- **THEN** all files from each directory SHALL be copied to the corresponding paths under `~/.lango/skills/<name>/`

#### Scenario: Skill with no resource directories
- **WHEN** a skill directory contains only SKILL.md and no recognized resource directories
- **THEN** no resource files SHALL be copied and no error SHALL occur

### Requirement: Resource directory import via HTTP API
The system SHALL fetch files from recognized resource directories via GitHub Contents API when git is unavailable.

#### Scenario: HTTP fallback fetches resources
- **WHEN** git is not available and a skill has a `scripts/` directory in the GitHub repo
- **THEN** the system SHALL use the GitHub Contents API to list and download each file in the directory
- **AND** save them to `~/.lango/skills/<name>/scripts/`

#### Scenario: HTTP resource directory not found
- **WHEN** the GitHub Contents API returns 404 for a resource directory
- **THEN** the system SHALL skip that directory without error

### Requirement: SaveResource persists resource files
The `SkillStore` interface SHALL provide a `SaveResource(ctx, skillName, relPath, data)` method that writes resource files under the skill's directory.

#### Scenario: Save a resource file
- **WHEN** `SaveResource` is called with skillName="my-skill", relPath="scripts/setup.sh", and data bytes
- **THEN** the file SHALL be written to `<skills-dir>/my-skill/scripts/setup.sh`

#### Scenario: Save resource creates parent directories
- **WHEN** the parent directory for the resource path does not exist
- **THEN** `SaveResource` SHALL create all necessary parent directories before writing
