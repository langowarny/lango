// Package embedding provides vector embedding generation and storage
// for RAG (Retrieval-Augmented Generation) capabilities.
package embedding

import "context"

// EmbeddingProvider generates vector embeddings from text.
type EmbeddingProvider interface {
	// ID returns the unique identifier for this provider.
	ID() string
	// Embed generates embeddings for the given texts in a single batch.
	Embed(ctx context.Context, texts []string) ([][]float32, error)
	// Dimensions returns the dimensionality of the generated embeddings.
	Dimensions() int
}
