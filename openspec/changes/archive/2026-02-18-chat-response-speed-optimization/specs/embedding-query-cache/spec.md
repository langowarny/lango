## ADDED Requirements

### Requirement: TTL-based embedding vector cache
The RAGService SHALL maintain an in-memory cache of query embedding vectors to avoid redundant embedding API calls. The cache SHALL use a TTL of 5 minutes and a maximum size of 100 entries.

#### Scenario: Cache hit returns cached vector
- **WHEN** a query is submitted that was previously embedded within the TTL window
- **THEN** the cached embedding vector SHALL be returned without calling the embedding API

#### Scenario: Cache miss triggers API call
- **WHEN** a query is submitted that is not in the cache
- **THEN** the embedding API SHALL be called and the result cached

#### Scenario: Expired entry triggers re-embedding
- **WHEN** a cached entry's TTL has expired
- **THEN** the entry SHALL be evicted and the embedding API SHALL be called for a fresh vector

### Requirement: Cache eviction at capacity
The cache SHALL evict expired entries first when at capacity. If still at capacity after expired eviction, it SHALL evict the oldest entry.

#### Scenario: Eviction of expired entries at capacity
- **WHEN** the cache reaches maxSize and new entry is added
- **THEN** expired entries SHALL be removed first

#### Scenario: Eviction of oldest entry when no expired entries
- **WHEN** the cache is at capacity with no expired entries
- **THEN** the oldest entry by creation time SHALL be evicted
