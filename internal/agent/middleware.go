package agent

import (
	"context"
)

// AgentRuntime defines the interface for an agent runtime.
type AgentRuntime interface {
	Run(ctx context.Context, sessionKey string, input string, events chan<- StreamEvent) error
	RegisterTool(tool *Tool) error
	GetTool(name string) (*Tool, bool)
	ListTools() []*Tool
	ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error)
}

// Ensure Runtime implements AgentRuntime
var _ AgentRuntime = (*Runtime)(nil)

// RuntimeMiddleware defines a function that wraps an AgentRuntime.
type RuntimeMiddleware func(next AgentRuntime) AgentRuntime

// BaseRuntimeMiddleware allows wrapping a Runtime to intercept calls.
type BaseRuntimeMiddleware struct {
	Next AgentRuntime
}

func (m *BaseRuntimeMiddleware) RegisterTool(tool *Tool) error {
	return m.Next.RegisterTool(tool)
}

func (m *BaseRuntimeMiddleware) Run(ctx context.Context, sessionKey string, input string, events chan<- StreamEvent) error {
	return m.Next.Run(ctx, sessionKey, input, events)
}

func (m *BaseRuntimeMiddleware) GetTool(name string) (*Tool, bool) {
	return m.Next.GetTool(name)
}

func (m *BaseRuntimeMiddleware) ListTools() []*Tool {
	return m.Next.ListTools()
}

func (m *BaseRuntimeMiddleware) ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	return m.Next.ExecuteTool(ctx, name, params)
}

// ChainMiddleware applies multiple middlewares to a Runtime.
func ChainMiddleware(r AgentRuntime, middlewares ...RuntimeMiddleware) AgentRuntime {
	for i := len(middlewares) - 1; i >= 0; i-- {
		r = middlewares[i](r)
	}
	return r
}

// PIIConfig defines configuration for PII redaction
type PIIConfig struct {
	RedactEmail bool
	RedactPhone bool
	CustomRegex []string
}

// ApprovalConfig defines configuration for tool approval
type ApprovalConfig struct {
	SensitiveTools []string
	NotifyChannel  func(ctx context.Context, msg string) (bool, error)
}
