## ADDED Requirements

### Requirement: README documents librarian configuration

README.md Configuration Reference table SHALL include all `librarian.*` fields matching `LibrarianConfig` in `internal/config/types.go`.

#### Scenario: Librarian config fields present
- **WHEN** a user reads the Configuration Reference in README.md
- **THEN** the table contains entries for `librarian.enabled`, `librarian.observationThreshold`, `librarian.inquiryCooldownTurns`, `librarian.maxPendingInquiries`, `librarian.autoSaveConfidence`, `librarian.provider`, `librarian.model`

### Requirement: README documents automation defaultDeliverTo

README.md Configuration Reference table SHALL include `defaultDeliverTo` fields for cron, background, and workflow sections.

#### Scenario: defaultDeliverTo fields present
- **WHEN** a user reads the Cron Scheduling, Background Execution, and Workflow Engine config sections
- **THEN** each section contains a `*.defaultDeliverTo` entry with type `[]string` and default `[]`

### Requirement: README multi-agent table reflects librarian tools

The multi-agent orchestration table SHALL list proactive knowledge extraction in the librarian role and include `librarian_pending_inquiries` and `librarian_dismiss_inquiry` in the tools column.

#### Scenario: Librarian row updated
- **WHEN** a user reads the Multi-Agent Orchestration table
- **THEN** the librarian row includes "proactive knowledge extraction" in Role and both `librarian_pending_inquiries` and `librarian_dismiss_inquiry` in Tools
