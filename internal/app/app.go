package app

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/langowarny/lango/internal/agent"
	"github.com/langowarny/lango/internal/cli/prompt"
	"github.com/langowarny/lango/internal/config"
	"github.com/langowarny/lango/internal/gateway"
	"github.com/langowarny/lango/internal/logging"
	"github.com/langowarny/lango/internal/security"
	"github.com/langowarny/lango/internal/session"
	"github.com/langowarny/lango/internal/supervisor"
	"github.com/langowarny/lango/internal/tools/browser"
	"github.com/langowarny/lango/internal/tools/crypto"

	// exec tool imported by supervisor
	"github.com/langowarny/lango/internal/tools/filesystem"
	"github.com/langowarny/lango/internal/tools/secrets"
)

var logger = logging.SubsystemSugar("app")

const DefaultSystemPrompt = "You are Lango, a powerful AI assistant. You have access to tools for web navigation (browser), secure secrets management (secrets), and cryptographic operations (crypto). Use them when appropriate to help the user."

// New creates a new application instance
func New(cfg *config.Config) (*App, error) {
	app := &App{
		Config: cfg,
	}

	// Initialize Supervisor (Holds Secrets)
	sv, err := supervisor.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create supervisor: %w", err)
	}

	// Initialize Session Store (Ent)
	// Prioritize LANGO_PASSPHRASE env var, fall back to config
	passphrase := os.Getenv("LANGO_PASSPHRASE")
	if passphrase == "" {
		passphrase = cfg.Security.Passphrase
	}

	var storeOpts []session.StoreOption
	if passphrase != "" {
		storeOpts = append(storeOpts, session.WithPassphrase(passphrase))
	}

	store, err := session.NewEntStore(cfg.Session.DatabasePath, storeOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create session store: %w", err)
	}
	app.Store = store

	// Provider initialization moved to Supervisor

	// Legacy Config Support check - handled by Supervisor

	// Initialize Agent
	agentCfg := agent.Config{
		Provider: cfg.Agent.Provider,
		Model:    cfg.Agent.Model,
		// APIKey removed
		MaxConversationTurns: cfg.Session.MaxHistoryTurns, // Use session config for history limit
		MaxTokens:            cfg.Agent.MaxTokens,
		Temperature:          cfg.Agent.Temperature,
		FallbackProvider:     cfg.Agent.FallbackProvider,
		FallbackModel:        cfg.Agent.FallbackModel,
		SystemPrompt:         DefaultSystemPrompt,
	}

	// Create Provider Proxy
	proxy := supervisor.NewProviderProxy(sv, cfg.Agent.Provider, cfg.Agent.Model)

	baseRuntime, err := agent.New(agentCfg, store, proxy)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent runtime: %w", err)
	}

	// Apply Middlewares
	var middlewares []agent.RuntimeMiddleware

	// 1. PII Redactor
	if cfg.Security.Interceptor.Enabled && cfg.Security.Interceptor.RedactPII {
		piiCfg := agent.PIIConfig{
			RedactEmail: true,
			RedactPhone: true,
			CustomRegex: cfg.Security.Interceptor.PIIRegexPatterns,
		}
		middlewares = append(middlewares, agent.NewPIIRedactor(piiCfg))
	}

	// 2. Approval Middleware
	if cfg.Security.Interceptor.Enabled && cfg.Security.Interceptor.ApprovalRequired {
		// Default sensitive tools that always require approval
		sensitiveTools := []string{"secrets.get"}
		// Append user-defined tools
		sensitiveTools = append(sensitiveTools, cfg.Security.Interceptor.SensitiveTools...)

		approvalCfg := agent.ApprovalConfig{
			SensitiveTools: sensitiveTools,
			// NotifyChannel must be wired to actual channel logic.
			// Ideally we have a ChannelManager to send messages.
			// For now, we will use a placeholder or log-only notifier if implementation is complex.
			// TODO: WIre this to actual Discord/Telegram channels via app.Channels
			NotifyChannel: func(ctx context.Context, msg string) (bool, error) {
				if app.Gateway == nil {
					return false, fmt.Errorf("gateway not initialized")
				}
				return app.Gateway.RequestApproval(ctx, msg)
			},
		}
		middlewares = append(middlewares, agent.NewApprovalMiddleware(approvalCfg))
	}

	app.Agent = agent.ChainMiddleware(baseRuntime, middlewares...)

	// Initialize Tools
	// Exec Tool - Delegated to Supervisor
	if err := app.Agent.RegisterTool(&agent.Tool{
		Name:        "exec",
		Description: "Execute shell commands",
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
			return sv.ExecuteTool(ctx, cmd)
		},
	}); err != nil {
		logger.Warnw("failed to register exec tool", "error", err)
	}

	if err := app.Agent.RegisterTool(&agent.Tool{
		Name:        "exec_bg",
		Description: "Execute a shell command in the background",
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
			return sv.StartBackground(cmd)
		},
	}); err != nil {
		logger.Warnw("failed to register exec_bg tool", "error", err)
	}

	if err := app.Agent.RegisterTool(&agent.Tool{
		Name:        "exec_status",
		Description: "Check the status of a background process",
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
	}); err != nil {
		logger.Warnw("failed to register exec_status tool", "error", err)
	}

	if err := app.Agent.RegisterTool(&agent.Tool{
		Name:        "exec_stop",
		Description: "Stop a background process",
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
	}); err != nil {
		logger.Warnw("failed to register exec_stop tool", "error", err)
	}

	// Browser Tool
	browserConfig := browser.Config{
		Headless:       cfg.Tools.Browser.Headless,
		SessionTimeout: cfg.Tools.Browser.SessionTimeout,
	}
	browserTool, err := browser.New(browserConfig)
	if err != nil {
		logger.Warnw("failed to initialize browser tool", "error", err)
	} else {
		if err := app.Agent.RegisterTool(&agent.Tool{
			Name:        "browser_navigate",
			Description: "Navigate to a URL and get page summary",
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
				if !ok {
					return nil, fmt.Errorf("missing url parameter")
				}
				sessionID, err := browserTool.NewSession()
				if err != nil {
					return nil, err
				}
				if err := browserTool.Navigate(ctx, sessionID, url); err != nil {
					return nil, err
				}
				app.BrowserSessionID = sessionID
				return browserTool.GetSnapshot(sessionID)
			},
		}); err != nil {
			logger.Warnw("failed to register browser_navigate tool", "error", err)
		}

		if err := app.Agent.RegisterTool(&agent.Tool{
			Name:        "browser_read",
			Description: "Read text content from the current browser session",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector to read (default: body)",
					},
				},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				if app.BrowserSessionID == "" {
					return nil, fmt.Errorf("no active browser session; call browser_navigate first")
				}
				selector, _ := params["selector"].(string)
				if selector == "" {
					selector = "body"
				}
				return browserTool.GetText(app.BrowserSessionID, selector)
			},
		}); err != nil {
			logger.Warnw("failed to register browser_read tool", "error", err)
		}

		if err := app.Agent.RegisterTool(&agent.Tool{
			Name:        "browser_screenshot",
			Description: "Capture a screenshot of the current page",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"fullPage": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether to capture the full scrollable page",
					},
				},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				if app.BrowserSessionID == "" {
					return nil, fmt.Errorf("no active browser session; call browser_navigate first")
				}
				fullPage, _ := params["fullPage"].(bool)
				return browserTool.Screenshot(app.BrowserSessionID, fullPage)
			},
		}); err != nil {
			logger.Warnw("failed to register browser_screenshot tool", "error", err)
		}
	}

	// Filesystem Tool
	fsConfig := filesystem.Config{
		MaxReadSize:  cfg.Tools.Filesystem.MaxReadSize,
		AllowedPaths: cfg.Tools.Filesystem.AllowedPaths,
	}
	fsTool := filesystem.New(fsConfig)
	if err := app.Agent.RegisterTool(&agent.Tool{
		Name:        "fs_read",
		Description: "Read a file",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The file path to read",
				},
			},
			"required": []string{"path"},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			path, ok := params["path"].(string)
			if !ok {
				return nil, fmt.Errorf("missing path parameter")
			}
			return fsTool.Read(path)
		},
	}); err != nil {
		logger.Warnw("failed to register filesystem tool", "error", err)
	}

	if err := app.Agent.RegisterTool(&agent.Tool{
		Name:        "fs_list",
		Description: "List contents of a directory",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The directory path to list",
				},
			},
			"required": []string{"path"},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			path, _ := params["path"].(string)
			if path == "" {
				path = "."
			}
			return fsTool.ListDir(path)
		},
	}); err != nil {
		logger.Warnw("failed to register fs_list tool", "error", err)
	}

	if err := app.Agent.RegisterTool(&agent.Tool{
		Name:        "fs_write",
		Description: "Write content to a file",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The file path to write to",
				},
				"content": map[string]interface{}{
					"type":        "string",
					"description": "The content to write",
				},
			},
			"required": []string{"path", "content"},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			path, _ := params["path"].(string)
			content, _ := params["content"].(string)
			if path == "" {
				return nil, fmt.Errorf("missing path parameter")
			}
			return nil, fsTool.Write(path, content)
		},
	}); err != nil {
		logger.Warnw("failed to register fs_write tool", "error", err)
	}

	if err := app.Agent.RegisterTool(&agent.Tool{
		Name:        "fs_edit",
		Description: "Edit a file by replacing a line range",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The file path to edit",
				},
				"startLine": map[string]interface{}{
					"type":        "integer",
					"description": "The starting line number (1-indexed)",
				},
				"endLine": map[string]interface{}{
					"type":        "integer",
					"description": "The ending line number (inclusive)",
				},
				"content": map[string]interface{}{
					"type":        "string",
					"description": "The new content for the specified range",
				},
			},
			"required": []string{"path", "startLine", "endLine", "content"},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			path, _ := params["path"].(string)
			content, _ := params["content"].(string)
			if path == "" {
				return nil, fmt.Errorf("missing path parameter")
			}

			var startLine, endLine int
			if sl, ok := params["startLine"].(float64); ok {
				startLine = int(sl)
			} else if sl, ok := params["startLine"].(int); ok {
				startLine = sl
			}

			if el, ok := params["endLine"].(float64); ok {
				endLine = int(el)
			} else if el, ok := params["endLine"].(int); ok {
				endLine = el
			}

			return nil, fsTool.Edit(path, startLine, endLine, content)
		},
	}); err != nil {
		logger.Warnw("failed to register fs_edit tool", "error", err)
	}

	if err := app.Agent.RegisterTool(&agent.Tool{
		Name:        "fs_mkdir",
		Description: "Create a directory",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The directory path to create",
				},
			},
			"required": []string{"path"},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			path, _ := params["path"].(string)
			if path == "" {
				return nil, fmt.Errorf("missing path parameter")
			}
			return nil, fsTool.Mkdir(path)
		},
	}); err != nil {
		logger.Warnw("failed to register fs_mkdir tool", "error", err)
	}

	if err := app.Agent.RegisterTool(&agent.Tool{
		Name:        "fs_delete",
		Description: "Delete a file or directory",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The path to delete",
				},
			},
			"required": []string{"path"},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			path, _ := params["path"].(string)
			if path == "" {
				return nil, fmt.Errorf("missing path parameter")
			}
			return nil, fsTool.Delete(path)
		},
	}); err != nil {
		logger.Warnw("failed to register fs_delete tool", "error", err)
	}

	// Initialize CryptoProvider
	var cryptoProvider security.CryptoProvider
	var rpcProvider *security.RPCProvider

	if cfg.Security.Signer.Provider == "rpc" {
		rpcProvider = security.NewRPCProvider()
		cryptoProvider = rpcProvider
		app.CryptoProvider = cryptoProvider
	} else if cfg.Security.Signer.Provider == "local" {
		if isRunningInDocker() {
			return nil, fmt.Errorf("Docker environment detected. LocalCryptoProvider is not supported in headless Docker environment. Please use 'rpc' provider with Companion app.")
		}
		localProvider := security.NewLocalCryptoProvider()

		// Try to load existing salt
		salt, err := store.GetSalt("default")
		if err == nil {
			// Salt exists - prompt for passphrase interactively
			if prompt.IsInteractive() {
				passphrase, err := prompt.Passphrase("Enter passphrase to unlock local key: ")
				if err != nil {
					return nil, fmt.Errorf("failed to read passphrase: %w", err)
				}

				// Verify checksum if exists
				storedChecksum, err := store.GetChecksum("default")
				if err == nil {
					// Checksum exists, verify it
					newChecksum := localProvider.CalculateChecksum(passphrase, salt)
					if !bytes.Equal(storedChecksum, newChecksum) {
						return nil, fmt.Errorf("incorrect passphrase")
					}
				} else {
					logger.Warn("no checksum found for local key, skipping verification")
				}

				if err := localProvider.InitializeWithSalt(passphrase, salt); err != nil {
					return nil, fmt.Errorf("failed to initialize local crypto provider with salt: %w", err)
				}
				cryptoProvider = localProvider
				app.CryptoProvider = cryptoProvider
				logger.Info("local crypto provider initialized")
			} else {
				// Non-interactive mode fallback (ENV var)
				passphrase := os.Getenv("LANGO_PASSPHRASE")
				if cfg.Security.Passphrase != "" {
					logger.Warn("security.passphrase in config is deprecated and ignored; use LANGO_PASSPHRASE env var")
				}

				if passphrase == "" {
					return nil, fmt.Errorf("local crypto enabled but no passphrase and not interactive (set LANGO_PASSPHRASE env var)")
				}

				// Verify checksum if exists
				storedChecksum, err := store.GetChecksum("default")
				if err == nil {
					newChecksum := localProvider.CalculateChecksum(passphrase, salt)
					if !bytes.Equal(storedChecksum, newChecksum) {
						return nil, fmt.Errorf("incorrect passphrase (from LANGO_PASSPHRASE)")
					}
				}

				if err := localProvider.InitializeWithSalt(passphrase, salt); err != nil {
					return nil, fmt.Errorf("failed to initialize local crypto provider with salt: %w", err)
				}
				cryptoProvider = localProvider
				app.CryptoProvider = cryptoProvider
				logger.Warn("local crypto provider initialized using env var (not recommended for production)")
			}
		} else {
			// First-time setup - generate new salt
			if prompt.IsInteractive() {
				passphrase, err := prompt.PassphraseConfirm("Enter new passphrase for local key: ", "Confirm passphrase: ")
				if err != nil {
					return nil, fmt.Errorf("passphrase setup failed: %w", err)
				}

				if err := localProvider.Initialize(passphrase); err != nil {
					return nil, fmt.Errorf("failed to initialize local crypto provider: %w", err)
				}

				// Store the salt
				if err := store.SetSalt("default", localProvider.Salt()); err != nil {
					logger.Warnw("failed to persist salt", "error", err)
				}

				// Store the checksum
				checksum := localProvider.CalculateChecksum(passphrase, localProvider.Salt())
				if err := store.SetChecksum("default", checksum); err != nil {
					logger.Warnw("failed to persist checksum", "error", err)
				}

				cryptoProvider = localProvider
				app.CryptoProvider = cryptoProvider
				logger.Info("local crypto provider initialized (new setup)")
			} else {
				// Non-interactive setup fallback
				passphrase := os.Getenv("LANGO_PASSPHRASE")
				if cfg.Security.Passphrase != "" {
					logger.Warn("security.passphrase in config is deprecated and ignored; use LANGO_PASSPHRASE env var")
				}

				if passphrase != "" {
					if err := localProvider.Initialize(passphrase); err != nil {
						logger.Warnw("failed to initialize local crypto provider", "error", err)
					} else {
						// Store the salt
						if err := store.SetSalt("default", localProvider.Salt()); err != nil {
							logger.Warnw("failed to persist salt", "error", err)
						}

						// Store the checksum
						checksum := localProvider.CalculateChecksum(passphrase, localProvider.Salt())
						if err := store.SetChecksum("default", checksum); err != nil {
							logger.Warnw("failed to persist checksum", "error", err)
						}

						cryptoProvider = localProvider
						app.CryptoProvider = cryptoProvider
						logger.Info("local crypto provider initialized (new)")
					}
				} else {
					logger.Warn("local crypto enabled but no passphrase, set LANGO_PASSPHRASE env var")
				}
			}
		}
	}

	// Register Security Tools (if crypto provider is available)
	if app.CryptoProvider != nil {
		keyRegistry := security.NewKeyRegistry(store.Client())
		secretsStore := security.NewSecretsStore(store.Client(), keyRegistry, app.CryptoProvider)

		// Register secrets tool
		secretsTool := secrets.New(secretsStore)
		if err := app.Agent.RegisterTool(&agent.Tool{
			Name:        "secrets",
			Description: "Securely store and retrieve secrets (API keys, tokens, etc.)",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"operation": map[string]interface{}{
						"type":        "string",
						"description": "Operation to perform: store, get, list, delete",
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Secret name",
					},
					"value": map[string]interface{}{
						"type":        "string",
						"description": "Secret value (for store operation)",
					},
				},
				"required": []string{"operation"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				op, _ := params["operation"].(string)
				switch op {
				case "store":
					return secretsTool.Store(ctx, params)
				case "get":
					return secretsTool.Get(ctx, params)
				case "list":
					return secretsTool.List(ctx, params)
				case "delete":
					return secretsTool.Delete(ctx, params)
				default:
					return nil, fmt.Errorf("unknown operation: %s", op)
				}
			},
		}); err != nil {
			logger.Warnw("failed to register secrets tool", "error", err)
		}

		// Register crypto tool
		cryptoTool := crypto.New(app.CryptoProvider, keyRegistry)
		if err := app.Agent.RegisterTool(&agent.Tool{
			Name:        "crypto",
			Description: "Cryptographic operations (encrypt, decrypt, sign, hash)",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"operation": map[string]interface{}{
						"type":        "string",
						"description": "Operation: encrypt, decrypt, sign, hash, keys",
					},
					"data": map[string]interface{}{
						"type":        "string",
						"description": "Data to process",
					},
					"ciphertext": map[string]interface{}{
						"type":        "string",
						"description": "Ciphertext to decrypt (base64)",
					},
					"keyId": map[string]interface{}{
						"type":        "string",
						"description": "Key ID to use",
					},
					"algorithm": map[string]interface{}{
						"type":        "string",
						"description": "Hash algorithm (sha256, sha512)",
					},
				},
				"required": []string{"operation"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				op, _ := params["operation"].(string)
				switch op {
				case "encrypt":
					return cryptoTool.Encrypt(ctx, params)
				case "decrypt":
					return cryptoTool.Decrypt(ctx, params)
				case "sign":
					return cryptoTool.Sign(ctx, params)
				case "hash":
					return cryptoTool.Hash(ctx, params)
				case "keys":
					return cryptoTool.Keys(ctx, params)
				default:
					return nil, fmt.Errorf("unknown operation: %s", op)
				}
			},
		}); err != nil {
			logger.Warnw("failed to register crypto tool", "error", err)
		}
	}

	// Initialize Channels
	if err := app.initChannels(); err != nil {
		logger.Errorw("failed to initialize channels", "error", err)
	}

	// Initialize AuthManager
	var authManager *gateway.AuthManager
	if len(cfg.Auth.Providers) > 0 {
		am, err := gateway.NewAuthManager(cfg.Auth, store)
		if err != nil {
			logger.Warnw("failed to initialize auth manager (OIDC discovery failed?)", "error", err)
		} else {
			authManager = am
		}
	}

	// Initialize Gateway (Server)
	app.Gateway = gateway.New(gateway.Config{
		Host:             cfg.Server.Host,
		Port:             cfg.Server.Port,
		HTTPEnabled:      cfg.Server.HTTPEnabled,
		WebSocketEnabled: cfg.Server.WebSocketEnabled,
	}, app.Agent, rpcProvider, app.Store, authManager)

	return app, nil
}

// Start starts the application services
func (a *App) Start(ctx context.Context) error {
	logger.Info("starting application")

	// Start Gateway
	go func() {
		if err := a.Gateway.Start(); err != nil {
			logger.Errorw("gateway server error", "error", err)
		}
	}()

	// Start Channels
	for _, ch := range a.Channels {
		go func(c Channel) {
			if err := c.Start(ctx); err != nil {
				logger.Errorw("channel start error", "error", err)
			}
		}(ch)
	}

	return nil
}

// Stop stops the application services
func (a *App) Stop(ctx context.Context) error {
	logger.Info("stopping application")

	// Stop Gateway
	if err := a.Gateway.Shutdown(ctx); err != nil {
		logger.Errorw("gateway shutdown error", "error", err)
	}

	// Stop Channels
	for _, ch := range a.Channels {
		ch.Stop()
	}

	// Close Session Store
	if a.Store != nil {
		if err := a.Store.Close(); err != nil {
			logger.Errorw("session store close error", "error", err)
		}
	}

	return nil
}

// isRunningInDocker checks if the application is running inside a Docker container
func isRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	content, err := os.ReadFile("/proc/1/cgroup")
	if err == nil && bytes.Contains(content, []byte("docker")) {
		return true
	}
	return false
}
