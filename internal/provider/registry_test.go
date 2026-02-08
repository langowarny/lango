package provider

import (
	"context"
	"iter"
	"testing"
)

type mockProvider struct {
	id string
}

func (m *mockProvider) ID() string {
	return m.id
}

func (m *mockProvider) Generate(ctx context.Context, params GenerateParams) (iter.Seq2[StreamEvent, error], error) {
	return nil, nil
}

func (m *mockProvider) ListModels(ctx context.Context) ([]ModelInfo, error) {
	return nil, nil
}

func TestRegistry(t *testing.T) {
	r := NewRegistry()

	openai := &mockProvider{id: "openai"}
	anthropic := &mockProvider{id: "anthropic"}

	r.Register(openai)
	r.Register(anthropic)

	tests := []struct {
		name   string
		query  string
		wantID string
		wantOK bool
	}{
		{"exact match openai", "openai", "openai", true},
		{"exact match anthropic", "anthropic", "anthropic", true},
		{"alias gpt", "gpt", "openai", true},
		{"alias claude", "claude", "anthropic", true},
		{"case insensitive", "OpenAI", "openai", true},
		{"unknown", "unknown", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := r.Get(tt.query)
			if ok != tt.wantOK {
				t.Errorf("Get(%q) ok = %v, want %v", tt.query, ok, tt.wantOK)
			}
			if ok && got.ID() != tt.wantID {
				t.Errorf("Get(%q) ID = %v, want %v", tt.query, got.ID(), tt.wantID)
			}
		})
	}
}
