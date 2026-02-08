package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/langowarny/lango/internal/logging"
	"github.com/langowarny/lango/internal/provider"
	"github.com/langowarny/lango/internal/session"
)

var logger = logging.SubsystemSugar("agent")

// Runtime represents the agent runtime
type Runtime struct {
	config       Config
	provider     provider.Provider // Replaces registry + string ID lookup
	tools        map[string]*Tool
	adkTools     []*AdkToolAdapter // Keep for compatibility if needed
	toolsMu      sync.RWMutex
	sessionStore session.Store
}

// Config holds agent runtime configuration
type Config struct {
	Provider string // Provider ID
	Model    string // Model ID
	// APIKey removed - handled by Supervisor
	MaxTokens            int     // context window limit
	Temperature          float64 // generation temperature
	MaxConversationTurns int     // max conversation turns
	FallbackProvider     string  // fallback provider ID
	FallbackModel        string  // fallback model ID
	SystemPrompt         string  // system prompt for the agent
}

// Tool represents a tool that can be invoked by the LLM
type Tool struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
	Handler     ToolHandler
}

// ParameterDef defines a tool parameter
type ParameterDef struct {
	Type        string
	Description string
	Required    bool
	Enum        []string
}

// ToolHandler is the function signature for tool implementations
type ToolHandler func(ctx context.Context, params map[string]interface{}) (interface{}, error)

// StreamEvent represents a streaming event
type StreamEvent struct {
	Type     string             // text_delta, tool_start, tool_end, error, done
	Text     string             // for text_delta
	ToolCall *provider.ToolCall // for tool events
	Error    error              // for error events
}

// New creates a new agent runtime
func New(cfg Config, store session.Store, p provider.Provider) (*Runtime, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return &Runtime{
		config:       cfg,
		provider:     p,
		tools:        make(map[string]*Tool),
		adkTools:     make([]*AdkToolAdapter, 0),
		sessionStore: store,
	}, nil
}

// validateConfig checks if the configuration is valid
func validateConfig(cfg Config) error {
	if cfg.Provider == "" {
		return fmt.Errorf("provider is required")
	}
	if cfg.Model == "" {
		return fmt.Errorf("model is required")
	}
	if cfg.MaxConversationTurns < 0 {
		return fmt.Errorf("max conversation turns must be non-negative")
	}
	return nil
}

// RegisterTool registers a tool with the runtime
func (r *Runtime) RegisterTool(tool *Tool) error {
	if tool.Name == "" {
		return fmt.Errorf("tool name is required")
	}
	if tool.Handler == nil {
		return fmt.Errorf("tool handler is required")
	}

	r.toolsMu.Lock()
	defer r.toolsMu.Unlock()

	if _, exists := r.tools[tool.Name]; exists {
		return fmt.Errorf("tool %s already registered", tool.Name)
	}

	r.tools[tool.Name] = tool

	// Adapter for backward compatibility
	adapter := &AdkToolAdapter{tool: tool}
	r.adkTools = append(r.adkTools, adapter)

	logger.Infow("tool registered", "name", tool.Name)
	return nil
}

// GetTool returns a registered tool by name
func (r *Runtime) GetTool(name string) (*Tool, bool) {
	r.toolsMu.RLock()
	defer r.toolsMu.RUnlock()
	tool, ok := r.tools[name]
	return tool, ok
}

// ListTools returns all registered tools
func (r *Runtime) ListTools() []*Tool {
	r.toolsMu.RLock()
	defer r.toolsMu.RUnlock()

	tools := make([]*Tool, 0, len(r.tools))
	for _, t := range r.tools {
		tools = append(tools, t)
	}
	return tools
}

// ExecuteTool executes a tool by name with given parameters
func (r *Runtime) ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	tool, ok := r.GetTool(name)
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	logger.Infow("executing tool", "name", name)
	return tool.Handler(ctx, params)
}

