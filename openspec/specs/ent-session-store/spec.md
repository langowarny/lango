## ADDED Requirements

### Requirement: Ent schema definition
The system SHALL define ent schemas for Session and Message entities with appropriate fields and relationships.

#### Scenario: Session schema generation
- **WHEN** `go generate ./ent` is run
- **THEN** the system SHALL generate type-safe Go code for Session CRUD operations

#### Scenario: Message relationship
- **WHEN** a Session is queried
- **THEN** the system SHALL be able to eagerly load associated Messages

### Requirement: Pure Go SQLite driver
The system SHALL use `modernc.org/sqlite` as the SQLite driver to eliminate CGO dependency.

#### Scenario: Building without CGO
- **WHEN** the project is built with `CGO_ENABLED=0`
- **THEN** the build SHALL succeed without errors

### Requirement: Store interface compatibility
The new ent-based store SHALL implement the existing `Store` interface.

#### Scenario: Interface compliance
- **WHEN** `EntStore` is used
- **THEN** it SHALL implement Create, Get, Update, Delete, AppendMessage, and Close methods

### Requirement: Automatic schema migration
The system SHALL automatically migrate the database schema on startup.

#### Scenario: New database
- **WHEN** the application starts with a new database file
- **THEN** ent SHALL create the required tables automatically

#### Scenario: Schema update
- **WHEN** the ent schema is modified and app restarts
- **THEN** ent SHALL apply additive migrations without data loss
