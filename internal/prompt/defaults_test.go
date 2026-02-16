package prompt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBuilder_ContainsAllSections(t *testing.T) {
	b := DefaultBuilder()
	assert.True(t, b.Has(SectionIdentity))
	assert.True(t, b.Has(SectionSafety))
	assert.True(t, b.Has(SectionConversationRules))
	assert.True(t, b.Has(SectionToolUsage))
}

func TestDefaultBuilder_IncludesConversationRules(t *testing.T) {
	result := DefaultBuilder().Build()
	assert.Contains(t, result, "Focus exclusively on the current question")
	assert.Contains(t, result, "Do not repeat or summarize your previous answers")
}

func TestDefaultBuilder_IncludesIdentity(t *testing.T) {
	result := DefaultBuilder().Build()
	assert.Contains(t, result, "You are Lango")
}

func TestDefaultBuilder_SectionOrder(t *testing.T) {
	result := DefaultBuilder().Build()
	idxIdentity := strings.Index(result, "You are Lango")
	idxSafety := strings.Index(result, "Safety Guidelines")
	idxConversation := strings.Index(result, "Conversation Rules")
	idxTool := strings.Index(result, "Tool Usage Guidelines")

	assert.Less(t, idxIdentity, idxSafety, "Identity should come before Safety")
	assert.Less(t, idxSafety, idxConversation, "Safety should come before Conversation Rules")
	assert.Less(t, idxConversation, idxTool, "Conversation Rules should come before Tool Usage")
}
