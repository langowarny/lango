package adk

import (
	"context"
	"encoding/json"
	"iter"
	"strings"

	"github.com/langowarny/lango/internal/provider"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

type ModelAdapter struct {
	p     provider.Provider
	model string
}

func NewModelAdapter(p provider.Provider, model string) *ModelAdapter {
	return &ModelAdapter{p: p, model: model}
}

func (m *ModelAdapter) Name() string {
	return m.model
}

func (m *ModelAdapter) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		msgs, err := convertMessages(req.Contents)
		if err != nil {
			yield(nil, err)
			return
		}

		tools, err := convertTools(req.Config)
		if err != nil {
			yield(nil, err)
			return
		}

		// Forward ADK system instruction as a system message for the provider.
		if req.Config != nil && req.Config.SystemInstruction != nil {
			sysText := extractSystemText(req.Config.SystemInstruction)
			if sysText != "" {
				sysMsg := provider.Message{Role: "system", Content: sysText}
				msgs = append([]provider.Message{sysMsg}, msgs...)
			}
		}

		params := provider.GenerateParams{
			Model:    req.Model,
			Messages: msgs,
			Tools:    tools,
		}

		if req.Config != nil {
			if req.Config.Temperature != nil {
				params.Temperature = float64(*req.Config.Temperature)
			}
			if req.Config.MaxOutputTokens != 0 {
				params.MaxTokens = int(req.Config.MaxOutputTokens)
			}
		}
		if params.Model == "" {
			// Fallback if not set in request (ADK might set it in client/factory)
			// But params must have it.
			// We can default or error.
			// provider usually requires it.
		}

		pSeq, err := m.p.Generate(ctx, params)
		if err != nil {
			yield(nil, err)
			return
		}

		for evt, err := range pSeq {
			if err != nil {
				yield(nil, err)
				return
			}

			resp := &model.LLMResponse{
				Content: &genai.Content{
					Role:  "model",
					Parts: []*genai.Part{},
				},
			}

			switch evt.Type {
			case provider.StreamEventPlainText:
				resp.Content.Parts = append(resp.Content.Parts, &genai.Part{Text: evt.Text})
				// For streaming, ADK usually expects partial responses?
				// resp.Partial = true (for text delta)
				// But ADK LLMResponse definition comment says:
				// "Only used for streaming mode and when the content is plain text."
				resp.Partial = true
			case provider.StreamEventToolCall:
				if evt.ToolCall != nil {
					args := make(map[string]any)
					_ = json.Unmarshal([]byte(evt.ToolCall.Arguments), &args)
					resp.Content.Parts = append(resp.Content.Parts, &genai.Part{
						FunctionCall: &genai.FunctionCall{
							Name: evt.ToolCall.Name,
							Args: args,
						},
					})
				}
			case provider.StreamEventDone:
				resp.TurnComplete = true
				resp.Partial = false
			case provider.StreamEventError:
				yield(nil, evt.Error)
				return
			}

			if !yield(resp, nil) {
				return
			}
		}
	}
}

func convertMessages(contents []*genai.Content) ([]provider.Message, error) {
	var msgs []provider.Message
	for _, c := range contents {
		role := c.Role
		if role == "model" {
			role = "assistant"
		} else if role == "function" {
			role = "tool"
		}

		msg := provider.Message{Role: role}
		for _, p := range c.Parts {
			if p.Text != "" {
				msg.Content += p.Text
			}
			if p.FunctionCall != nil {
				b, _ := json.Marshal(p.FunctionCall.Args)
				msg.ToolCalls = append(msg.ToolCalls, provider.ToolCall{
					ID:        "call_" + p.FunctionCall.Name, // Synthetic ID if needed
					Name:      p.FunctionCall.Name,
					Arguments: string(b),
				})
			}
			if p.FunctionResponse != nil {
				// Tool response content
				// Should be handled as content/metadata?
				// Provider expects content string usually for tool result.
				// We map it to content.
				b, _ := json.Marshal(p.FunctionResponse.Response)
				msg.Content += string(b)
				// Add tool call ID to metadata if provider supports linking
				if msg.Metadata == nil {
					msg.Metadata = make(map[string]interface{})
				}
				msg.Metadata["tool_call_id"] = p.FunctionResponse.Name // Or ID if available?
			}
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// extractSystemText concatenates all text parts from a genai.Content into a single string.
func extractSystemText(content *genai.Content) string {
	var parts []string
	for _, p := range content.Parts {
		if p.Text != "" {
			parts = append(parts, p.Text)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, "\n")
}

func convertTools(cfg *genai.GenerateContentConfig) ([]provider.Tool, error) {
	var tools []provider.Tool
	if cfg == nil || cfg.Tools == nil {
		return tools, nil
	}

	for _, t := range cfg.Tools {
		if t.FunctionDeclarations != nil {
			for _, fd := range t.FunctionDeclarations {
				// Convert Schema to map
				schemaMap := make(map[string]interface{})
				if fd.Parameters != nil {
					// genai.Schema to map is complex if we recurse.
					// But we can json marshal/unmarshal
					b, err := json.Marshal(fd.Parameters)
					if err == nil {
						_ = json.Unmarshal(b, &schemaMap)
					}
				}

				tools = append(tools, provider.Tool{
					Name:        fd.Name,
					Description: fd.Description,
					Parameters:  schemaMap,
				})
			}
		}
	}
	return tools, nil
}
