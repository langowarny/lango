package orchestration

import (
	"strings"

	"github.com/langowarny/lango/internal/agent"
)

// RoleToolSet defines which tools belong to each sub-agent role.
type RoleToolSet struct {
	Executor      []*agent.Tool
	Researcher    []*agent.Tool
	Planner       []*agent.Tool
	MemoryManager []*agent.Tool
}

// executorPrefixes are tool name prefixes assigned to the Executor sub-agent.
var executorPrefixes = []string{"exec", "fs_", "browser_", "crypto_", "skill_"}

// researcherPrefixes are tool name prefixes assigned to the Researcher sub-agent.
var researcherPrefixes = []string{"search_", "rag_", "graph_", "save_knowledge", "save_learning"}

// memoryPrefixes are tool name prefixes assigned to the Memory Manager sub-agent.
var memoryPrefixes = []string{"memory_", "observe_", "reflect_"}

// PartitionTools splits tools into role-specific sets based on tool name prefixes.
// Tools not matching any prefix are assigned to the Executor (default).
func PartitionTools(tools []*agent.Tool) RoleToolSet {
	var rs RoleToolSet
	for _, t := range tools {
		switch {
		case matchesPrefix(t.Name, researcherPrefixes):
			rs.Researcher = append(rs.Researcher, t)
		case matchesPrefix(t.Name, memoryPrefixes):
			rs.MemoryManager = append(rs.MemoryManager, t)
		case matchesPrefix(t.Name, executorPrefixes):
			rs.Executor = append(rs.Executor, t)
		default:
			// Unmatched tools default to Executor.
			rs.Executor = append(rs.Executor, t)
		}
	}
	// Planner has no tools (LLM-only reasoning).
	return rs
}

// matchesPrefix returns true if name starts with any of the given prefixes.
func matchesPrefix(name string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}
