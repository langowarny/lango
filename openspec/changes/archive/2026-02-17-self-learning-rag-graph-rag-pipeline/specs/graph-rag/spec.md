## MODIFIED Requirements

### Requirement: VectorRetrieveOptions supports MaxDistance
VectorRetrieveOptions SHALL include a MaxDistance field that is passed through to the underlying vector retrieval.

#### Scenario: MaxDistance passed to vector retriever
- **WHEN** Graph RAG retrieval is invoked with MaxDistance set
- **THEN** the MaxDistance value SHALL be forwarded to the VectorRetriever's RetrieveOptions
