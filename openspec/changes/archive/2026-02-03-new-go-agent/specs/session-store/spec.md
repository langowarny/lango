## ADDED Requirements

### Requirement: Session creation
The system SHALL create new sessions with unique identifiers and store them in SQLite.

#### Scenario: Create new session
- **WHEN** a new conversation begins
- **THEN** a session record SHALL be created with a unique key

#### Scenario: Session with agent assignment
- **WHEN** creating a session for a specific agent
- **THEN** the agent ID SHALL be associated with the session

### Requirement: Message history storage
The system SHALL store conversation message history in the session.

#### Scenario: Store user message
- **WHEN** a user message is processed
- **THEN** the message SHALL be appended to the session history

#### Scenario: Store assistant response
- **WHEN** the assistant generates a response
- **THEN** the response SHALL be appended to the session history

### Requirement: Session retrieval
The system SHALL retrieve session data including message history.

#### Scenario: Load session by key
- **WHEN** a session key is provided
- **THEN** the full session data SHALL be loaded

#### Scenario: Session not found
- **WHEN** an invalid session key is provided
- **THEN** a session-not-found error SHALL be returned

### Requirement: Session metadata
The system SHALL store and retrieve session metadata (model, settings).

#### Scenario: Store session settings
- **WHEN** session settings are updated (model, thinking level)
- **THEN** the settings SHALL be persisted

#### Scenario: Retrieve session settings
- **WHEN** a session is loaded
- **THEN** the current settings SHALL be included

### Requirement: Session cleanup
The system SHALL support session deletion and expiration.

#### Scenario: Delete session
- **WHEN** session deletion is requested
- **THEN** all session data SHALL be removed from storage

#### Scenario: Session expiration
- **WHEN** a session exceeds its TTL
- **THEN** the session MAY be marked for cleanup
