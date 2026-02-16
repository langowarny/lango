package orchestration

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	adk_tool "google.golang.org/adk/tool"

	"github.com/langowarny/lango/internal/agent"
)

func newTestTool(name string) *agent.Tool {
	return &agent.Tool{
		Name:        name,
		Description: "test tool " + name,
		SafetyLevel: agent.SafetyLevelSafe,
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return nil, nil
		},
	}
}

// stubTool implements adk_tool.Tool for testing.
type stubTool struct {
	name string
}

func (s *stubTool) Name() string       { return s.name }
func (s *stubTool) Description() string { return "stub " + s.name }
func (s *stubTool) IsLongRunning() bool { return false }

// stubAdapter is a ToolAdapter that returns a stubTool without real ADK wiring.
func stubAdapter(t *agent.Tool) (adk_tool.Tool, error) {
	return &stubTool{name: t.Name}, nil
}

// failingAdapter always returns an error.
func failingAdapter(t *agent.Tool) (adk_tool.Tool, error) {
	return nil, fmt.Errorf("adapter error for %s", t.Name)
}

func TestPartitionTools(t *testing.T) {
	tests := []struct {
		name         string
		give         []*agent.Tool
		wantExecutor []string
		wantResearch []string
		wantMemory   []string
		wantPlanner  []string
	}{
		{
			name: "executor prefixes",
			give: []*agent.Tool{
				newTestTool("exec_shell"),
				newTestTool("fs_read"),
				newTestTool("browser_navigate"),
				newTestTool("crypto_sign"),
				newTestTool("skill_deploy"),
			},
			wantExecutor: []string{"exec_shell", "fs_read", "browser_navigate", "crypto_sign", "skill_deploy"},
		},
		{
			name: "researcher prefixes",
			give: []*agent.Tool{
				newTestTool("search_web"),
				newTestTool("rag_query"),
				newTestTool("graph_traverse"),
			},
			wantResearch: []string{"search_web", "rag_query", "graph_traverse"},
		},
		{
			name: "memory prefixes",
			give: []*agent.Tool{
				newTestTool("memory_store"),
				newTestTool("observe_event"),
				newTestTool("reflect_summary"),
			},
			wantMemory: []string{"memory_store", "observe_event", "reflect_summary"},
		},
		{
			name: "unmatched tools go to executor",
			give: []*agent.Tool{
				newTestTool("custom_action"),
				newTestTool("do_something"),
			},
			wantExecutor: []string{"custom_action", "do_something"},
		},
		{
			name: "mixed tools partitioned correctly",
			give: []*agent.Tool{
				newTestTool("exec_run"),
				newTestTool("search_docs"),
				newTestTool("memory_save"),
				newTestTool("unknown_tool"),
			},
			wantExecutor: []string{"exec_run", "unknown_tool"},
			wantResearch: []string{"search_docs"},
			wantMemory:   []string{"memory_save"},
		},
		{
			name: "empty input",
			give: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PartitionTools(tt.give)

			assert.Equal(t, tt.wantExecutor, toolNames(got.Executor), "executor tools")
			assert.Equal(t, tt.wantResearch, toolNames(got.Researcher), "researcher tools")
			assert.Equal(t, tt.wantMemory, toolNames(got.MemoryManager), "memory tools")
			assert.Equal(t, tt.wantPlanner, toolNames(got.Planner), "planner tools")
		})
	}
}

func TestBuildAgentTree_NilAdaptTool(t *testing.T) {
	_, err := BuildAgentTree(Config{
		Tools:        []*agent.Tool{newTestTool("exec_shell")},
		Model:        nil,
		SystemPrompt: "test",
		AdaptTool:    nil,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "AdaptTool is required")
}

func TestBuildAgentTree_Success(t *testing.T) {
	tools := []*agent.Tool{
		newTestTool("exec_shell"),
		newTestTool("search_web"),
		newTestTool("memory_store"),
		newTestTool("custom_tool"),
	}

	root, err := BuildAgentTree(Config{
		Tools:        tools,
		Model:        nil, // nil is accepted at construction time
		SystemPrompt: "test prompt",
		AdaptTool:    stubAdapter,
	})
	require.NoError(t, err)
	require.NotNil(t, root)

	assert.Equal(t, "lango-orchestrator", root.Name())
	assert.Len(t, root.SubAgents(), 4, "orchestrator should have 4 sub-agents")

	subNames := make([]string, len(root.SubAgents()))
	for i, sa := range root.SubAgents() {
		subNames[i] = sa.Name()
	}
	assert.Contains(t, subNames, "executor")
	assert.Contains(t, subNames, "researcher")
	assert.Contains(t, subNames, "planner")
	assert.Contains(t, subNames, "memory-manager")
}

func TestBuildAgentTree_AdapterError(t *testing.T) {
	tools := []*agent.Tool{newTestTool("exec_shell")}

	_, err := BuildAgentTree(Config{
		Tools:        tools,
		Model:        nil,
		SystemPrompt: "test",
		AdaptTool:    failingAdapter,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "adapt executor tools")
}

func TestPartitionTools_PrefixPriority(t *testing.T) {
	// Verify researcher prefixes are checked before executor defaults.
	tools := []*agent.Tool{
		newTestTool("search_rag"),
		newTestTool("graph_node"),
	}

	got := PartitionTools(tools)

	assert.Empty(t, got.Executor, "no tools should go to executor")
	assert.Len(t, got.Researcher, 2, "both should be researcher")
}

// toolNames extracts names from a tool slice for assertions.
func toolNames(tools []*agent.Tool) []string {
	if len(tools) == 0 {
		return nil
	}
	names := make([]string, len(tools))
	for i, t := range tools {
		names[i] = t.Name
	}
	return names
}
