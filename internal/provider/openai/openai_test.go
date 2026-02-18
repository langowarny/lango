package openai

import (
	"context"
	"testing"
)

func TestNewProvider(t *testing.T) {
	p := NewProvider("openai", "test-key", "http://localhost:1234")
	if p.ID() != "openai" {
		t.Errorf("expected ID 'openai', got %s", p.ID())
	}
}

func TestOpenAIProvider_ListModels(t *testing.T) {
	// ListModels calls the real API, so we test that it returns an error
	// when the server is not available (invalid base URL)
	p := NewProvider("openai", "test-key", "http://localhost:1/v1")
	_, err := p.ListModels(context.Background())
	if err == nil {
		t.Error("expected error when connecting to unavailable server")
	}
}
