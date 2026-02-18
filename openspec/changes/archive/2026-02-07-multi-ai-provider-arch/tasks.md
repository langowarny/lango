## 1. Provider Core & Registry
- [x] 1.1 Define `Provider` interface and common types in `internal/provider/provider.go`
- [x] 1.2 Implement `Registry` with thread-safe registration and ID normalization
- [x] 1.3 Add unit tests for `Registry` behavior

## 2. OpenAI-Compatible Provider
- [x] 2.1 Add `github.com/openai/openai-go/v3` dependency
- [x] 2.2 Implement `OpenAICompatibleProvider` struct in `internal/provider/openai/`
- [x] 2.3 Implement `Generate` method with streaming support using iterators
- [x] 2.4 Implement tool calling parameter conversion
- [x] 2.5 Add integration test with Ollama (if available) or mock

## 3. Anthropic Provider
- [x] 3.1 Add `github.com/anthropics/anthropic-sdk-go` dependency
- [x] 3.2 Implement `AnthropicProvider` struct in `internal/provider/anthropic/`
- [x] 3.3 Implement `Generate` method with streaming support
- [x] 3.4 Implement tool calling conversion for Anthropic format
- [x] 3.5 Verify compilation and basic instantiation

## 4. Config & Runtime Integration
- [x] 4.1 Update `internal/config` to support `providers` section
- [x] 4.2 Initialize providers registry in `internal/app/app.go`
- [x] 4.3 Refactor `internal/agent/runtime.go` to use `Provider` interface
- [x] 4.4 Implement `GeminiProvider` wrapper for backward compatibility
- [x] 4.5 Update `lango.json` in local environment for testing

## 5. Model Fallback & Finalization
- [x] 5.1 Implement fallback logic in `agent.Runtime`
- [x] 5.2 Add error classification for retryable errors
- [x] 5.3 Verify fallback behavior with simulated failures
- [x] 5.4 Update documentation and archive change
