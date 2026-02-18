## MODIFIED Requirements

### Requirement: Skill type validation
The registry SHALL accept skill types `composite`, `script`, `template`, and `instruction`. The `instruction` type SHALL NOT require a non-empty Definition.

#### Scenario: Valid instruction type
- **WHEN** a skill with type `instruction` is created
- **THEN** the registry SHALL accept it even with empty Definition

#### Scenario: Invalid type
- **WHEN** a skill with an unrecognized type is created
- **THEN** the registry SHALL return an error listing all valid types including `instruction`

### Requirement: Skill entry structure
`SkillEntry` SHALL include a `Source` string field to track the import origin URL. The field SHALL be empty for locally created skills.

#### Scenario: Imported skill
- **WHEN** a skill is imported from an external URL
- **THEN** its `Source` field SHALL contain the origin URL

#### Scenario: Local skill
- **WHEN** a skill is created locally
- **THEN** its `Source` field SHALL be empty

## ADDED Requirements

### Requirement: Instruction skill type
The system SHALL support an `instruction` skill type for agent reference documents. Instruction skills store their entire markdown body as `definition["content"]`.

#### Scenario: Parse instruction SKILL.md
- **WHEN** a SKILL.md has no explicit `type` in frontmatter
- **THEN** the parser SHALL default to type `instruction` and store the body as content

#### Scenario: Parse instruction with explicit type
- **WHEN** a SKILL.md has `type: instruction` in frontmatter
- **THEN** the parser SHALL store the body as `definition["content"]`

### Requirement: Instruction skill as tool
Instruction skills SHALL be converted to agent tools. The tool description SHALL use the skill's original description (for agent reasoning). The handler SHALL return the full reference content, source, and description.

#### Scenario: Tool registration
- **WHEN** an instruction skill is loaded
- **THEN** it SHALL be registered as `skill_{name}` with the original description

#### Scenario: Tool invocation
- **WHEN** the agent calls an instruction skill tool
- **THEN** the handler SHALL return content, source, description, and type fields

#### Scenario: Empty description fallback
- **WHEN** an instruction skill has no description
- **THEN** the tool description SHALL default to "Reference guide for {name}"

### Requirement: Instruction skill execution
The executor SHALL handle the `instruction` type by returning the content directly without script execution.

#### Scenario: Execute instruction skill
- **WHEN** the executor is called with an instruction skill
- **THEN** it SHALL return a map with skill name, type, and content

### Requirement: Source field roundtrip in SKILL.md
The parser and renderer SHALL preserve the `source` field in YAML frontmatter through parse/render cycles.

#### Scenario: Source field roundtrip
- **WHEN** a SKILL.md with a `source` field is parsed and re-rendered
- **THEN** the re-parsed entry SHALL have the same `source` value

### Requirement: Instruction skill rendering
The renderer SHALL output instruction skill content as plain markdown (no code block wrapping).

#### Scenario: Render instruction
- **WHEN** an instruction skill is rendered to SKILL.md format
- **THEN** the body SHALL contain the raw content without code block delimiters

### Requirement: BuildInstructionSkill builder
The system SHALL provide a `BuildInstructionSkill` function that creates an instruction SkillEntry with name, description, content, and source.

#### Scenario: Build instruction skill
- **WHEN** `BuildInstructionSkill("name", "desc", "content", "https://source")` is called
- **THEN** it SHALL return a SkillEntry with type=instruction, status=active, createdBy=import

### Requirement: Registry Store accessor
The registry SHALL expose a `Store()` method returning the underlying `SkillStore` interface.

#### Scenario: Store access
- **WHEN** `registry.Store()` is called
- **THEN** it SHALL return the configured SkillStore instance
