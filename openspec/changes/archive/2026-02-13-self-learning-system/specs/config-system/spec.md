## ADDED Requirements

### Requirement: Knowledge Configuration Section
The system SHALL support a `knowledge` section in the configuration for self-learning settings.

#### Scenario: Knowledge config fields
- **WHEN** `knowledge` section is present in configuration
- **THEN** it SHALL support the following fields:
  - `enabled` (bool): Enable the knowledge/learning system (default: false)
  - `maxLearnings` (int): Maximum learning entries per session (default: 10)
  - `maxKnowledge` (int): Maximum knowledge entries per session (default: 20)
  - `maxContextPerLayer` (int): Maximum context items per layer in retrieval (default: 5)
  - `autoApproveSkills` (bool): Auto-approve new skills without human review (default: false)
  - `maxSkillsPerDay` (int): Maximum new skills per day

#### Scenario: Knowledge disabled by default
- **WHEN** `knowledge` section is omitted from configuration
- **THEN** the system SHALL treat knowledge as disabled
- **AND** no knowledge-related initialization SHALL occur

#### Scenario: Knowledge config validation
- **WHEN** `knowledge.enabled` is true
- **THEN** the system SHALL apply default values for any omitted numeric fields
- **AND** `maxLearnings` SHALL default to 10 if not specified or <= 0
- **AND** `maxKnowledge` SHALL default to 20 if not specified or <= 0
- **AND** `maxContextPerLayer` SHALL default to 5 if not specified or <= 0
