package adk

import (
	"github.com/google/jsonschema-go/jsonschema"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/langowarny/lango/internal/agent"
)

// AdaptTool converts an internal agent.Tool to an ADK tool.Tool
func AdaptTool(t *agent.Tool) (tool.Tool, error) {
	// Build input schema from parameters
	props := make(map[string]*jsonschema.Schema)
	var required []string

	for name, paramDef := range t.Parameters {
		s := &jsonschema.Schema{}

		// Attempt to parse ParameterDef
		// Since it is stored as interface{}, we need to handle potential map conversions if it came from JSON
		// But in-memory tools usually use the struct.
		if pd, ok := paramDef.(agent.ParameterDef); ok {
			s.Type = pd.Type
			s.Description = pd.Description
			if len(pd.Enum) > 0 {
				s.Enum = make([]any, len(pd.Enum))
				for i, v := range pd.Enum {
					s.Enum[i] = v
				}
			}
			if pd.Required {
				required = append(required, name)
			}
		} else if pdMap, ok := paramDef.(map[string]interface{}); ok {
			// Handle map (e.g. from JSON config)
			if t, ok := pdMap["type"].(string); ok {
				s.Type = t
			}
			if d, ok := pdMap["description"].(string); ok {
				s.Description = d
			}
			if r, ok := pdMap["required"].(bool); ok && r {
				required = append(required, name)
			}
		} else {
			// Fallback or skip
			s.Type = "string" // default
		}
		props[name] = s
	}

	inputSchema := &jsonschema.Schema{
		Type:       "object",
		Properties: props,
		Required:   required,
	}

	cfg := functiontool.Config{
		Name:        t.Name,
		Description: t.Description,
		InputSchema: inputSchema,
	}

	// Wrapper handler
	handler := func(ctx tool.Context, args map[string]any) (any, error) {
		return t.Handler(ctx, args)
	}

	return functiontool.New(cfg, handler)
}
