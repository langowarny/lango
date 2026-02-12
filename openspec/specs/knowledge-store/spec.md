## ADDED Requirements

### Requirement: Knowledge CRUD Operations
The system SHALL provide persistent CRUD operations for knowledge entries identified by a unique key.

#### Scenario: Save new knowledge entry
- **WHEN** `SaveKnowledge` is called with a key that does not exist
- **THEN** the system SHALL create a new knowledge entry with the given key, category, content, tags, and source

#### Scenario: Update existing knowledge entry
- **WHEN** `SaveKnowledge` is called with a key that already exists
- **THEN** the system SHALL update the existing entry's category, content, tags, and source

#### Scenario: Get knowledge by key
- **WHEN** `GetKnowledge` is called with an existing key
- **THEN** the system SHALL return the knowledge entry
- **AND** if the key does not exist, SHALL return an error

#### Scenario: Delete knowledge by key
- **WHEN** `DeleteKnowledge` is called with an existing key
- **THEN** the system SHALL remove the entry from the store

#### Scenario: Increment knowledge use count
- **WHEN** `IncrementKnowledgeUseCount` is called with a valid key
- **THEN** the system SHALL increment the use count by 1

### Requirement: Knowledge Search
The system SHALL support keyword-based search across knowledge entries.

#### Scenario: Search by query
- **WHEN** `SearchKnowledge` is called with a query string
- **THEN** the system SHALL return entries where the content or key contains the query
- **AND** results SHALL be ordered by relevance score descending
- **AND** results SHALL be limited to the specified limit (default 10)

#### Scenario: Search with category filter
- **WHEN** `SearchKnowledge` is called with a query and a category
- **THEN** the system SHALL return only entries matching both the query and the category

### Requirement: Learning CRUD Operations
The system SHALL provide persistent CRUD operations for learning entries.

#### Scenario: Save new learning
- **WHEN** `SaveLearning` is called with trigger, error pattern, diagnosis, fix, and category
- **THEN** the system SHALL create a new learning entry

#### Scenario: Search learnings
- **WHEN** `SearchLearnings` is called with an error pattern query
- **THEN** the system SHALL return learnings where the error pattern or trigger contains the query
- **AND** results SHALL be ordered by confidence descending

#### Scenario: Boost learning confidence
- **WHEN** `BoostLearningConfidence` is called with a learning ID and success delta
- **THEN** the system SHALL increment success count and recalculate confidence as `success / (success + occurrence)`
- **AND** confidence SHALL NOT drop below 0.1

### Requirement: Skill Persistence
The system SHALL provide persistent CRUD operations for skill entries.

#### Scenario: Save new skill
- **WHEN** `SaveSkill` is called with a skill entry
- **THEN** the system SHALL create the skill with default status `draft`

#### Scenario: Activate skill
- **WHEN** `ActivateSkill` is called with a skill name
- **THEN** the system SHALL set the skill status to `active`

#### Scenario: List active skills
- **WHEN** `ListActiveSkills` is called
- **THEN** the system SHALL return all skills with status `active`

#### Scenario: Increment skill usage
- **WHEN** `IncrementSkillUsage` is called with a skill name and success flag
- **THEN** the system SHALL increment use count by 1 and update last used timestamp
- **AND** if success is true, SHALL also increment success count

### Requirement: Audit Logging
The system SHALL maintain an append-only audit log of knowledge operations.

#### Scenario: Save audit log entry
- **WHEN** `SaveAuditLog` is called with an action, actor, target, and optional details
- **THEN** the system SHALL create an immutable audit log record

### Requirement: External Reference Management
The system SHALL support CRUD operations for external references (docs, APIs, wiki links).

#### Scenario: Save or update external reference
- **WHEN** `SaveExternalRef` is called with a name that does not exist
- **THEN** the system SHALL create the external reference
- **AND** if the name already exists, SHALL update the existing reference

#### Scenario: Search external references
- **WHEN** `SearchExternalRefs` is called with a query
- **THEN** the system SHALL return references where name or summary contains the query

### Requirement: Per-Session Rate Limiting
The system SHALL enforce per-session rate limits for knowledge and learning saves to prevent runaway accumulation.

#### Scenario: Knowledge save limit
- **WHEN** a session has reached the configured maximum knowledge saves
- **THEN** the system SHALL reject further knowledge saves with an error

#### Scenario: Learning save limit
- **WHEN** a session has reached the configured maximum learning saves
- **THEN** the system SHALL reject further learning saves with an error

#### Scenario: Atomic rate limit check
- **WHEN** checking and incrementing the rate limit counter
- **THEN** the system SHALL perform the check and increment atomically within a single lock acquisition to prevent TOCTOU race conditions

#### Scenario: Empty session key bypasses rate limit
- **WHEN** session key is empty
- **THEN** the system SHALL skip rate limit checks

### Requirement: Ent Schema Definitions
The system SHALL define Ent ORM schemas for the 5 knowledge entities.

#### Scenario: Knowledge schema
- **WHEN** the database is migrated
- **THEN** a `Knowledge` table SHALL exist with fields: key (unique), category (enum: rule/definition/preference/fact), content, tags (JSON), source, relevance_score, use_count, created_at, updated_at

#### Scenario: Learning schema
- **WHEN** the database is migrated
- **THEN** a `Learning` table SHALL exist with fields: trigger, error_pattern, diagnosis, fix, category (enum: tool_error/provider_error/user_correction/timeout/permission/general), tags (JSON), confidence, occurrence_count, success_count, created_at, updated_at

#### Scenario: Skill schema
- **WHEN** the database is migrated
- **THEN** a `Skill` table SHALL exist with fields: name (unique), description, skill_type (enum: composite/script/template), definition (JSON), parameters (JSON), status (enum: draft/active/deprecated), created_by, requires_approval, use_count, success_count, last_used_at, created_at, updated_at

#### Scenario: AuditLog schema
- **WHEN** the database is migrated
- **THEN** an `AuditLog` table SHALL exist with fields: session_key, action (enum), actor, target, details (JSON), created_at

#### Scenario: ExternalRef schema
- **WHEN** the database is migrated
- **THEN** an `ExternalRef` table SHALL exist with fields: name (unique), ref_type (enum: documentation/api/wiki/repository/tool), location, summary, metadata (JSON), created_at, updated_at
