package app

import (
	"context"
	"fmt"
	"os"

	"github.com/langowarny/lango/internal/adk"
	"github.com/langowarny/lango/internal/agent"
	"github.com/langowarny/lango/internal/config"
	"github.com/langowarny/lango/internal/gateway"
	"github.com/langowarny/lango/internal/session"
	"github.com/langowarny/lango/internal/supervisor"
	adk_tool "google.golang.org/adk/tool"
)

const _defaultSystemPrompt = "You are Lango, a powerful AI assistant. You have access to tools for shell command execution and file system operations. Use them when appropriate to help the user."

// initSupervisor creates and initializes the Supervisor.
func initSupervisor(cfg *config.Config) (*supervisor.Supervisor, error) {
	logger().Info("initializing supervisor...")
	sv, err := supervisor.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("supervisor: %w", err)
	}
	return sv, nil
}

// initSessionStore creates and initializes the session store.
func initSessionStore(cfg *config.Config) (session.Store, error) {
	passphrase := os.Getenv("LANGO_PASSPHRASE")
	if passphrase == "" {
		passphrase = cfg.Security.Passphrase
	}

	var storeOpts []session.StoreOption
	if passphrase != "" {
		storeOpts = append(storeOpts, session.WithPassphrase(passphrase))
	}

	logger().Info("initializing session store...")
	store, err := session.NewEntStore(cfg.Session.DatabasePath, storeOpts...)
	if err != nil {
		return nil, fmt.Errorf("session store: %w", err)
	}
	return store, nil
}

// initAgent creates the ADK agent with the given tools and provider proxy.
func initAgent(ctx context.Context, sv *supervisor.Supervisor, cfg *config.Config, store session.Store, tools []*agent.Tool) (*adk.Agent, error) {
	// Adapt tools to ADK format
	var adkTools []adk_tool.Tool
	for _, t := range tools {
		at, err := adk.AdaptTool(t)
		if err != nil {
			logger().Warnw("failed to adapt tool", "name", t.Name, "error", err)
			continue
		}
		adkTools = append(adkTools, at)
	}

	// Create provider proxy and model adapter
	proxy := supervisor.NewProviderProxy(sv, cfg.Agent.Provider, cfg.Agent.Model)
	modelAdapter := adk.NewModelAdapter(proxy)

	logger().Info("initializing agent runtime (ADK)...")
	adkAgent, err := adk.NewAgent(ctx, adkTools, modelAdapter, _defaultSystemPrompt, store)
	if err != nil {
		return nil, fmt.Errorf("adk agent: %w", err)
	}
	return adkAgent, nil
}

// initGateway creates the gateway server.
func initGateway(cfg *config.Config, adkAgent *adk.Agent, store session.Store) *gateway.Server {
	return gateway.New(gateway.Config{
		Host:             cfg.Server.Host,
		Port:             cfg.Server.Port,
		HTTPEnabled:      cfg.Server.HTTPEnabled,
		WebSocketEnabled: cfg.Server.WebSocketEnabled,
	}, adkAgent, nil, store, nil)
}
