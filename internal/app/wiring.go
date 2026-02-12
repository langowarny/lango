package app

import (
	"context"
	"crypto/hmac"
	"fmt"
	"os"

	"github.com/langowarny/lango/internal/adk"
	"github.com/langowarny/lango/internal/agent"
	"github.com/langowarny/lango/internal/config"
	"github.com/langowarny/lango/internal/gateway"
	"github.com/langowarny/lango/internal/knowledge"
	"github.com/langowarny/lango/internal/learning"
	"github.com/langowarny/lango/internal/security"
	"github.com/langowarny/lango/internal/session"
	"github.com/langowarny/lango/internal/skill"
	"github.com/langowarny/lango/internal/supervisor"
	"google.golang.org/adk/model"
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

// initSecurity creates and initializes the security stack.
func initSecurity(cfg *config.Config, store session.Store) (security.CryptoProvider, *security.KeyRegistry, *security.SecretsStore, error) {
	if cfg.Security.Signer.Provider == "" {
		return nil, nil, nil, nil
	}

	switch cfg.Security.Signer.Provider {
	case "local":
		passphrase := os.Getenv("LANGO_PASSPHRASE")
		if passphrase == "" {
			passphrase = cfg.Security.Passphrase
		}
		if passphrase == "" {
			return nil, nil, nil, fmt.Errorf("local security provider requires a passphrase")
		}

		provider := security.NewLocalCryptoProvider()

		entStore, ok := store.(*session.EntStore)
		if !ok {
			return nil, nil, nil, fmt.Errorf("local security provider requires EntStore")
		}

		salt, err := entStore.GetSalt("default")
		if err != nil {
			// First-time setup: initialize with new salt
			if err := provider.Initialize(passphrase); err != nil {
				return nil, nil, nil, fmt.Errorf("initialize crypto provider: %w", err)
			}
			if err := entStore.SetSalt("default", provider.Salt()); err != nil {
				return nil, nil, nil, fmt.Errorf("store salt: %w", err)
			}
			checksum := provider.CalculateChecksum(passphrase, provider.Salt())
			if err := entStore.SetChecksum("default", checksum); err != nil {
				return nil, nil, nil, fmt.Errorf("store checksum: %w", err)
			}
		} else {
			// Existing salt found: initialize with it
			if err := provider.InitializeWithSalt(passphrase, salt); err != nil {
				return nil, nil, nil, fmt.Errorf("initialize crypto provider with salt: %w", err)
			}
			// Verify checksum
			storedChecksum, err := entStore.GetChecksum("default")
			if err == nil {
				computed := provider.CalculateChecksum(passphrase, salt)
				if !hmac.Equal(storedChecksum, computed) {
					return nil, nil, nil, fmt.Errorf("passphrase checksum mismatch: incorrect passphrase")
				}
			}
		}

		client := entStore.Client()
		keys := security.NewKeyRegistry(client)
		secrets := security.NewSecretsStore(client, keys, provider)

		// Register default encryption key
		ctx := context.Background()
		if _, err := keys.RegisterKey(ctx, "default", "local", security.KeyTypeEncryption); err != nil {
			return nil, nil, nil, fmt.Errorf("register default key: %w", err)
		}

		logger().Info("security initialized (local provider)")
		return provider, keys, secrets, nil

	case "rpc":
		provider := security.NewRPCProvider()

		entStore, ok := store.(*session.EntStore)
		if !ok {
			return nil, nil, nil, fmt.Errorf("rpc security provider requires EntStore")
		}

		client := entStore.Client()
		keys := security.NewKeyRegistry(client)
		secrets := security.NewSecretsStore(client, keys, provider)

		logger().Info("security initialized (rpc provider)")
		return provider, keys, secrets, nil

	case "enclave":
		return nil, nil, nil, fmt.Errorf("enclave provider not yet implemented")

	default:
		return nil, nil, nil, fmt.Errorf("unknown security provider: %s", cfg.Security.Signer.Provider)
	}
}

// knowledgeComponents holds optional self-learning components.
type knowledgeComponents struct {
	store    *knowledge.Store
	engine   *learning.Engine
	registry *skill.Registry
}

// initKnowledge creates the self-learning components if enabled.
func initKnowledge(cfg *config.Config, store session.Store, baseTools []*agent.Tool) *knowledgeComponents {
	if !cfg.Knowledge.Enabled {
		logger().Info("knowledge system disabled")
		return nil
	}

	entStore, ok := store.(*session.EntStore)
	if !ok {
		logger().Warn("knowledge system requires EntStore, skipping")
		return nil
	}

	client := entStore.Client()
	kLogger := logger()

	kStore := knowledge.NewStore(
		client, kLogger,
		cfg.Knowledge.MaxKnowledge,
		cfg.Knowledge.MaxLearnings,
	)

	engine := learning.NewEngine(kStore, kLogger)
	registry, err := skill.NewRegistry(kStore, baseTools, kLogger)
	if err != nil {
		logger().Warnw("skill registry init error, skipping knowledge system", "error", err)
		return nil
	}

	ctx := context.Background()
	if err := registry.LoadSkills(ctx); err != nil {
		logger().Warnw("load skills error", "error", err)
	}

	logger().Info("knowledge system initialized")
	return &knowledgeComponents{
		store:    kStore,
		engine:   engine,
		registry: registry,
	}
}

// initAgent creates the ADK agent with the given tools and provider proxy.
func initAgent(ctx context.Context, sv *supervisor.Supervisor, cfg *config.Config, store session.Store, tools []*agent.Tool, kc *knowledgeComponents) (*adk.Agent, error) {
	// Adapt tools to ADK format
	var adkTools []adk_tool.Tool
	for _, t := range tools {
		at, err := adk.AdaptTool(t)
		if err != nil {
			logger().Warnw("adapt tool error", "name", t.Name, "error", err)
			continue
		}
		adkTools = append(adkTools, at)
	}

	// Create provider proxy and model adapter
	proxy := supervisor.NewProviderProxy(sv, cfg.Agent.Provider, cfg.Agent.Model)
	modelAdapter := adk.NewModelAdapter(proxy)

	// If knowledge is enabled, wrap with context-aware adapter
	var llm model.LLM = modelAdapter
	if kc != nil {
		retriever := knowledge.NewContextRetriever(
			kc.store,
			cfg.Knowledge.MaxContextPerLayer,
			logger(),
		)
		llm = adk.NewContextAwareModelAdapter(modelAdapter, retriever, _defaultSystemPrompt, logger())
	}

	// If PII redaction is enabled, wrap with PII-redacting adapter
	if cfg.Security.Interceptor.Enabled && cfg.Security.Interceptor.RedactPII {
		redactor := agent.NewPIIRedactor(agent.PIIConfig{
			RedactEmail: true,
			RedactPhone: true,
			CustomRegex: cfg.Security.Interceptor.PIIRegexPatterns,
		})
		llm = adk.NewPIIRedactingModelAdapter(llm, redactor)
		logger().Info("PII redaction interceptor enabled")
	}

	logger().Info("initializing agent runtime (ADK)...")
	adkAgent, err := adk.NewAgent(ctx, adkTools, llm, _defaultSystemPrompt, store)
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
