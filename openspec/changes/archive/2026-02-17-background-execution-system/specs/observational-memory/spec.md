## MODIFIED Requirements

### Requirement: Memory buffer compaction
The observational memory buffer SHALL support a compaction callback that deletes observed messages and replaces them with a summary after successful observation.

#### Scenario: Compaction after observation
- **WHEN** the memory buffer processes messages and generates an observation
- **THEN** the buffer SHALL invoke the compactor function to replace observed messages with a summary

#### Scenario: Compactor not set
- **WHEN** the buffer processes messages but no compactor is configured
- **THEN** the buffer SHALL skip compaction and retain all original messages

#### Scenario: Index reset after compaction
- **WHEN** compaction deletes messages up to index N and inserts a summary
- **THEN** the buffer SHALL reset its lastObserved index to 0 since message indices have shifted

## ADDED Requirements

### Requirement: SetCompactor configuration
The memory buffer SHALL provide a SetCompactor() method to configure the compaction callback at runtime.

#### Scenario: Wire compactor during app initialization
- **WHEN** the app wires memory buffer with session store
- **THEN** SetCompactor SHALL be called with EntStore.CompactMessages as the compaction function

### Requirement: Session store CompactMessages
The session store SHALL provide a CompactMessages(key, upToIndex, summary) method that atomically deletes messages up to the given index and inserts a summary message.

#### Scenario: Compact messages
- **WHEN** CompactMessages is called with key "session:123", upToIndex 5, and summary "User discussed weather"
- **THEN** messages at indices 0-5 SHALL be deleted and a new message with the summary SHALL be inserted with an early timestamp

#### Scenario: Compact with no messages to delete
- **WHEN** CompactMessages is called with upToIndex 0
- **THEN** the operation SHALL still insert the summary message
