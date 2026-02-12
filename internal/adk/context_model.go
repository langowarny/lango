package adk

import (
	"context"
	"iter"

	"go.uber.org/zap"
	"google.golang.org/adk/model"
	"google.golang.org/genai"

	"github.com/langowarny/lango/internal/knowledge"
)

// ContextAwareModelAdapter wraps a ModelAdapter with context retrieval.
// Before each LLM call, it retrieves relevant knowledge and injects it
// into the system instruction.
type ContextAwareModelAdapter struct {
	inner      *ModelAdapter
	retriever  *knowledge.ContextRetriever
	basePrompt string
	logger     *zap.SugaredLogger
}

// NewContextAwareModelAdapter creates a context-aware model adapter.
func NewContextAwareModelAdapter(
	inner *ModelAdapter,
	retriever *knowledge.ContextRetriever,
	basePrompt string,
	logger *zap.SugaredLogger,
) *ContextAwareModelAdapter {
	return &ContextAwareModelAdapter{
		inner:      inner,
		retriever:  retriever,
		basePrompt: basePrompt,
		logger:     logger,
	}
}

// Name delegates to the inner adapter.
func (m *ContextAwareModelAdapter) Name() string {
	return m.inner.Name()
}

// GenerateContent retrieves context and injects an augmented system prompt before delegating to the inner adapter.
func (m *ContextAwareModelAdapter) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	// Extract last user message for context retrieval
	userQuery := extractLastUserMessage(req.Contents)

	if userQuery != "" && m.retriever != nil {
		retrieved, err := m.retriever.Retrieve(ctx, knowledge.RetrievalRequest{
			Query: userQuery,
		})
		if err != nil {
			m.logger.Warnw("context retrieval error", "error", err)
		} else if retrieved != nil && retrieved.TotalItems > 0 {
			augmented := m.retriever.AssemblePrompt(m.basePrompt, retrieved)
			if req.Config == nil {
				req.Config = &genai.GenerateContentConfig{}
			}
			req.Config.SystemInstruction = &genai.Content{
				Parts: []*genai.Part{{Text: augmented}},
			}
		}
	}

	return m.inner.GenerateContent(ctx, req, stream)
}

// extractLastUserMessage finds the last user message from the content history.
func extractLastUserMessage(contents []*genai.Content) string {
	for i := len(contents) - 1; i >= 0; i-- {
		c := contents[i]
		if c.Role == "user" {
			for _, p := range c.Parts {
				if p.Text != "" {
					return p.Text
				}
			}
		}
	}
	return ""
}
