package adk

import (
	"context"
	"fmt"
	"iter"
	"regexp"
	"strings"

	"go.uber.org/zap"
	adk_agent "google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/langowarny/lango/internal/logging"
	internal "github.com/langowarny/lango/internal/session"
)

func logger() *zap.SugaredLogger { return logging.Agent() }

// Agent wraps the ADK runner for integration with Lango.
type Agent struct {
	runner   *runner.Runner
	adkAgent adk_agent.Agent
}

// NewAgent creates a new Agent instance.
func NewAgent(ctx context.Context, tools []tool.Tool, mod model.LLM, systemPrompt string, store internal.Store) (*Agent, error) {
	// Create LLM Agent
	cfg := llmagent.Config{
		Name:        "lango-agent",
		Description: "Lango Assistant",
		Model:       mod,
		Tools:       tools,
		Instruction: systemPrompt,
	}

	adkAgent, err := llmagent.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create llm agent: %w", err)
	}

	// Create Session Service
	sessService := NewSessionServiceAdapter(store, "lango-agent")

	// Create Runner
	runnerCfg := runner.Config{
		AppName:        "lango",
		Agent:          adkAgent,
		SessionService: sessService,
	}

	r, err := runner.New(runnerCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create runner: %w", err)
	}

	return &Agent{
		runner:   r,
		adkAgent: adkAgent,
	}, nil
}

// NewAgentFromADK creates a Lango Agent wrapping a pre-built ADK agent.
// Used for multi-agent orchestration where the agent tree is built externally.
func NewAgentFromADK(adkAgent adk_agent.Agent, store internal.Store) (*Agent, error) {
	sessService := NewSessionServiceAdapter(store, adkAgent.Name())

	runnerCfg := runner.Config{
		AppName:        "lango",
		Agent:          adkAgent,
		SessionService: sessService,
	}

	r, err := runner.New(runnerCfg)
	if err != nil {
		return nil, fmt.Errorf("create runner: %w", err)
	}

	return &Agent{runner: r, adkAgent: adkAgent}, nil
}

// ADKAgent returns the underlying ADK agent, or nil if not available.
func (a *Agent) ADKAgent() adk_agent.Agent {
	return a.adkAgent
}

// Run executes the agent for a given session and returns an event iterator.
func (a *Agent) Run(ctx context.Context, sessionID string, input string) iter.Seq2[*session.Event, error] {
	// Create user content
	userMsg := &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: input}},
	}

	// Config for run
	runCfg := adk_agent.RunConfig{
		// Defaults
	}

	// Execute via Runner
	return a.runner.Run(ctx, "user", sessionID, userMsg, runCfg)
}

// RunAndCollect executes the agent and returns the full text response.
// If the agent encounters a "failed to find agent" error (hallucinated agent
// name), it sends a correction message and retries once.
func (a *Agent) RunAndCollect(ctx context.Context, sessionID, input string) (string, error) {
	resp, err := a.runAndCollectOnce(ctx, sessionID, input)
	if err == nil {
		return resp, nil
	}

	badAgent := extractMissingAgent(err)
	if badAgent == "" || len(a.adkAgent.SubAgents()) == 0 {
		return "", err
	}

	// Build correction message and retry once.
	names := subAgentNames(a.adkAgent)
	correction := fmt.Sprintf(
		"[System: Agent %q does not exist. Valid agents: %s. Please retry using one of the valid agent names listed above.]",
		badAgent, strings.Join(names, ", "))
	logger().Warnw("agent name hallucination detected, retrying",
		"hallucinated", badAgent,
		"valid_agents", names,
		"session", sessionID)

	return a.runAndCollectOnce(ctx, sessionID, correction)
}

// runAndCollectOnce executes a single agent run and collects text output.
// It tracks whether partial (streaming) events were seen to avoid
// double-counting text that appears in both partial chunks and the
// final non-partial response.
func (a *Agent) runAndCollectOnce(ctx context.Context, sessionID, input string) (string, error) {
	var b strings.Builder
	var sawPartial bool

	for event, err := range a.Run(ctx, sessionID, input) {
		if err != nil {
			return "", fmt.Errorf("agent error: %w", err)
		}

		// Log agent event for multi-agent observability.
		if event.Author != "" {
			if event.Actions.TransferToAgent != "" {
				logger().Debugw("agent delegation",
					"from", event.Author,
					"to", event.Actions.TransferToAgent,
					"session", sessionID)
			} else if hasText(event) {
				logger().Debugw("agent response",
					"agent", event.Author,
					"session", sessionID)
			}
		}

		if event.Content == nil {
			continue
		}

		if event.Partial {
			// Streaming text chunk — collect incrementally.
			sawPartial = true
			for _, part := range event.Content.Parts {
				if part.Text != "" {
					b.WriteString(part.Text)
				}
			}
		} else if !sawPartial {
			// Non-streaming mode: no partial events were seen,
			// so collect from the final complete response.
			for _, part := range event.Content.Parts {
				if part.Text != "" {
					b.WriteString(part.Text)
				}
			}
		}
		// If sawPartial && !event.Partial: this is the final done event
		// in streaming mode. Its text duplicates partial chunks, so skip.
	}

	return b.String(), nil
}

// reAgentNotFound matches ADK's "failed to find agent: <name>" error.
var reAgentNotFound = regexp.MustCompile(`failed to find agent: (\S+)`)

// extractMissingAgent returns the hallucinated agent name from an error,
// or an empty string if the error does not match the pattern.
func extractMissingAgent(err error) string {
	m := reAgentNotFound.FindStringSubmatch(err.Error())
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

// subAgentNames returns the names of all immediate sub-agents.
func subAgentNames(a adk_agent.Agent) []string {
	subs := a.SubAgents()
	names := make([]string, len(subs))
	for i, s := range subs {
		names[i] = s.Name()
	}
	return names
}

// hasText reports whether the event contains any non-empty text part.
func hasText(e *session.Event) bool {
	if e.Content == nil {
		return false
	}
	for _, p := range e.Content.Parts {
		if p.Text != "" {
			return true
		}
	}
	return false
}
