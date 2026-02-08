package supervisor

import (
	"context"
	"iter"

	"github.com/langowarny/lango/internal/provider"
)

// ProviderProxy implements provider.Provider but forwards requests to the Supervisor.
type ProviderProxy struct {
	supervisor   *Supervisor
	providerID   string
	defaultModel string
}

// NewProviderProxy creates a new proxy for a specific provider.
func NewProviderProxy(sv *Supervisor, providerID, defaultModel string) *ProviderProxy {
	return &ProviderProxy{
		supervisor:   sv,
		providerID:   providerID,
		defaultModel: defaultModel,
	}
}

// ID returns the provider ID.
func (p *ProviderProxy) ID() string {
	return p.providerID
}

// Generate forwards the request to the Supervisor.
func (p *ProviderProxy) Generate(ctx context.Context, params provider.GenerateParams) (iter.Seq2[provider.StreamEvent, error], error) {
	return p.supervisor.Generate(ctx, p.providerID, p.defaultModel, params)
}

// ListModels is not yet fully implemented via proxy, returning empty for now or could similarly proxy.
func (p *ProviderProxy) ListModels(ctx context.Context) ([]provider.ModelInfo, error) {
	// TODO: Implement ListModels proxying if needed
	return []provider.ModelInfo{}, nil
}
