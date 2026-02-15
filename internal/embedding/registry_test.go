package embedding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewRegistry_Local(t *testing.T) {
	logger := zap.NewNop().Sugar()

	reg, err := NewRegistry(ProviderConfig{
		Provider:   "local",
		Model:      "test-model",
		Dimensions: 128,
		BaseURL:    "http://localhost:11434/v1",
	}, nil, logger)

	require.NoError(t, err)
	assert.Equal(t, "local", reg.Provider().ID())
	assert.Equal(t, 128, reg.Provider().Dimensions())
	assert.Nil(t, reg.Fallback())
}

func TestNewRegistry_UnknownProvider(t *testing.T) {
	logger := zap.NewNop().Sugar()

	_, err := NewRegistry(ProviderConfig{
		Provider: "unknown",
	}, nil, logger)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown embedding provider")
}

func TestNewRegistry_OpenAI_NoKey(t *testing.T) {
	logger := zap.NewNop().Sugar()

	_, err := NewRegistry(ProviderConfig{
		Provider: "openai",
	}, nil, logger)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "API key")
}

func TestNewRegistry_WithFallback(t *testing.T) {
	logger := zap.NewNop().Sugar()

	reg, err := NewRegistry(
		ProviderConfig{Provider: "local", Dimensions: 128},
		[]ProviderConfig{
			{Provider: "local", Dimensions: 256, BaseURL: "http://other:11434/v1"},
		},
		logger,
	)

	require.NoError(t, err)
	assert.Equal(t, 128, reg.Provider().Dimensions())
	require.NotNil(t, reg.Fallback())
	assert.Equal(t, 256, reg.Fallback().Dimensions())
}
