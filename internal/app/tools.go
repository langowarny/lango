package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/langoai/lango/internal/agent"
	"github.com/langoai/lango/internal/approval"
	"github.com/langoai/lango/internal/config"
	"github.com/langoai/lango/internal/learning"
	"github.com/langoai/lango/internal/session"
	"github.com/langoai/lango/internal/supervisor"
	"github.com/langoai/lango/internal/tools/browser"
	"github.com/langoai/lango/internal/tools/filesystem"
	"github.com/langoai/lango/internal/types"
	"github.com/langoai/lango/internal/wallet"
)

// buildTools creates the set of tools available to the agent.
// When browserSM is non-nil, browser tools are included.
// automationAvailable indicates which automation features are enabled (cron, background, workflow).
func buildTools(sv *supervisor.Supervisor, fsCfg filesystem.Config, browserSM *browser.SessionManager, automationAvailable map[string]bool) []*agent.Tool {
	var tools []*agent.Tool

	// Exec tools (delegated to Supervisor for security isolation)
	tools = append(tools, buildExecTools(sv, automationAvailable)...)

	// Filesystem tools
	fsTool := filesystem.New(fsCfg)
	tools = append(tools, buildFilesystemTools(fsTool)...)

	// Browser tools (opt-in), wrapped with panic recovery
	if browserSM != nil {
		for _, bt := range buildBrowserTools(browserSM) {
			tools = append(tools, wrapBrowserHandler(bt, browserSM))
		}
	}

	return tools
}

// blockLangoExec checks if the command attempts to invoke the lango CLI.
// ALL lango CLI commands require passphrase authentication via bootstrap and
// will fail when spawned as a subprocess (non-interactive stdin). Returns a
// guidance message if blocked, or empty string if allowed.
func blockLangoExec(cmd string, automationAvailable map[string]bool) string {
	lower := strings.ToLower(strings.TrimSpace(cmd))

	// --- Phase 1: Subcommands with in-process tool equivalents ---
	type guard struct {
		prefix  string
		feature string // key in automationAvailable; empty = always available
		tools   string
	}
	guards := []guard{
		{"lango cron", "cron", "cron_add, cron_list, cron_pause, cron_resume, cron_remove, cron_history"},
		{"lango bg", "background", "bg_submit, bg_status, bg_list, bg_result, bg_cancel"},
		{"lango background", "background", "bg_submit, bg_status, bg_list, bg_result, bg_cancel"},
		{"lango workflow", "workflow", "workflow_run, workflow_status, workflow_list, workflow_cancel, workflow_save"},
		{"lango graph", "", "graph_traverse, graph_query, rag_retrieve"},
		{"lango memory", "", "memory_list_observations, memory_list_reflections"},
		{"lango p2p", "", "p2p_status, p2p_connect, p2p_disconnect, p2p_peers, p2p_query, p2p_discover, p2p_firewall_rules, p2p_firewall_add, p2p_firewall_remove, p2p_reputation, p2p_pay, p2p_price_query"},
		{"lango security", "", "crypto_encrypt, crypto_decrypt, crypto_sign, crypto_hash, crypto_keys, secrets_store, secrets_get, secrets_list, secrets_delete"},
		{"lango payment", "", "payment_send, payment_create_wallet, payment_x402_fetch"},
	}

	for _, g := range guards {
		if strings.HasPrefix(lower, g.prefix) {
			if g.feature == "" || automationAvailable[g.feature] {
				return fmt.Sprintf(
					"Do not use exec to run '%s' — use the built-in tools instead (%s). "+
						"Spawning a new lango process requires passphrase authentication and will fail in non-interactive mode.",
					g.prefix, g.tools)
			}
			return fmt.Sprintf(
				"Cannot run '%s' via exec — spawning a new lango process requires passphrase authentication. "+
					"Enable the %s feature in Settings to use the built-in tools (%s).",
				g.prefix, g.feature, g.tools)
		}
	}

	// --- Phase 2: Catch-all for any remaining lango subcommand ---
	if strings.HasPrefix(lower, "lango ") || lower == "lango" {
		return "Do not use exec to run the lango CLI — every lango command requires passphrase authentication " +
			"via bootstrap and will fail when spawned as a subprocess. " +
			"Use the built-in tools for the operation you need, or ask the user to run this command directly in their terminal."
	}

	// Redirect skill-related git clone to import_skill tool.
	if strings.HasPrefix(lower, "git clone") && strings.Contains(lower, "skill") {
		return "Use the built-in import_skill tool instead of manual git clone — " +
			"it automatically uses git clone internally when available and stores skills in the correct location (~/.lango/skills/). " +
			"Example: import_skill(url: \"<github-repo-url>\")"
	}

	// Redirect skill-related curl/wget to import_skill tool.
	if (strings.HasPrefix(lower, "curl ") || strings.HasPrefix(lower, "wget ")) &&
		strings.Contains(lower, "skill") {
		return "Use the built-in import_skill tool instead of manual curl/wget — " +
			"it handles downloads internally and stores skills correctly. " +
			"Example: import_skill(url: \"<url>\")"
	}

	return ""
}

