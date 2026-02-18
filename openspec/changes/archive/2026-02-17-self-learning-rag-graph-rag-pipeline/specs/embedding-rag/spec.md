## MODIFIED Requirements

### Requirement: RAG filters results by maximum distance
The system SHALL support a MaxDistance configuration (default 0.0 = disabled). When enabled, vector search results with distance exceeding MaxDistance SHALL be excluded from RAG context.

#### Scenario: MaxDistance filters irrelevant results
- **WHEN** MaxDistance is set to 0.5 and a search result has distance 0.7
- **THEN** that result SHALL be excluded from the returned results

#### Scenario: MaxDistance disabled by default
- **WHEN** MaxDistance is 0.0 (default)
- **THEN** all results SHALL be returned regardless of distance (backward compatible)

### Requirement: RAG supports session-scoped metadata filtering
The system SHALL support filtering vector search results by metadata key-value pairs, enabling session-scoped retrieval.

#### Scenario: Filter by session key
- **WHEN** a RAG query includes a session key
- **THEN** results SHALL be filtered to include only entries matching that session's metadata

#### Scenario: No filter when session key is empty
- **WHEN** a RAG query has no session key
- **THEN** all results SHALL be returned without metadata filtering (backward compatible)

### Requirement: VectorStore Search accepts optional SearchOptions
The VectorStore.Search method SHALL accept an optional `*SearchOptions` parameter for metadata filtering. Nil means no filtering.

#### Scenario: Nil options preserves existing behavior
- **WHEN** Search is called with nil SearchOptions
- **THEN** behavior SHALL be identical to the current implementation

#### Scenario: MetadataFilter post-filters results
- **WHEN** Search is called with MetadataFilter containing key-value pairs
- **THEN** results SHALL be post-filtered to match all specified metadata pairs, with 3x over-fetch to compensate
