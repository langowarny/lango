## ADDED Requirements

### Requirement: Two-phase hybrid retrieval
The GraphRAGService SHALL perform 2-phase retrieval: Phase 1 vector search (sqlite-vec cosine similarity), Phase 2 graph expansion (BFS traversal from Phase 1 results).

#### Scenario: Vector results expanded via graph
- **WHEN** a query returns vector results with source IDs matching graph nodes
- **THEN** the service SHALL traverse the graph from each result node and append discovered nodes as `GraphNode` entries

#### Scenario: No graph store available
- **WHEN** the graph store is nil or vector results are empty
- **THEN** the service SHALL return vector results only without graph expansion

#### Scenario: Expansion limit respected
- **WHEN** graph traversal yields more nodes than `maxExpand`
- **THEN** the service SHALL stop expanding and return at most `maxExpand` graph results

### Requirement: LLM-based entity extraction
The system SHALL use an LLM to extract entities and relationships from saved knowledge and memory content, producing triples for the graph store.

#### Scenario: Entity extraction on knowledge save
- **WHEN** a knowledge entry is saved with `graph.enabled: true`
- **THEN** an async goroutine SHALL extract entities via the Extractor and enqueue resulting triples to GraphBuffer

#### Scenario: No meaningful relationships
- **WHEN** the LLM returns "NONE" for a content piece
- **THEN** no triples SHALL be enqueued (only the basic Contains triple from the direct callback)

### Requirement: Context injection for Graph RAG
The ContextAwareModelAdapter SHALL inject Graph RAG results into the system prompt when both graph store and RAG are enabled.

#### Scenario: Graph RAG section in system prompt
- **WHEN** a query triggers RAG retrieval with graph enabled
- **THEN** the system prompt SHALL include both "Semantic Context (RAG)" and "Graph-Expanded Context" sections
