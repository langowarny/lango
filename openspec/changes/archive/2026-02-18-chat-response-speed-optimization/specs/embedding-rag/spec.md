## MODIFIED Requirements

### Requirement: Cross-collection semantic search
The RAGService SHALL search all configured collections in parallel using errgroup. Each collection search SHALL run concurrently, with results merged, sorted by distance, and limited after all collections complete. Individual collection search errors SHALL be logged and treated as non-fatal.

#### Scenario: Parallel collection search
- **WHEN** a query is submitted against multiple collections
- **THEN** all collections SHALL be searched concurrently and results merged after all complete

#### Scenario: Single collection failure
- **WHEN** one collection search fails during parallel execution
- **THEN** the error SHALL be logged as a warning and results from other collections SHALL still be returned

#### Scenario: Results sorted and limited
- **WHEN** parallel searches complete
- **THEN** results SHALL be sorted by ascending distance and limited to the configured maximum