// Run executes the agent with input and streams results
func (r *Runtime) Run(ctx context.Context, sessionKey string, input string, events chan<- StreamEvent) error {
	defer close(events)

	// 1. Get Provider (Already injected)
	p := r.provider

	// 2. Load or Create Session
	sess, err := r.sessionStore.Get(sessionKey)
	if err != nil {
		sess = &session.Session{
			Key:     sessionKey,
			History: []session.Message{},
			Model:   r.config.Model,
		}

		// Prepend system prompt if available
		if r.config.SystemPrompt != "" {
			sysMsg := session.Message{
				Role:      "system",
				Content:   r.config.SystemPrompt,
				Timestamp: time.Now(),
			}
			// Note: We don't AppendMessage here because Create(sess) below will persist the entire history
			sess.History = append(sess.History, sysMsg)
		}

		if err := r.sessionStore.Create(sess); err != nil {
			events <- StreamEvent{Type: "error", Error: err}
			return err
		}
	}

	// 3. Append User Message
	userMsg := session.Message{
		Role:      "user",
		Content:   input,
		Timestamp: time.Now(),
	}
	if err := r.sessionStore.AppendMessage(sessionKey, userMsg); err != nil {
		events <- StreamEvent{Type: "error", Error: err}
		return err
	}
	sess.History = append(sess.History, userMsg)

	// 4. Generation Loop
	maxTurns := 10
	currentTurn := 0

	for currentTurn < maxTurns {
		currentTurn++

		var providerMsgs []provider.Message
		for _, m := range sess.History {
			providerMsgs = append(providerMsgs, provider.Message{
				Role:    m.Role,
				Content: m.Content,
			})
		}

		var providerTools []provider.Tool
		r.toolsMu.RLock()
		for _, t := range r.tools {
			providerTools = append(providerTools, provider.Tool{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
			})
		}
		r.toolsMu.RUnlock()

		params := provider.GenerateParams{
			Model:       r.config.Model,
			Messages:    providerMsgs,
			Tools:       providerTools,
			MaxTokens:   r.config.MaxTokens,
			Temperature: r.config.Temperature,
		}

		stream, err := p.Generate(ctx, params)
		if err != nil {
			events <- StreamEvent{Type: "error", Error: err}
			return err
		}

		fullResponseText := ""
		var toolCalls []*provider.ToolCall
		toolCalled := false

		// Streaming Loop
		for event, err := range stream {
			if err != nil {
				events <- StreamEvent{Type: "error", Error: err}
				return err
			}

			switch event.Type {
			case provider.StreamEventPlainText:
				fullResponseText += event.Text
				events <- StreamEvent{Type: "text_delta", Text: event.Text}
			case provider.StreamEventToolCall:
				if event.ToolCall != nil {
					toolCalls = append(toolCalls, event.ToolCall)
					toolCalled = true
					events <- StreamEvent{Type: "tool_start", ToolCall: event.ToolCall}
				}
			case provider.StreamEventError:
				events <- StreamEvent{Type: "error", Error: event.Error}
				return event.Error
			}
		}

		// Save response
		if fullResponseText != "" {
			asstMsg := session.Message{
				Role:      "assistant",
				Content:   fullResponseText,
				Timestamp: time.Now(),
			}
			if err := r.sessionStore.AppendMessage(sessionKey, asstMsg); err != nil {
				logger.Warnw("failed to save response", "error", err)
			}
			sess.History = append(sess.History, asstMsg)
		}

		// Execute Tools
		if toolCalled {
			for _, tc := range toolCalls {
				var args map[string]interface{}
				if err := json.Unmarshal([]byte(tc.Arguments), &args); err != nil {
					logger.Errorw("failed to parse tool arguments", "tool", tc.Name, "args", tc.Arguments, "error", err)
					args = make(map[string]interface{})
				}

				tool, ok := r.GetTool(tc.Name)
				var result interface{}
				var err error
				if !ok {
					err = fmt.Errorf("tool not found: %s", tc.Name)
				} else {
					result, err = tool.Handler(ctx, args)
				}

				resultStr := ""
				if err != nil {
					resultStr = fmt.Sprintf("Error: %v", err)
				} else {
					resultJSON, _ := json.Marshal(result)
					resultStr = string(resultJSON)
				}

				events <- StreamEvent{
					Type: "tool_end",
					ToolCall: &provider.ToolCall{
						ID:   tc.ID,
						Name: tc.Name,
					},
					Text: resultStr,
				}

				toolMsg := session.Message{
					Role:      "tool",
					Content:   resultStr,
					Timestamp: time.Now(),
				}
				if err := r.sessionStore.AppendMessage(sessionKey, toolMsg); err != nil {
					logger.Warnw("failed to save tool result", "error", err)
				}
				sess.History = append(sess.History, toolMsg)
			}
			continue
		}

		break
	}

	events <- StreamEvent{Type: "done"}
	return nil
}

// AdkToolAdapter adapts our Tool to ADK's tool.Tool interface (kept for testing/compatibility)
type AdkToolAdapter struct {
	tool *Tool
}

// Name returns the tool name
func (a *AdkToolAdapter) Name() string {
	return a.tool.Name
}

// Description returns the tool description
func (a *AdkToolAdapter) Description() string {
	return a.tool.Description
}

// Run executes the tool with the given input
func (a *AdkToolAdapter) Run(ctx context.Context, input any) (any, error) {
	// Debug logging
	logger.Infow("AdkToolAdapter.Run called", "tool", a.tool.Name, "input", input)

	params, ok := input.(map[string]interface{})
	if !ok {
		if input == nil {
			params = make(map[string]interface{})
		} else {
			return nil, fmt.Errorf("invalid input type for tool %s: expected map[string]interface{}, got %T", a.tool.Name, input)
		}
	}

	return a.tool.Handler(ctx, params)
}

// IsLongRunning returns true if the tool takes a long time to run
func (a *AdkToolAdapter) IsLongRunning() bool {
	return false
}
