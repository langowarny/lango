package openai

import (
	"testing"
)

func TestNewProvider(t *testing.T) {
	p := NewProvider("openai", "test-key", "http://localhost:1234")
	if p.ID() != "openai" {
		t.Errorf("expected ID 'openai', got %s", p.ID())
	}
}
