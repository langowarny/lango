package embedding

import "context"

// VectorRecord represents a single vector entry in the store.
type VectorRecord struct {
	// ID is the source entity identifier.
	ID string
	// Collection groups records (e.g., "knowledge", "observation", "reflection", "learning").
	Collection string
	// Embedding is the vector representation.
	Embedding []float32
	// Metadata holds filterable key-value pairs (e.g., session_key, category).
	Metadata map[string]string
}

// SearchResult represents a single result from a vector similarity search.
type SearchResult struct {
	ID         string
	Collection string
	Distance   float32
	Metadata   map[string]string
}

// VectorStore provides vector storage and similarity search.
type VectorStore interface {
	// Upsert inserts or replaces vector records.
	Upsert(ctx context.Context, records []VectorRecord) error
	// Search finds the most similar vectors in a collection.
	Search(ctx context.Context, collection string, query []float32, limit int) ([]SearchResult, error)
	// Delete removes vectors by collection and source IDs.
	Delete(ctx context.Context, collection string, ids []string) error
	// Close releases any resources held by the store.
	Close() error
}
