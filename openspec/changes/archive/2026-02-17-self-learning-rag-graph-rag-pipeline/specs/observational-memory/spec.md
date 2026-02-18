## MODIFIED Requirements

### Requirement: Buffer drops logged at warn level with counters
EmbeddingBuffer and GraphBuffer SHALL log dropped requests at warn level (not debug) and track drop counts via atomic counters accessible through a DroppedCount() method.

#### Scenario: Queue full logs warning
- **WHEN** a buffer's queue is full and a new request arrives
- **THEN** the request SHALL be dropped with a warn-level log entry including the request ID

#### Scenario: Drop counter increments
- **WHEN** a buffer drops a request
- **THEN** the atomic drop counter SHALL increment and be readable via DroppedCount()
