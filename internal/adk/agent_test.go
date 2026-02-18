package adk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractMissingAgent(t *testing.T) {
	tests := []struct {
		name string
		give error
		want string
	}{
		{
			name: "standard ADK error",
			give: fmt.Errorf("agent error: failed to find agent: browser_agent"),
			want: "browser_agent",
		},
		{
			name: "wrapped error",
			give: fmt.Errorf("outer: %w", fmt.Errorf("failed to find agent: exec")),
			want: "exec",
		},
		{
			name: "unrelated error",
			give: fmt.Errorf("connection refused"),
			want: "",
		},
		{
			name: "partial match no agent name",
			give: fmt.Errorf("failed to find agent: "),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractMissingAgent(tt.give)
			assert.Equal(t, tt.want, got)
		})
	}
}
