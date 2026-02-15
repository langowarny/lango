package checks

import (
	"context"
	"fmt"

	"github.com/langowarny/lango/internal/config"
)

// EmbeddingCheck validates embedding/RAG configuration.
type EmbeddingCheck struct{}

// Name returns the check name.
func (c *EmbeddingCheck) Name() string {
	return "Embedding / RAG"
}

// Run checks embedding configuration validity.
func (c *EmbeddingCheck) Run(_ context.Context, cfg *config.Config) Result {
	if cfg == nil {
		return Result{Name: c.Name(), Status: StatusSkip, Message: "Configuration not loaded"}
	}

	emb := cfg.Embedding
	if emb.Provider == "" {
		return Result{
			Name:    c.Name(),
			Status:  StatusSkip,
			Message: "Embedding provider not configured",
		}
	}

	var issues []string
	status := StatusPass

	// Validate provider type.
	switch emb.Provider {
	case "openai", "google", "local":
		// ok
	default:
		issues = append(issues, fmt.Sprintf("unknown provider: %s", emb.Provider))
		status = StatusFail
	}

	// Check API key availability for cloud providers.
	if emb.Provider == "openai" {
		if p, ok := cfg.Providers["openai"]; !ok || p.APIKey == "" {
			issues = append(issues, "openai embedding requires an API key in providers.openai")
			status = StatusFail
		}
	}
	if emb.Provider == "google" {
		hasKey := false
		if p, ok := cfg.Providers["google"]; ok && p.APIKey != "" {
			hasKey = true
		}
		if p, ok := cfg.Providers["gemini"]; ok && p.APIKey != "" {
			hasKey = true
		}
		if !hasKey {
			issues = append(issues, "google embedding requires an API key in providers.google or providers.gemini")
			status = StatusFail
		}
	}

	// Check dimensions.
	if emb.Dimensions <= 0 {
		issues = append(issues, "dimensions should be set (e.g. 1536 for OpenAI, 768 for Google/local)")
		if status < StatusWarn {
			status = StatusWarn
		}
	}

	// Check RAG settings.
	if emb.RAG.Enabled && emb.RAG.MaxResults <= 0 {
		issues = append(issues, "rag.maxResults should be positive when RAG is enabled")
		if status < StatusWarn {
			status = StatusWarn
		}
	}

	if len(issues) == 0 {
		msg := fmt.Sprintf("Embedding configured (provider=%s, dimensions=%d, rag=%v)",
			emb.Provider, emb.Dimensions, emb.RAG.Enabled)
		return Result{Name: c.Name(), Status: StatusPass, Message: msg}
	}

	message := "Embedding issues:\n"
	for _, issue := range issues {
		message += fmt.Sprintf("- %s\n", issue)
	}
	return Result{Name: c.Name(), Status: status, Message: message}
}

// Fix delegates to Run as automatic fixing is not supported.
func (c *EmbeddingCheck) Fix(ctx context.Context, cfg *config.Config) Result {
	return c.Run(ctx, cfg)
}
