## ADDED Requirements

### Requirement: BoltDB-backed triple store
The system SHALL provide a graph store backed by BoltDB that stores Subject-Predicate-Object triples with SPO, POS, and OSP index buckets for efficient lookups.

#### Scenario: Add and query a triple
- **WHEN** a triple `{Subject: "error:timeout", Predicate: "caused_by", Object: "tool:exec"}` is added
- **THEN** the triple SHALL be retrievable by subject, by object, and by subject+predicate queries

#### Scenario: Batch add triples atomically
- **WHEN** multiple triples are submitted via `AddTriples`
- **THEN** all triples SHALL be persisted in a single BoltDB transaction

#### Scenario: BFS traversal
- **WHEN** `Traverse("node-A", maxDepth=2, predicates=["related_to"])` is called
- **THEN** the store SHALL return all triples reachable within 2 hops following only `related_to` edges

### Requirement: Predicate vocabulary
The graph store SHALL define these standard predicates: `related_to`, `caused_by`, `resolved_by`, `follows`, `similar_to`, `contains`, `in_session`, `reflects_on`, `learned_from`.

#### Scenario: Invalid predicate rejected by extractor
- **WHEN** the entity extractor produces a triple with predicate `"foo_bar"`
- **THEN** the triple SHALL be skipped during parsing

### Requirement: Async graph buffer
The system SHALL provide a `GraphBuffer` that batches graph update requests and processes them on a background goroutine with Start/Enqueue/Stop lifecycle.

#### Scenario: GraphBuffer lifecycle
- **WHEN** the application starts with `graph.enabled: true`
- **THEN** GraphBuffer.Start() SHALL be called, and on shutdown GraphBuffer.Stop() SHALL drain remaining items before exiting

#### Scenario: Non-blocking enqueue
- **WHEN** the buffer queue is full (capacity 256)
- **THEN** Enqueue SHALL drop the request without blocking

### Requirement: Graph store close on shutdown
The system SHALL call `GraphStore.Close()` during application shutdown after the WaitGroup completes.

#### Scenario: Clean shutdown
- **WHEN** the application is stopping
- **THEN** GraphBuffer.Stop() SHALL be called before wg.Wait(), and GraphStore.Close() SHALL be called after wg.Wait()
