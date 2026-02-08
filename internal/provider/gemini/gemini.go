package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"

	"github.com/langowarny/lango/internal/provider"
	"google.golang.org/genai"
)

type GeminiProvider struct {
	client *genai.Client
	id     string
	model  string
}

func NewProvider(ctx context.Context, apiKey string, model string) (*GeminiProvider, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	return &GeminiProvider{
		client: client,
		id:     "gemini",
		model:  model,
	}, nil
}

func (p *GeminiProvider) ID() string {
	return p.id
}

func (p *GeminiProvider) Generate(ctx context.Context, params provider.GenerateParams) (iter.Seq2[provider.StreamEvent, error], error) {
	// Convert messages to genai.Content
	var contents []*genai.Content
	for _, m := range params.Messages {
		contents = append(contents, &genai.Content{
			Role:  m.Role,
			Parts: []*genai.Part{{Text: m.Content}},
		})
	}

	// Tools
	var tools []*genai.Tool
	if len(params.Tools) > 0 {
		var funcDecls []*genai.FunctionDeclaration
		for _, t := range params.Tools {
			schema, err := convertSchema(t.Parameters)
			if err != nil {
				return nil, fmt.Errorf("failed to convert tool schema: %w", err)
			}
			funcDecls = append(funcDecls, &genai.FunctionDeclaration{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  schema,
			})
		}
		tools = append(tools, &genai.Tool{FunctionDeclarations: funcDecls})
	}

	model := p.model
	if params.Model != "" {
		model = params.Model
	}

	temp := float32(params.Temperature)
	maxTokens := int32(params.MaxTokens)

	conf := &genai.GenerateContentConfig{
		Temperature:     &temp,
		MaxOutputTokens: maxTokens, // Updated: int32 value
		Tools:           tools,
	}

	// ToolConfig (auto)
	// if len(params.Tools) > 0 {
	// 	conf.ToolConfig = &genai.ToolConfig{FunctionCallingConfig: &genai.FunctionCallingConfig{Mode: "AUTO"}}
	// }

	// Streaming
	streamIter := p.client.Models.GenerateContentStream(ctx, model, contents, conf)

	return func(yield func(provider.StreamEvent, error) bool) {
		for resp, err := range streamIter {
			if err != nil {
				yield(provider.StreamEvent{Type: provider.StreamEventError, Error: err}, err)
				return
			}

			// Handle response parts
			for _, cand := range resp.Candidates {
				if cand.Content != nil {
					for _, part := range cand.Content.Parts {
						if part.Text != "" {
							if !yield(provider.StreamEvent{
								Type: provider.StreamEventPlainText,
								Text: part.Text,
							}, nil) {
								return
							}
						}
						if part.FunctionCall != nil {
							argsJSON, _ := json.Marshal(part.FunctionCall.Args)
							if !yield(provider.StreamEvent{
								Type: provider.StreamEventToolCall,
								ToolCall: &provider.ToolCall{
									ID:        part.FunctionCall.Name, // Use name as ID if ID missing
									Name:      part.FunctionCall.Name,
									Arguments: string(argsJSON),
								},
							}, nil) {
								return
							}
						}
					}
				}
			}
		}
		yield(provider.StreamEvent{Type: provider.StreamEventDone}, nil)
	}, nil
}

func (p *GeminiProvider) ListModels(ctx context.Context) ([]provider.ModelInfo, error) {
	// Basic implementation using configured client
	// p.client.Models.List(ctx, nil) returns iterator

	// Example hardcoded for now as API exploration might take time
	return []provider.ModelInfo{
		{ID: "gemini-2.0-flash-exp", Name: "Gemini 2.0 Flash Exp"},
		{ID: "gemini-1.5-pro", Name: "Gemini 1.5 Pro"},
		{ID: "gemini-1.5-flash", Name: "Gemini 1.5 Flash"},
	}, nil
}

func convertSchema(schemaMap map[string]interface{}) (*genai.Schema, error) {
	b, err := json.Marshal(schemaMap)
	if err != nil {
		return nil, err
	}
	var s genai.Schema
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return &s, nil
}
