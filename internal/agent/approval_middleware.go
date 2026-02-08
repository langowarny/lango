package agent

import (
	"context"
	"fmt"
	"slices"
)

// ApprovalMiddleware implements AgentRuntime to require approval for sensitive tools.
type ApprovalMiddleware struct {
	BaseRuntimeMiddleware
	Config ApprovalConfig
}

// NewApprovalMiddleware creates a new ApprovalMiddleware.
func NewApprovalMiddleware(cfg ApprovalConfig) RuntimeMiddleware {
	return func(next AgentRuntime) AgentRuntime {
		return &ApprovalMiddleware{
			BaseRuntimeMiddleware: BaseRuntimeMiddleware{Next: next},
			Config:                cfg,
		}
	}
}

// RegisterTool intercepts tool registration to wrap handlers with approval logic.
func (mw *ApprovalMiddleware) RegisterTool(tool *Tool) error {
	requiresApproval := false
	var restrictedOps []string

	// Check for full tool blocking
	if slices.Contains(mw.Config.SensitiveTools, tool.Name) {
		requiresApproval = true
	} else {
		// Check for operation-specific blocking
		prefix := tool.Name + "."
		for _, s := range mw.Config.SensitiveTools {
			if len(s) > len(prefix) && s[:len(prefix)] == prefix {
				restrictedOps = append(restrictedOps, s[len(prefix):])
			}
		}
	}

	if requiresApproval || len(restrictedOps) > 0 {
		originalHandler := tool.Handler

		// Wrap handler
		tool.Handler = func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			// Check if this execution requires approval
			shouldBlock := requiresApproval
			if !shouldBlock && len(restrictedOps) > 0 {
				if op, ok := params["operation"].(string); ok {
					if slices.Contains(restrictedOps, op) {
						shouldBlock = true
					}
				}
			}

			if shouldBlock {
				logger.Infow("intercepting sensitive tool execution", "tool", tool.Name, "params", params)

				if mw.Config.NotifyChannel != nil {
					msg := fmt.Sprintf("Approval required for tool '%s' with params: %v", tool.Name, params)
					approved, err := mw.Config.NotifyChannel(ctx, msg)
					if err != nil {
						return nil, fmt.Errorf("approval check failed: %w", err)
					}
					if !approved {
						return nil, fmt.Errorf("execution denied by user")
					}
				} else {
					// Fail safe if no notifier configured but approval required
					logger.Warnw("approval required but no notify channel configured", "tool", tool.Name)
					return nil, fmt.Errorf("approval required but system not configured for notifications")
				}
				logger.Infow("tool execution approved", "tool", tool.Name)
			}

			return originalHandler(ctx, params)
		}
	}

	return mw.Next.RegisterTool(tool)
}
