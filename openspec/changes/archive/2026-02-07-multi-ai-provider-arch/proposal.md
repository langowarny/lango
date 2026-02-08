## Why

Lango currently supports only Gemini as an LLM provider. Users need flexibility to choose different AI providers (OpenAI, Anthropic, Ollama for local models, etc.) based on their needs, preferences, or existing API subscriptions. Multi-provider support also enables model fallback for resilience and cost optimization.

## What Changes

- Add unified `Provider` interface with `ID()`, `Generate()`, and `ListModels()` methods
- Implement `ProviderRegistry` for dynamic provider registration and lookup
- Add OpenAI-compatible provider supporting 10+ services (OpenAI, Ollama, Groq, Together AI, etc.)
- Add Anthropic provider for Claude models with streaming and tool support
- Extend `lango.json` with `providers` section for multi-provider configuration
- Add model fallback system for automatic failover between providers

## Capabilities

### New Capabilities
- `provider-interface`: Unified abstraction for all LLM providers with streaming support
- `provider-openai-compatible`: OpenAI API compatible provider supporting multiple backends
- `provider-anthropic`: Native Anthropic/Claude provider implementation
- `provider-registry`: Dynamic provider registration, lookup, and lifecycle management

### Modified Capabilities
- `config-system`: Extended with `providers` configuration section and environment variable support
- `agent-runtime`: Refactored to use Provider abstraction instead of direct Gemini calls

## Impact

- **Code**: `internal/agent/runtime.go` refactored to use Provider interface
- **New packages**: `internal/provider/` with implementations for each provider
- **Config**: `lango.json` schema extended with providers section
- **Dependencies**: Add `openai-go` and `anthropic-sdk-go` modules
- **Breaking**: None - existing Gemini-only configs remain backwards compatible
