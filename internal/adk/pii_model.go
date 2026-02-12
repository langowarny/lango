package adk

import (
	"context"
	"iter"

	"github.com/langowarny/lango/internal/agent"
	"github.com/langowarny/lango/internal/logging"
	"google.golang.org/adk/model"
)

var piiModelLogger = logging.SubsystemSugar("pii-model")

// PIIRedactingModelAdapter wraps an LLM and redacts PII from user messages
// before forwarding them to the underlying model.
type PIIRedactingModelAdapter struct {
	inner    model.LLM
	redactor *agent.PIIRedactor
}

// NewPIIRedactingModelAdapter creates a PII-redacting model adapter.
func NewPIIRedactingModelAdapter(inner model.LLM, redactor *agent.PIIRedactor) *PIIRedactingModelAdapter {
	return &PIIRedactingModelAdapter{
		inner:    inner,
		redactor: redactor,
	}
}

// Name delegates to the inner adapter.
func (m *PIIRedactingModelAdapter) Name() string {
	return m.inner.Name()
}

// GenerateContent redacts PII from user messages before delegating to the inner model.
func (m *PIIRedactingModelAdapter) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	// Redact PII from user messages
	for _, content := range req.Contents {
		if content.Role != "user" {
			continue
		}
		for _, part := range content.Parts {
			if part.Text == "" {
				continue
			}
			redacted := m.redactor.RedactInput(part.Text)
			if redacted != part.Text {
				piiModelLogger.Debugw("redacted PII from user message")
				part.Text = redacted
			}
		}
	}

	return m.inner.GenerateContent(ctx, req, stream)
}
