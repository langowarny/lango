package settings

import (
	"context"
	"sort"
	"time"

	"github.com/langoai/lango/internal/config"
	"github.com/langoai/lango/internal/provider"
	provanthropic "github.com/langoai/lango/internal/provider/anthropic"
	provgemini "github.com/langoai/lango/internal/provider/gemini"
	provopenai "github.com/langoai/lango/internal/provider/openai"
	"github.com/langoai/lango/internal/types"
)

const modelFetchTimeout = 5 * time.Second

// newProviderFromConfig creates a lightweight provider instance from config.
// Returns nil if the provider cannot be created (missing API key, unknown type, etc.).
func newProviderFromConfig(id string, pCfg config.ProviderConfig) provider.Provider {
	apiKey := pCfg.APIKey
	if apiKey == "" && pCfg.Type != types.ProviderOllama {
		return nil
	}

	switch pCfg.Type {
	case types.ProviderOpenAI:
		return provopenai.NewProvider(id, apiKey, pCfg.BaseURL)
	case types.ProviderAnthropic:
		return provanthropic.NewProvider(id, apiKey)
	case types.ProviderGemini, types.ProviderGoogle:
		p, err := provgemini.NewProvider(context.Background(), id, apiKey, "")
		if err != nil {
			return nil
		}
		return p
	case types.ProviderOllama:
		baseURL := pCfg.BaseURL
		if baseURL == "" {
			baseURL = "http://localhost:11434/v1"
		}
		return provopenai.NewProvider(id, apiKey, baseURL)
	case types.ProviderGitHub:
		baseURL := pCfg.BaseURL
		if baseURL == "" {
			baseURL = "https://models.inference.ai.azure.com"
		}
		return provopenai.NewProvider(id, apiKey, baseURL)
	default:
		return nil
	}
}

// fetchModelOptions fetches available models from a provider.
// Returns a sorted list of model IDs, or nil if fetching fails.
// The currentModel (if non-empty) is always included in the result.
func fetchModelOptions(providerID string, cfg *config.Config, currentModel string) []string {
	pCfg, ok := cfg.Providers[providerID]
	if !ok {
		return nil
	}

	p := newProviderFromConfig(providerID, pCfg)
	if p == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), modelFetchTimeout)
	defer cancel()

	models, err := p.ListModels(ctx)
	if err != nil || len(models) == 0 {
		return nil
	}

	seen := make(map[string]bool, len(models))
	opts := make([]string, 0, len(models))
	for _, m := range models {
		if !seen[m.ID] {
			seen[m.ID] = true
			opts = append(opts, m.ID)
		}
	}
	sort.Strings(opts)

	// Ensure current model is included
	if currentModel != "" && !seen[currentModel] {
		opts = append([]string{currentModel}, opts...)
	}

	return opts
}
