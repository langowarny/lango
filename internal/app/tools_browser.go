package app

import (
	"context"
	"fmt"
	"time"

	"github.com/langoai/lango/internal/agent"
	"github.com/langoai/lango/internal/tools/browser"
)

func buildBrowserTools(sm *browser.SessionManager) []*agent.Tool {
	return []*agent.Tool{
		{
			Name:        "browser_navigate",
			Description: "Navigate the browser to a URL and return the page title, URL, and a text snippet",
			SafetyLevel: agent.SafetyLevelDangerous,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "The URL to navigate to",
					},
				},
				"required": []string{"url"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				url, ok := params["url"].(string)
				if !ok || url == "" {
					return nil, fmt.Errorf("missing url parameter")
				}

				sessionID, err := sm.EnsureSession()
				if err != nil {
					return nil, err
				}

				if err := sm.Tool().Navigate(ctx, sessionID, url); err != nil {
					return nil, err
				}

				return sm.Tool().GetSnapshot(sessionID)
			},
		},
		{
			Name:        "browser_action",
			Description: "Perform an action on the current browser page: click, type, eval, get_text, get_element_info, or wait",
			SafetyLevel: agent.SafetyLevelDangerous,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"description": "The action to perform",
						"enum":        []string{"click", "type", "eval", "get_text", "get_element_info", "wait"},
					},
					"selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector for the target element (required for click, type, get_text, get_element_info, wait)",
					},
					"text": map[string]interface{}{
						"type":        "string",
						"description": "Text to type (required for type action) or JavaScript to evaluate (required for eval action)",
					},
					"timeout": map[string]interface{}{
						"type":        "integer",
						"description": "Timeout in seconds for wait action (default: 10)",
					},
				},
				"required": []string{"action"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				action, ok := params["action"].(string)
				if !ok || action == "" {
					return nil, fmt.Errorf("missing action parameter")
				}

				sessionID, err := sm.EnsureSession()
				if err != nil {
					return nil, err
				}

				selector, _ := params["selector"].(string)
				text, _ := params["text"].(string)

				switch action {
				case "click":
					if selector == "" {
						return nil, fmt.Errorf("selector required for click action")
					}
					return nil, sm.Tool().Click(ctx, sessionID, selector)

				case "type":
					if selector == "" {
						return nil, fmt.Errorf("selector required for type action")
					}
					if text == "" {
						return nil, fmt.Errorf("text required for type action")
					}
					return nil, sm.Tool().Type(ctx, sessionID, selector, text)

				case "eval":
					if text == "" {
						return nil, fmt.Errorf("text (JavaScript) required for eval action")
					}
					return sm.Tool().Eval(sessionID, text)

				case "get_text":
					if selector == "" {
						return nil, fmt.Errorf("selector required for get_text action")
					}
					return sm.Tool().GetText(sessionID, selector)

				case "get_element_info":
					if selector == "" {
						return nil, fmt.Errorf("selector required for get_element_info action")
					}
					return sm.Tool().GetElementInfo(sessionID, selector)

				case "wait":
					if selector == "" {
						return nil, fmt.Errorf("selector required for wait action")
					}
					timeout := 10 * time.Second
					if t, ok := params["timeout"].(float64); ok && t > 0 {
						timeout = time.Duration(t) * time.Second
					}
					return nil, sm.Tool().WaitForSelector(ctx, sessionID, selector, timeout)

				default:
					return nil, fmt.Errorf("unknown action: %s", action)
				}
			},
		},
		{
			Name:        "browser_screenshot",
			Description: "Capture a screenshot of the current browser page as base64 PNG",
			SafetyLevel: agent.SafetyLevelSafe,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"fullPage": map[string]interface{}{
						"type":        "boolean",
						"description": "Capture the full scrollable page (default: false)",
					},
				},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				sessionID, err := sm.EnsureSession()
				if err != nil {
					return nil, err
				}

				fullPage, _ := params["fullPage"].(bool)
				return sm.Tool().Screenshot(sessionID, fullPage)
			},
		},
	}
}
