package embedding

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// RAGResult represents a single retrieval result with original content.
type RAGResult struct {
	Collection string
	SourceID   string
	Content    string
	Distance   float32
}

// RetrieveOptions configures a RAG retrieval query.
type RetrieveOptions struct {
	// Collections to search (empty means all).
	Collections []string
	// Maximum results to return.
	Limit int
	// SessionKey filters for session-specific results.
	SessionKey string
	// MaxDistance is the maximum cosine distance for results (0.0 = disabled).
	MaxDistance float32
}

// ContentResolver looks up the original text content by collection and ID.
type ContentResolver interface {
	ResolveContent(ctx context.Context, collection, id string) (string, error)
}

// RAGService provides semantic retrieval across all embedded collections.
type RAGService struct {
	provider EmbeddingProvider
	store    VectorStore
	resolver ContentResolver
	logger   *zap.SugaredLogger
}

// NewRAGService creates a new RAG retrieval service.
func NewRAGService(
	provider EmbeddingProvider,
	store VectorStore,
	resolver ContentResolver,
	logger *zap.SugaredLogger,
) *RAGService {
	return &RAGService{
		provider: provider,
		store:    store,
		resolver: resolver,
		logger:   logger,
	}
}

// allCollections lists all supported embedding collections.
var allCollections = []string{"knowledge", "observation", "reflection", "learning"}

// Retrieve finds relevant context across collections for a given query.
func (r *RAGService) Retrieve(ctx context.Context, query string, opts RetrieveOptions) ([]RAGResult, error) {
	if query == "" {
		return nil, nil
	}

	if opts.Limit <= 0 {
		opts.Limit = 5
	}

	// Embed the query text.
	embeddings, err := r.provider.Embed(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}
	if len(embeddings) == 0 {
		return nil, nil
	}
	queryVec := embeddings[0]

	collections := opts.Collections
	if len(collections) == 0 {
		collections = allCollections
	}

	// Search each collection and merge results.
	var results []RAGResult
	perCollectionLimit := opts.Limit
	if len(collections) > 1 {
		// Fetch more per collection to allow cross-collection ranking.
		perCollectionLimit = opts.Limit * 2
	}

	// Build search options from session key.
	var searchOpts *SearchOptions
	if opts.SessionKey != "" {
		searchOpts = &SearchOptions{
			MetadataFilter: map[string]string{"session_key": opts.SessionKey},
		}
	}

	for _, col := range collections {
		hits, err := r.store.Search(ctx, col, queryVec, perCollectionLimit, searchOpts)
		if err != nil {
			r.logger.Warnw("rag search error", "collection", col, "error", err)
			continue
		}

		for _, hit := range hits {
			content := ""
			if r.resolver != nil {
				resolved, err := r.resolver.ResolveContent(ctx, col, hit.ID)
				if err != nil {
					r.logger.Debugw("content resolve failed", "collection", col, "id", hit.ID, "error", err)
					continue
				}
				content = resolved
			}

			results = append(results, RAGResult{
				Collection: col,
				SourceID:   hit.ID,
				Content:    content,
				Distance:   hit.Distance,
			})
		}
	}

	// Sort by distance and limit.
	sortByDistance(results)
	if len(results) > opts.Limit {
		results = results[:opts.Limit]
	}

	// Filter by MaxDistance if configured.
	if opts.MaxDistance > 0 {
		results = filterByMaxDistance(results, opts.MaxDistance)
	}

	return results, nil
}

// sortByDistance sorts results by ascending distance (most similar first).
func sortByDistance(results []RAGResult) {
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Distance < results[j-1].Distance; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
}

// filterByMaxDistance removes results whose distance exceeds maxDist.
func filterByMaxDistance(results []RAGResult, maxDist float32) []RAGResult {
	filtered := make([]RAGResult, 0, len(results))
	for _, r := range results {
		if r.Distance <= maxDist {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
