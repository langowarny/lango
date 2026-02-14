## ADDED Requirements

### Requirement: Knowledge configuration exposed in TUI
The Onboard TUI SHALL provide a dedicated Knowledge configuration form accessible from the main menu.

#### Scenario: Knowledge menu category
- **WHEN** user views the Configuration Menu in the onboard wizard
- **THEN** a "Knowledge" category SHALL appear with label "🧠 Knowledge" and description "Learning, Skills, Context limits"

#### Scenario: Knowledge form fields
- **WHEN** user selects the Knowledge category
- **THEN** the form SHALL display 6 fields:
  - knowledge_enabled (boolean toggle)
  - knowledge_max_learnings (integer input)
  - knowledge_max_knowledge (integer input)
  - knowledge_max_context (integer input for MaxContextPerLayer)
  - knowledge_auto_approve (boolean toggle for AutoApproveSkills)
  - knowledge_max_skills_day (integer input for MaxSkillsPerDay)

#### Scenario: Knowledge config persistence
- **WHEN** user modifies Knowledge form fields and saves
- **THEN** values SHALL be written to the `knowledge` section of `lango.json`
