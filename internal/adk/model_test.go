package adk

import (
	"context"
	"iter"
	"testing"

	"github.com/langowarny/lango/internal/provider"
	"google.golang.org/adk/model"
)

type mockProvider struct {
	id     string
	events []provider.StreamEvent
	err    error
}

func (m *mockProvider) ID() string { return m.id }

func (m *mockProvider) Generate(_ context.Context, _ provider.GenerateParams) (iter.Seq2[provider.StreamEvent, error], error) {
	if m.err != nil {
		return nil, m.err
	}
	return func(yield func(provider.StreamEvent, error) bool) {
		for _, evt := range m.events {
			if !yield(evt, nil) {
				return
			}
		}
	}, nil
}

func (m *mockProvider) ListModels(_ context.Context) ([]provider.ModelInfo, error) {
	return nil, nil
}

func TestModelAdapter_Name(t *testing.T) {
	p := &mockProvider{id: "test-provider"}
	adapter := NewModelAdapter(p)

	if adapter.Name() != "test-provider" {
		t.Errorf("expected 'test-provider', got %q", adapter.Name())
	}
}

func TestModelAdapter_GenerateContent_TextDelta(t *testing.T) {
	p := &mockProvider{
		id: "test",
		events: []provider.StreamEvent{
			{Type: provider.StreamEventPlainText, Text: "Hello "},
			{Type: provider.StreamEventPlainText, Text: "world"},
			{Type: provider.StreamEventDone},
		},
	}
	adapter := NewModelAdapter(p)

	req := &model.LLMRequest{Model: "test-model"}
	seq := adapter.GenerateContent(context.Background(), req, true)

	var responses []*model.LLMResponse
	for resp, err := range seq {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		responses = append(responses, resp)
	}

	if len(responses) != 3 {
		t.Fatalf("expected 3 responses, got %d", len(responses))
	}

	// First two should be partial text
	if !responses[0].Partial {
		t.Error("expected first response to be partial")
	}
	if responses[0].Content.Parts[0].Text != "Hello " {
		t.Errorf("expected 'Hello ', got %q", responses[0].Content.Parts[0].Text)
	}

	// Last should be turn complete
	if !responses[2].TurnComplete {
		t.Error("expected last response to be turn complete")
	}
	if responses[2].Partial {
		t.Error("expected last response to not be partial")
	}
}

func TestModelAdapter_GenerateContent_ProviderError(t *testing.T) {
	p := &mockProvider{
		id:  "test",
		err: context.DeadlineExceeded,
	}
	adapter := NewModelAdapter(p)

	req := &model.LLMRequest{Model: "test-model"}
	seq := adapter.GenerateContent(context.Background(), req, false)

	for _, err := range seq {
		if err == nil {
			t.Fatal("expected error from provider")
		}
		return // Only check first yield
	}
	t.Fatal("expected at least one yield")
}

func TestModelAdapter_GenerateContent_ToolCall(t *testing.T) {
	p := &mockProvider{
		id: "test",
		events: []provider.StreamEvent{
			{
				Type: provider.StreamEventToolCall,
				ToolCall: &provider.ToolCall{
					ID:        "call_1",
					Name:      "exec",
					Arguments: `{"command":"ls"}`,
				},
			},
			{Type: provider.StreamEventDone},
		},
	}
	adapter := NewModelAdapter(p)

	req := &model.LLMRequest{Model: "test-model"}
	seq := adapter.GenerateContent(context.Background(), req, false)

	var responses []*model.LLMResponse
	for resp, err := range seq {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		responses = append(responses, resp)
	}

	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}

	// First should have function call
	parts := responses[0].Content.Parts
	hasFuncCall := false
	for _, p := range parts {
		if p.FunctionCall != nil {
			hasFuncCall = true
			if p.FunctionCall.Name != "exec" {
				t.Errorf("expected function name 'exec', got %q", p.FunctionCall.Name)
			}
			if p.FunctionCall.Args["command"] != "ls" {
				t.Errorf("expected arg command='ls', got %v", p.FunctionCall.Args["command"])
			}
		}
	}
	if !hasFuncCall {
		t.Error("expected a FunctionCall part")
	}
}

func TestModelAdapter_GenerateContent_StreamError(t *testing.T) {
	p := &mockProvider{
		id: "test",
		events: []provider.StreamEvent{
			{Type: provider.StreamEventPlainText, Text: "partial"},
			{Type: provider.StreamEventError, Error: context.Canceled},
		},
	}
	adapter := NewModelAdapter(p)

	req := &model.LLMRequest{Model: "test-model"}
	seq := adapter.GenerateContent(context.Background(), req, true)

	gotError := false
	for _, err := range seq {
		if err != nil {
			gotError = true
			break
		}
	}
	if !gotError {
		t.Error("expected error event to propagate")
	}
}
