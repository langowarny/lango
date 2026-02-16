package orchestration

import (
	"fmt"

	adk_agent "google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	adk_tool "google.golang.org/adk/tool"

	"github.com/langowarny/lango/internal/agent"
)

// ToolAdapter converts an internal agent.Tool to an ADK tool.Tool.
// This is injected to avoid a direct dependency on the adk package,
// which carries transitive imports that may cause import cycles.
type ToolAdapter func(t *agent.Tool) (adk_tool.Tool, error)

// Config holds orchestration configuration.
type Config struct {
	// Tools is the full set of available tools.
	Tools []*agent.Tool
	// Model is the primary LLM model adapter.
	Model model.LLM
	// SystemPrompt is the base system instruction.
	SystemPrompt string
	// AdaptTool converts an internal tool to an ADK tool.
	// Callers should pass adk.AdaptTool.
	AdaptTool ToolAdapter
	// RemoteAgents are external A2A agents to include as sub-agents.
	RemoteAgents []adk_agent.Agent
}

// BuildAgentTree creates a hierarchical agent tree with an orchestrator root
// and specialized sub-agents (executor, researcher, planner, memory-manager).
func BuildAgentTree(cfg Config) (adk_agent.Agent, error) {
	if cfg.AdaptTool == nil {
		return nil, fmt.Errorf("build agent tree: AdaptTool is required")
	}

	rs := PartitionTools(cfg.Tools)

	executorTools, err := adaptTools(cfg.AdaptTool, rs.Executor)
	if err != nil {
		return nil, fmt.Errorf("adapt executor tools: %w", err)
	}

	researcherTools, err := adaptTools(cfg.AdaptTool, rs.Researcher)
	if err != nil {
		return nil, fmt.Errorf("adapt researcher tools: %w", err)
	}

	memoryTools, err := adaptTools(cfg.AdaptTool, rs.MemoryManager)
	if err != nil {
		return nil, fmt.Errorf("adapt memory tools: %w", err)
	}

	executorAgent, err := llmagent.New(llmagent.Config{
		Name:        "executor",
		Description: "Executes tools including shell commands, file operations, browser automation, and cryptographic operations. Delegate to this agent when the user needs to perform actions.",
		Model:       cfg.Model,
		Tools:       executorTools,
		Instruction: "You are the Executor agent. Execute tool calls precisely as requested. Report results accurately. If a tool fails, report the error without retrying unless explicitly asked.",
	})
	if err != nil {
		return nil, fmt.Errorf("create executor agent: %w", err)
	}

	researcherAgent, err := llmagent.New(llmagent.Config{
		Name:        "researcher",
		Description: "Searches knowledge bases, performs RAG retrieval, and traverses the knowledge graph. Delegate to this agent for information lookup and research tasks.",
		Model:       cfg.Model,
		Tools:       researcherTools,
		Instruction: "You are the Researcher agent. Search and retrieve relevant information from knowledge bases, semantic search, and the knowledge graph. Provide comprehensive, well-organized results.",
	})
	if err != nil {
		return nil, fmt.Errorf("create researcher agent: %w", err)
	}

	plannerAgent, err := llmagent.New(llmagent.Config{
		Name:        "planner",
		Description: "Decomposes complex tasks into steps and designs execution plans. Delegate to this agent when the user needs task planning or strategy.",
		Model:       cfg.Model,
		Instruction: "You are the Planner agent. Break complex tasks into clear, actionable steps. Consider dependencies between steps. Output structured plans.",
	})
	if err != nil {
		return nil, fmt.Errorf("create planner agent: %w", err)
	}

	memoryAgent, err := llmagent.New(llmagent.Config{
		Name:        "memory-manager",
		Description: "Manages conversational memory including observations, reflections, and the memory graph. Delegate to this agent for memory-related operations.",
		Model:       cfg.Model,
		Tools:       memoryTools,
		Instruction: "You are the Memory Manager agent. Manage observations, reflections, and the memory graph. Organize and retrieve relevant past interactions.",
	})
	if err != nil {
		return nil, fmt.Errorf("create memory agent: %w", err)
	}

	subAgents := []adk_agent.Agent{executorAgent, researcherAgent, plannerAgent, memoryAgent}

	// Append remote A2A agents if configured.
	subAgents = append(subAgents, cfg.RemoteAgents...)

	orchestratorInstruction := cfg.SystemPrompt + "\n\nYou are the orchestrator. Route tasks to the most appropriate sub-agent based on the task type. You have these sub-agents:\n- executor: for running tools and performing actions\n- researcher: for searching knowledge and information retrieval\n- planner: for task decomposition and planning\n- memory-manager: for memory operations"

	// Mention remote agents in instruction.
	for _, ra := range cfg.RemoteAgents {
		orchestratorInstruction += fmt.Sprintf("\n- %s: %s (remote A2A agent)", ra.Name(), ra.Description())
	}
	orchestratorInstruction += "\n\nDelegate tasks appropriately and synthesize results for the user."

	orchestrator, err := llmagent.New(llmagent.Config{
		Name:        "lango-orchestrator",
		Description: "Lango Assistant Orchestrator",
		Model:       cfg.Model,
		SubAgents:   subAgents,
		Instruction: orchestratorInstruction,
	})
	if err != nil {
		return nil, fmt.Errorf("create orchestrator agent: %w", err)
	}

	return orchestrator, nil
}

// adaptTools converts a slice of internal agent tools to ADK tools using the provided adapter.
func adaptTools(adapt ToolAdapter, tools []*agent.Tool) ([]adk_tool.Tool, error) {
	result := make([]adk_tool.Tool, 0, len(tools))
	for _, t := range tools {
		adapted, err := adapt(t)
		if err != nil {
			return nil, fmt.Errorf("adapt tool %q: %w", t.Name, err)
		}
		result = append(result, adapted)
	}
	return result, nil
}
