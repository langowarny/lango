## Purpose

Documentation accuracy requirements ensuring README.md stays in sync with codebase configuration and feature state.

## Requirements

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

### Requirement: README documents streaming in gateway feature

README.md Features list SHALL describe the Gateway as supporting real-time streaming.

#### Scenario: Gateway feature line updated
- **WHEN** a user reads the Features list in README.md
- **THEN** the Gateway bullet reads "WebSocket/HTTP server with real-time streaming"

### Requirement: README documents observational memory context limit configs

README.md Configuration Reference table SHALL include `observationalMemory.maxReflectionsInContext` and `observationalMemory.maxObservationsInContext` fields matching `ObservationalMemoryConfig` in `internal/config/types.go`.

#### Scenario: Context limit config fields present
- **WHEN** a user reads the Observational Memory config section in README.md
- **THEN** the table contains `observationalMemory.maxReflectionsInContext` (int, default `5`) and `observationalMemory.maxObservationsInContext` (int, default `20`)

### Requirement: README documents embedding cache

README.md Embedding & RAG section SHALL include an Embedding Cache subsection describing in-memory caching with 5-minute TTL and 100-entry limit.

#### Scenario: Embedding cache subsection present
- **WHEN** a user reads the Embedding & RAG section in README.md
- **THEN** there is an "Embedding Cache" heading describing automatic in-memory caching with 5-minute TTL and 100-entry limit

### Requirement: README documents observational memory context limits

README.md Observational Memory section SHALL describe context limits for reflections and observations.

#### Scenario: Context limits bullet present
- **WHEN** a user reads the Observational Memory component list in README.md
- **THEN** there is a "Context Limits" bullet describing default limits of 5 reflections and 20 observations

### Requirement: README documents WebSocket events

README.md SHALL include a WebSocket Events subsection documenting `agent.thinking`, `agent.chunk`, and `agent.done` events with their payloads.

#### Scenario: WebSocket events table present
- **WHEN** a user reads the WebSocket section in README.md
- **THEN** there is a "WebSocket Events" heading with a table listing `agent.thinking`, `agent.chunk`, and `agent.done` events

#### Scenario: Backward compatibility noted
- **WHEN** a user reads the WebSocket Events section
- **THEN** there is a note that clients not handling `agent.chunk` will still receive the full response in the RPC result

### Requirement: Documentation accuracy

Documentation, prompts, and CLI help text SHALL accurately reflect all implemented features including P2P REST API endpoints, CLI flags, and example projects.

#### Scenario: P2P REST API documented
- **WHEN** a user reads the HTTP API documentation
- **THEN** the P2P REST endpoints (`/api/p2p/status`, `/api/p2p/peers`, `/api/p2p/identity`) SHALL be documented with request/response examples

#### Scenario: Secrets --value-hex documented
- **WHEN** a user reads the secrets set CLI documentation
- **THEN** the `--value-hex` flag SHALL be documented with non-interactive usage examples

#### Scenario: P2P trading example discoverable
- **WHEN** a user reads the README
- **THEN** the `examples/p2p-trading/` directory SHALL be referenced in an Examples section
