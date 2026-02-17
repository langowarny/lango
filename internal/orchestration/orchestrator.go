package orchestration

import (
	"fmt"
	"strings"

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
	// MaxDelegationRounds limits the number of orchestrator→sub-agent
	// delegation rounds per user turn. Zero means use default (5).
	MaxDelegationRounds int
}

// BuildAgentTree creates a hierarchical agent tree with an orchestrator root
// and specialized sub-agents. Sub-agents are only created when they have
// tools assigned (except Planner which is LLM-only and always included).
func BuildAgentTree(cfg Config) (adk_agent.Agent, error) {
	if cfg.AdaptTool == nil {
		return nil, fmt.Errorf("build agent tree: AdaptTool is required")
	}

	rs := PartitionTools(cfg.Tools)

	var subAgents []adk_agent.Agent
	var agentDescriptions []string

	// Executor: only if tools are assigned.
	if len(rs.Executor) > 0 {
		executorTools, err := adaptTools(cfg.AdaptTool, rs.Executor)
		if err != nil {
			return nil, fmt.Errorf("adapt executor tools: %w", err)
		}
		caps := capabilityDescription(rs.Executor)
		a, err := llmagent.New(llmagent.Config{
			Name:        "executor",
			Description: fmt.Sprintf("Executes actions. Capabilities: %s", caps),
			Model:       cfg.Model,
			Tools:       executorTools,
			Instruction: "You are the Executor agent. Execute tool calls precisely as requested. Report results accurately. After completing the requested action, provide the results clearly. If a tool fails, report the error without retrying unless explicitly asked.",
		})
		if err != nil {
			return nil, fmt.Errorf("create executor agent: %w", err)
		}
		subAgents = append(subAgents, a)
		agentDescriptions = append(agentDescriptions,
			fmt.Sprintf("- \"executor\": %s", caps))
	}

	// Researcher: only if tools are assigned.
	if len(rs.Researcher) > 0 {
		researcherTools, err := adaptTools(cfg.AdaptTool, rs.Researcher)
		if err != nil {
			return nil, fmt.Errorf("adapt researcher tools: %w", err)
		}
		caps := capabilityDescription(rs.Researcher)
		a, err := llmagent.New(llmagent.Config{
			Name:        "researcher",
			Description: fmt.Sprintf("Searches knowledge bases and retrieves relevant information. Capabilities: %s", caps),
			Model:       cfg.Model,
			Tools:       researcherTools,
			Instruction: "You are the Researcher agent. Search and retrieve relevant information from knowledge bases, semantic search, and the knowledge graph. Provide comprehensive, well-organized results. After completing research, summarize your findings clearly.",
		})
		if err != nil {
			return nil, fmt.Errorf("create researcher agent: %w", err)
		}
		subAgents = append(subAgents, a)
		agentDescriptions = append(agentDescriptions,
			fmt.Sprintf("- \"researcher\": %s", caps))
	}

	// Planner: always included (LLM-only reasoning agent).
	{
		a, err := llmagent.New(llmagent.Config{
			Name:        "planner",
			Description: "Decomposes complex tasks into steps and designs execution plans. Has no tools — uses LLM reasoning only.",
			Model:       cfg.Model,
			Instruction: "You are the Planner agent. Break complex tasks into clear, actionable steps. Consider dependencies between steps. Output structured plans. After planning, present the plan for review.",
		})
		if err != nil {
			return nil, fmt.Errorf("create planner agent: %w", err)
		}
		subAgents = append(subAgents, a)
		agentDescriptions = append(agentDescriptions,
			"- \"planner\": task decomposition and planning (no tools, LLM reasoning only)")
	}

	// Memory Manager: only if memory tools are assigned.
	if len(rs.MemoryManager) > 0 {
		memoryTools, err := adaptTools(cfg.AdaptTool, rs.MemoryManager)
		if err != nil {
			return nil, fmt.Errorf("adapt memory tools: %w", err)
		}
		caps := capabilityDescription(rs.MemoryManager)
		a, err := llmagent.New(llmagent.Config{
			Name:        "memory-manager",
			Description: fmt.Sprintf("Manages conversational memory including observations and reflections. Capabilities: %s", caps),
			Model:       cfg.Model,
			Tools:       memoryTools,
			Instruction: "You are the Memory Manager agent. Manage observations, reflections, and the memory graph. Organize and retrieve relevant past interactions. After completing memory operations, report what was stored or retrieved.",
		})
		if err != nil {
			return nil, fmt.Errorf("create memory agent: %w", err)
		}
		subAgents = append(subAgents, a)
		agentDescriptions = append(agentDescriptions,
			fmt.Sprintf("- \"memory-manager\": %s", caps))
	}

	// Append remote A2A agents if configured.
	subAgents = append(subAgents, cfg.RemoteAgents...)
	for _, ra := range cfg.RemoteAgents {
		agentDescriptions = append(agentDescriptions,
			fmt.Sprintf("- %s: %s (remote A2A agent)", ra.Name(), ra.Description()))
	}

	// Build orchestrator instruction focused on delegation.
	// The orchestrator has NO tools — it must delegate all tool-requiring
	// tasks to the appropriate sub-agent.
	maxRounds := cfg.MaxDelegationRounds
	if maxRounds <= 0 {
		maxRounds = 5
	}

	orchestratorInstruction := cfg.SystemPrompt + fmt.Sprintf(`

You are the orchestrator. You coordinate specialized sub-agents to fulfill user requests.

## Your Role
You do NOT have tools. You MUST delegate all tool-requiring tasks to the appropriate sub-agent using transfer_to_agent.

## Available Sub-Agents (use EXACTLY these names)
%s

## Delegation Rules
1. For any action that requires tools: delegate to the sub-agent listed above whose description best matches.
2. For simple conversational messages (greetings, opinions, general knowledge): respond directly without delegation.
3. Maximum %d delegation rounds per user turn.

## CRITICAL
- You MUST use the EXACT agent name from the list above (e.g. "executor", NOT "exec", "browser", or any abbreviation).
- NEVER invent or abbreviate agent names. If unsure, pick the closest match from the list above.`, strings.Join(agentDescriptions, "\n"), maxRounds)

	orchestrator, err := llmagent.New(llmagent.Config{
		Name:        "lango-orchestrator",
		Description: "Lango Assistant Orchestrator",
		Model:       cfg.Model,
		Tools:       nil,
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
