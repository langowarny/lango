package anthropic

import (
	"testing"
)

func TestNewProvider(t *testing.T) {
	p := NewProvider("test-key")
	if p.ID() != "anthropic" {
		t.Errorf("expected ID 'anthropic', got %s", p.ID())
	}
}
