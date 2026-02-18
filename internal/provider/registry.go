package provider

import (
	"strings"
	"sync"
)

// Registry manages the registration and lookup of providers.
type Registry struct {
	providers map[string]Provider
	mu        sync.RWMutex
}

// NewRegistry creates a new provider registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry.
func (r *Registry) Register(p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[p.ID()] = p
}

// Get returns a provider by ID. It handles aliases (e.g., "gpt" -> "openai").
func (r *Registry) Get(id string) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	normalized := normalizeID(id)
	p, ok := r.providers[normalized]
	return p, ok
}

// List returns a list of all registered providers.
func (r *Registry) List() []Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var providers []Provider
	for _, p := range r.providers {
		providers = append(providers, p)
	}
	return providers
}

// normalizeID normalizes the provider ID and resolves aliases.
func normalizeID(id string) string {
	lower := strings.ToLower(strings.TrimSpace(id))
	switch lower {
	case "gpt", "chatgpt":
		return "openai"
	case "claude":
		return "anthropic"
	case "llama":
		return "ollama"
	case "bard":
		return "gemini"
	default:
		return lower
	}
}

// GetSupportedProviders returns a list of all supported provider IDs.
func GetSupportedProviders() []string {
	return []string{
		"gemini",
		"openai",
		"anthropic",
		"ollama",
	}
}
