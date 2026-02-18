## Purpose

Allowed tools frontmatter support — enables SKILL.md files to declare pre-approved tools via the `allowed-tools` YAML frontmatter field.

## Requirements

### Requirement: AllowedTools field in SkillEntry
`SkillEntry` SHALL include an `AllowedTools []string` field that holds pre-approved tool names parsed from the SKILL.md frontmatter.

#### Scenario: Parse allowed-tools from frontmatter
- **WHEN** a SKILL.md contains `allowed-tools: exec fs_read browser_navigate` in its YAML frontmatter
- **THEN** `ParseSkillMD` SHALL set `AllowedTools` to `["exec", "fs_read", "browser_navigate"]`

#### Scenario: No allowed-tools in frontmatter
- **WHEN** a SKILL.md does not contain an `allowed-tools` field
- **THEN** `AllowedTools` SHALL be nil or empty

### Requirement: AllowedTools serialization in RenderSkillMD
`RenderSkillMD` SHALL serialize the `AllowedTools` field as a space-separated string in the `allowed-tools` YAML frontmatter key.

#### Scenario: Render allowed-tools roundtrip
- **WHEN** a `SkillEntry` with `AllowedTools: ["exec", "fs_read"]` is rendered and re-parsed
- **THEN** the re-parsed entry SHALL have `AllowedTools: ["exec", "fs_read"]`

#### Scenario: Render without allowed-tools
- **WHEN** a `SkillEntry` has empty `AllowedTools`
- **THEN** the rendered SKILL.md SHALL NOT include the `allowed-tools` frontmatter key
