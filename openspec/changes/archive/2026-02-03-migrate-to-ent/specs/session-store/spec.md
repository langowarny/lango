## MODIFIED Requirements

### Requirement: Session storage implementation (MODIFIED)
The session store implementation SHALL use entgo.io instead of raw SQL queries.

#### Scenario: Create session (unchanged interface)
- **WHEN** `Store.Create(session)` is called
- **THEN** the session SHALL be persisted using ent client

#### Scenario: Get session (unchanged interface)
- **WHEN** `Store.Get(key)` is called
- **THEN** the session SHALL be retrieved using ent query with Message eager loading

#### Scenario: AppendMessage (unchanged interface)
- **WHEN** `Store.AppendMessage(key, msg)` is called
- **THEN** a new Message entity SHALL be created linked to the Session
