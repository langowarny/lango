package prompts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFS_ContainsAllPromptFiles(t *testing.T) {
	files := []string{
		"AGENTS.md",
		"SAFETY.md",
		"CONVERSATION_RULES.md",
		"TOOL_USAGE.md",
	}

	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			data, err := FS.ReadFile(name)
			require.NoError(t, err, "embedded file %s must exist", name)
			assert.NotEmpty(t, data, "embedded file %s must not be empty", name)
		})
	}
}
