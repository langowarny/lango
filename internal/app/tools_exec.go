package app

import (
	"context"
	"fmt"

	"github.com/langoai/lango/internal/agent"
	"github.com/langoai/lango/internal/supervisor"
)

func buildExecTools(sv *supervisor.Supervisor, automationAvailable map[string]bool) []*agent.Tool {
	return []*agent.Tool{
		{
			Name:        "exec",
			Description: "Execute shell commands",
			SafetyLevel: agent.SafetyLevelDangerous,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "The shell command to execute",
					},
				},
				"required": []string{"command"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				cmd, ok := params["command"].(string)
				if !ok {
					return nil, fmt.Errorf("missing command parameter")
				}
				if msg := blockLangoExec(cmd, automationAvailable); msg != "" {
					return map[string]interface{}{"blocked": true, "message": msg}, nil
				}
				return sv.ExecuteTool(ctx, cmd)
			},
		},
		{
			Name:        "exec_bg",
			Description: "Execute a shell command in the background",
			SafetyLevel: agent.SafetyLevelDangerous,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "The shell command to execute",
					},
				},
				"required": []string{"command"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				cmd, ok := params["command"].(string)
				if !ok {
					return nil, fmt.Errorf("missing command parameter")
				}
				if msg := blockLangoExec(cmd, automationAvailable); msg != "" {
					return map[string]interface{}{"blocked": true, "message": msg}, nil
				}
				return sv.StartBackground(cmd)
			},
		},
		{
			Name:        "exec_status",
			Description: "Check the status of a background process",
			SafetyLevel: agent.SafetyLevelSafe,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "The background process ID returned by exec_bg",
					},
				},
				"required": []string{"id"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				id, ok := params["id"].(string)
				if !ok {
					return nil, fmt.Errorf("missing id parameter")
				}
				return sv.GetBackgroundStatus(id)
			},
		},
		{
			Name:        "exec_stop",
			Description: "Stop a background process",
			SafetyLevel: agent.SafetyLevelDangerous,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "The background process ID returned by exec_bg",
					},
				},
				"required": []string{"id"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				id, ok := params["id"].(string)
				if !ok {
					return nil, fmt.Errorf("missing id parameter")
				}
				return nil, sv.StopBackground(id)
			},
		},
	}
}
