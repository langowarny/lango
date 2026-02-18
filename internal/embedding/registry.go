package embedding

import (
	"fmt"

	"go.uber.org/zap"
)

// ProviderConfig holds configuration for creating an embedding provider.
type ProviderConfig struct {
	// Provider selects the embedding backend: "openai", "google", or "local".
	Provider string
	// Model is the embedding model identifier.
	Model string
	// Dimensions is the embedding vector dimensionality.
	Dimensions int
	// APIKey is required for openai and google providers.
	APIKey string
	// BaseURL is used by the local provider (Ollama endpoint).
	BaseURL string
}

// Registry manages embedding provider selection with fallback.
type Registry struct {
	primary  EmbeddingProvider
	fallback EmbeddingProvider
	logger   *zap.SugaredLogger
}

// NewRegistry creates a provider registry with the given primary config.
// If fallback configs are provided, the first successful one becomes the fallback.
func NewRegistry(primary ProviderConfig, fallbacks []ProviderConfig, logger *zap.SugaredLogger) (*Registry, error) {
	p, err := createProvider(primary)
	if err != nil {
		return nil, fmt.Errorf("primary provider (%s): %w", primary.Provider, err)
	}

	r := &Registry{
		primary: p,
		logger:  logger,
	}

	for _, fc := range fallbacks {
		fb, err := createProvider(fc)
		if err != nil {
			logger.Warnw("fallback provider init failed, skipping", "provider", fc.Provider, "error", err)
			continue
		}
		r.fallback = fb
		break
	}

	return r, nil
}

// Provider returns the active embedding provider.
func (r *Registry) Provider() EmbeddingProvider {
	return r.primary
}

// Fallback returns the fallback provider, or nil if none.
func (r *Registry) Fallback() EmbeddingProvider {
	return r.fallback
}

func createProvider(cfg ProviderConfig) (EmbeddingProvider, error) {
	switch cfg.Provider {
	case "openai":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("openai provider requires API key")
		}
		return NewOpenAIProvider(cfg.APIKey, cfg.Model, cfg.Dimensions), nil

	case "google":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("google provider requires API key")
		}
		return NewGoogleProvider(cfg.APIKey, cfg.Model, cfg.Dimensions)

	case "local":
		return NewLocalProvider(cfg.BaseURL, cfg.Model, cfg.Dimensions), nil

	default:
		return nil, fmt.Errorf("unknown embedding provider: %s", cfg.Provider)
	}
}