// wrapBrowserHandler wraps a browser tool handler with panic recovery and auto-reconnect.
// On panic, it converts to an error. On ErrBrowserPanic, it closes the session and retries once.
func wrapBrowserHandler(t *agent.Tool, sm *browser.SessionManager) *agent.Tool {
	original := t.Handler
	return &agent.Tool{
		Name:        t.Name,
		Description: t.Description,
		Parameters:  t.Parameters,
		SafetyLevel: t.SafetyLevel,
		Handler: func(ctx context.Context, params map[string]interface{}) (result interface{}, retErr error) {
			defer func() {
				if r := recover(); r != nil {
					logger().Errorw("browser tool panic recovered", "tool", t.Name, "panic", r)
					retErr = fmt.Errorf("%w: %v", browser.ErrBrowserPanic, r)
				}
			}()

			result, retErr = original(ctx, params)
			if retErr != nil && errors.Is(retErr, browser.ErrBrowserPanic) {
				// Connection likely dead — close and retry once
				logger().Warnw("browser panic detected, closing session and retrying", "tool", t.Name, "error", retErr)
				_ = sm.Close()
				result, retErr = original(ctx, params)
			}
			return
		},
	}
}

// wrapWithLearning wraps a tool's handler to call the learning observer after each execution.
func wrapWithLearning(t *agent.Tool, observer learning.ToolResultObserver) *agent.Tool {
	original := t.Handler
	return &agent.Tool{
		Name:        t.Name,
		Description: t.Description,
		Parameters:  t.Parameters,
		SafetyLevel: t.SafetyLevel,
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			result, err := original(ctx, params)
			sessionKey := session.SessionKeyFromContext(ctx)
			observer.OnToolResult(ctx, sessionKey, t.Name, params, result, err)
			return result, err
		},
	}
}

// detectChannelFromContext extracts the delivery target from the session key in context.
// Returns "channel:targetID" (e.g. "telegram:123456789") or "" if no known channel prefix is found.
func detectChannelFromContext(ctx context.Context) string {
	sessionKey := session.SessionKeyFromContext(ctx)
	if sessionKey == "" {
		return ""
	}
	// Session key format: "channel:targetID:userID"
	parts := strings.SplitN(sessionKey, ":", 3)
	if len(parts) < 2 {
		return ""
	}
	ch := types.ChannelType(parts[0])
	if ch.Valid() {
		return parts[0] + ":" + parts[1]
	}
	return ""
}

// needsApproval determines whether a tool requires approval based on the
// configured policy, explicit exemptions, and sensitive tool lists.
func needsApproval(t *agent.Tool, ic config.InterceptorConfig) bool {
	// ExemptTools always bypass approval.
	for _, name := range ic.ExemptTools {
		if name == t.Name {
			return false
		}
	}

	// SensitiveTools always require approval.
	for _, name := range ic.SensitiveTools {
		if name == t.Name {
			return true
		}
	}

	switch ic.ApprovalPolicy {
	case config.ApprovalPolicyAll:
		return true
	case config.ApprovalPolicyDangerous:
		return t.SafetyLevel.IsDangerous()
	case config.ApprovalPolicyConfigured:
		return false // only SensitiveTools (handled above)
	case config.ApprovalPolicyNone:
		return false
	default:
		return true // unknown policy → fail-safe
	}
}

