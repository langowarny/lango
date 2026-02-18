## MODIFIED Requirements

### Requirement: Parallel context retrieval
The ContextAwareModelAdapter GenerateContent method SHALL execute knowledge retrieval, RAG/GraphRAG retrieval, and memory retrieval in parallel using errgroup. Each retrieval SHALL run as a separate goroutine with context cancellation propagation. Errors in individual retrievals SHALL be logged and treated as non-fatal (existing degradation pattern preserved).

#### Scenario: All three retrievals run concurrently
- **WHEN** GenerateContent is called with knowledge, RAG, and memory providers configured
- **THEN** all three retrievals SHALL execute concurrently and their results combined after completion

#### Scenario: One retrieval fails
- **WHEN** knowledge retrieval fails but RAG and memory succeed
- **THEN** the error SHALL be logged and the prompt SHALL include RAG and memory sections only

#### Scenario: Memory limits applied during parallel retrieval
- **WHEN** maxReflections and maxObservations are configured
- **THEN** assembleMemorySection SHALL use ListRecentReflections/ListRecentObservations with the configured limits
