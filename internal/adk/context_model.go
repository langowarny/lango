package adk

import (
	"context"
	"fmt"
	"iter"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/adk/model"
	"google.golang.org/genai"

	"github.com/langowarny/lango/internal/embedding"
	"github.com/langowarny/lango/internal/knowledge"
	"github.com/langowarny/lango/internal/memory"
	"github.com/langowarny/lango/internal/prompt"
)

// MemoryProvider retrieves observations and reflections for a session.
type MemoryProvider interface {
	ListObservations(ctx context.Context, sessionKey string) ([]memory.Observation, error)
	ListReflections(ctx context.Context, sessionKey string) ([]memory.Reflection, error)
}

// ContextAwareModelAdapter wraps a ModelAdapter with context retrieval.
// Before each LLM call, it retrieves relevant knowledge and injects it
// into the system instruction.
type ContextAwareModelAdapter struct {
	inner          *ModelAdapter
	retriever      *knowledge.ContextRetriever
	memoryProvider MemoryProvider
	ragService     *embedding.RAGService
	ragOpts        embedding.RetrieveOptions
	runtimeAdapter *RuntimeContextAdapter
	sessionKey     string
	basePrompt     string
	logger         *zap.SugaredLogger
}

// NewContextAwareModelAdapter creates a context-aware model adapter.
// The builder is used to produce the base system prompt; dynamic context
// (knowledge, memory, RAG) is still appended at call time.
func NewContextAwareModelAdapter(
	inner *ModelAdapter,
	retriever *knowledge.ContextRetriever,
	builder *prompt.Builder,
	logger *zap.SugaredLogger,
) *ContextAwareModelAdapter {
	return &ContextAwareModelAdapter{
		inner:      inner,
		retriever:  retriever,
		basePrompt: builder.Build(),
		logger:     logger,
	}
}

// WithMemory adds observational memory support to the adapter.
func (m *ContextAwareModelAdapter) WithMemory(provider MemoryProvider, sessionKey string) *ContextAwareModelAdapter {
	m.memoryProvider = provider
	m.sessionKey = sessionKey
	return m
}

// WithRuntimeAdapter adds runtime context support to the adapter.
func (m *ContextAwareModelAdapter) WithRuntimeAdapter(adapter *RuntimeContextAdapter) *ContextAwareModelAdapter {
	m.runtimeAdapter = adapter
	return m
}

// WithRAG adds RAG (retrieval-augmented generation) support.
func (m *ContextAwareModelAdapter) WithRAG(svc *embedding.RAGService, opts embedding.RetrieveOptions) *ContextAwareModelAdapter {
	m.ragService = svc
	m.ragOpts = opts
	return m
}

// Name delegates to the inner adapter.
func (m *ContextAwareModelAdapter) Name() string {
	return m.inner.Name()
}

// GenerateContent retrieves context and injects an augmented system prompt before delegating to the inner adapter.
func (m *ContextAwareModelAdapter) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	prompt := m.basePrompt

	// Update runtime session state before retrieval
	if m.runtimeAdapter != nil && m.sessionKey != "" {
		m.runtimeAdapter.SetSession(m.sessionKey)
	}

	// Retrieve knowledge context
	userQuery := extractLastUserMessage(req.Contents)
	if userQuery != "" && m.retriever != nil {
		layers := []knowledge.ContextLayer{
			knowledge.LayerRuntimeContext,
			knowledge.LayerToolRegistry,
			knowledge.LayerUserKnowledge,
			knowledge.LayerSkillPatterns,
			knowledge.LayerExternalKnowledge,
			knowledge.LayerAgentLearnings,
		}
		retrieved, err := m.retriever.Retrieve(ctx, knowledge.RetrievalRequest{
			Query:  userQuery,
			Layers: layers,
		})
		if err != nil {
			m.logger.Warnw("context retrieval error", "error", err)
		} else if retrieved != nil && retrieved.TotalItems > 0 {
			prompt = m.retriever.AssemblePrompt(prompt, retrieved)
		}
	}

	// Retrieve RAG context (semantic search across all collections).
	if m.ragService != nil && userQuery != "" {
		ragSection := m.assembleRAGSection(ctx, userQuery)
		if ragSection != "" {
			prompt = fmt.Sprintf("%s\n\n%s", prompt, ragSection)
		}
	}

	// Retrieve observational memory
	if m.memoryProvider != nil && m.sessionKey != "" {
		memorySection := m.assembleMemorySection(ctx)
		if memorySection != "" {
			prompt = fmt.Sprintf("%s\n\n%s", prompt, memorySection)
		}
	}

	// Set the augmented system instruction
	if prompt != m.basePrompt {
		if req.Config == nil {
			req.Config = &genai.GenerateContentConfig{}
		}
		req.Config.SystemInstruction = &genai.Content{
			Parts: []*genai.Part{{Text: prompt}},
		}
	}

	return m.inner.GenerateContent(ctx, req, stream)
}

// assembleMemorySection builds the "Conversation Memory" section from observations and reflections.
func (m *ContextAwareModelAdapter) assembleMemorySection(ctx context.Context) string {
	reflections, err := m.memoryProvider.ListReflections(ctx, m.sessionKey)
	if err != nil {
		m.logger.Warnw("memory reflection retrieval error", "error", err)
	}

	observations, err := m.memoryProvider.ListObservations(ctx, m.sessionKey)
	if err != nil {
		m.logger.Warnw("memory observation retrieval error", "error", err)
	}

	if len(reflections) == 0 && len(observations) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("## Conversation Memory\n")

	if len(reflections) > 0 {
		b.WriteString("\n### Summary\n")
		for _, ref := range reflections {
			b.WriteString(ref.Content)
			b.WriteString("\n")
		}
	}

	if len(observations) > 0 {
		b.WriteString("\n### Recent Observations\n")
		for _, obs := range observations {
			b.WriteString("- ")
			b.WriteString(obs.Content)
			b.WriteString("\n")
		}
	}

	return b.String()
}

// assembleRAGSection builds a "Semantic Context" section from RAG retrieval results.
func (m *ContextAwareModelAdapter) assembleRAGSection(ctx context.Context, query string) string {
	results, err := m.ragService.Retrieve(ctx, query, m.ragOpts)
	if err != nil {
		m.logger.Warnw("rag retrieval error", "error", err)
		return ""
	}
	if len(results) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("## Semantic Context (RAG)\n")
	for _, r := range results {
		if r.Content == "" {
			continue
		}
		b.WriteString(fmt.Sprintf("\n### [%s] %s\n", r.Collection, r.SourceID))
		b.WriteString(r.Content)
		b.WriteString("\n")
	}
	return b.String()
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