// buildApprovalSummary returns a human-readable description of what a tool
// invocation will do, suitable for display in approval messages.
func buildApprovalSummary(toolName string, params map[string]interface{}) string {
	switch toolName {
	case "exec", "exec_bg":
		if cmd, ok := params["command"].(string); ok {
			return "Execute: " + truncate(cmd, 200)
		}
	case "fs_write":
		path, _ := params["path"].(string)
		content, _ := params["content"].(string)
		return fmt.Sprintf("Write to %s (%d bytes)", path, len(content))
	case "fs_edit":
		path, _ := params["path"].(string)
		return "Edit file: " + path
	case "fs_delete":
		path, _ := params["path"].(string)
		return "Delete: " + path
	case "browser_navigate":
		url, _ := params["url"].(string)
		return "Navigate to: " + truncate(url, 200)
	case "browser_action":
		action, _ := params["action"].(string)
		selector, _ := params["selector"].(string)
		if selector != "" {
			return fmt.Sprintf("Browser %s on: %s", action, truncate(selector, 100))
		}
		return "Browser action: " + action
	case "secrets_store":
		name, _ := params["name"].(string)
		return "Store secret: " + name
	case "secrets_get":
		name, _ := params["name"].(string)
		return "Retrieve secret: " + name
	case "secrets_delete":
		name, _ := params["name"].(string)
		return "Delete secret: " + name
	case "crypto_encrypt":
		return "Encrypt data"
	case "crypto_decrypt":
		return "Decrypt ciphertext"
	case "crypto_sign":
		return "Generate digital signature"
	case "payment_send":
		amount, _ := params["amount"].(string)
		to, _ := params["to"].(string)
		purpose, _ := params["purpose"].(string)
		return fmt.Sprintf("Send %s USDC to %s (%s)", amount, truncate(to, 12), truncate(purpose, 50))
	case "payment_create_wallet":
		return "Create new blockchain wallet"
	case "payment_x402_fetch":
		url, _ := params["url"].(string)
		method, _ := params["method"].(string)
		if method == "" {
			method = "GET"
		}
		return fmt.Sprintf("X402 %s %s (auto-pay enabled)", method, truncate(url, 150))
	case "cron_add":
		name, _ := params["name"].(string)
		scheduleType, _ := params["schedule_type"].(string)
		schedule, _ := params["schedule"].(string)
		return fmt.Sprintf("Create cron job: %s (%s=%s)", name, scheduleType, schedule)
	case "cron_remove":
		id, _ := params["id"].(string)
		return "Remove cron job: " + id
	case "bg_submit":
		prompt, _ := params["prompt"].(string)
		return "Submit background task: " + truncate(prompt, 100)
	case "workflow_run":
		filePath, _ := params["file_path"].(string)
		if filePath != "" {
			return "Run workflow: " + filePath
		}
		return "Run inline workflow"
	case "workflow_cancel":
		runID, _ := params["run_id"].(string)
		return "Cancel workflow: " + runID
	case "p2p_pay":
		amount, _ := params["amount"].(string)
		peerDID, _ := params["peer_did"].(string)
		memo, _ := params["memo"].(string)
		if memo == "" {
			memo = "P2P payment"
		}
		return fmt.Sprintf("Pay %s USDC to peer %s (%s)", amount, truncate(peerDID, 16), truncate(memo, 50))
	}
	return "Tool: " + toolName
}

// truncate shortens s to maxLen characters, appending "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// wrapWithApproval wraps a tool to require approval based on the configured policy.
// Uses fail-closed: denies execution unless explicitly approved.
// The approval.Provider routes requests to the appropriate channel (Gateway, Telegram, Discord, Slack, TTY).
// The GrantStore tracks "always allow" grants to auto-approve repeat invocations within a session.
// When limiter is non-nil, payment tools with an amount below the auto-approve threshold
// are executed without explicit user confirmation.
func wrapWithApproval(t *agent.Tool, ic config.InterceptorConfig, ap approval.Provider, gs *approval.GrantStore, limiter wallet.SpendingLimiter) *agent.Tool {
	if !needsApproval(t, ic) {
		return t
	}

	original := t.Handler
	return &agent.Tool{
		Name:        t.Name,
		Description: t.Description,
		Parameters:  t.Parameters,
		SafetyLevel: t.SafetyLevel,
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			sessionKey := session.SessionKeyFromContext(ctx)
			if target := approval.ApprovalTargetFromContext(ctx); target != "" {
				sessionKey = target
			}

			// Check persistent grant — auto-approve if previously "always allowed"
			if gs != nil && gs.IsGranted(sessionKey, t.Name) {
				return original(ctx, params)
			}

			// Auto-approve small payments via spending limiter threshold.
			if limiter != nil && (t.Name == "p2p_pay" || t.Name == "payment_send") {
				if amountStr, ok := params["amount"].(string); ok && amountStr != "" {
					amt, err := wallet.ParseUSDC(amountStr)
					if err == nil {
						if autoOK, checkErr := limiter.IsAutoApprovable(ctx, amt); checkErr == nil && autoOK {
							return original(ctx, params)
						}
					}
				}
			}

			req := approval.ApprovalRequest{
				ID:         fmt.Sprintf("req-%d", time.Now().UnixNano()),
				ToolName:   t.Name,
				SessionKey: sessionKey,
				Params:     params,
				Summary:    buildApprovalSummary(t.Name, params),
				CreatedAt:  time.Now(),
			}
			resp, err := ap.RequestApproval(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("tool '%s' approval: %w", t.Name, err)
			}
			if !resp.Approved {
				sk := session.SessionKeyFromContext(ctx)
				if sk == "" {
					return nil, fmt.Errorf("tool '%s' execution denied: no approval channel available (session key missing)", t.Name)
				}
				return nil, fmt.Errorf("tool '%s' execution denied: user did not approve the action", t.Name)
			}

			// Record persistent grant for this session+tool
			if resp.AlwaysAllow && gs != nil {
				gs.Grant(sessionKey, t.Name)
			}

			return original(ctx, params)
		},
	}
}
