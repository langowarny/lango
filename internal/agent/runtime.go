package agent

import (
	"context"
)

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
