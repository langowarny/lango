## Why

The `agent-runtime` spec requires ADK-Go framework integration, but the current implementation uses a placeholder runtime without actual LLM provider integration. This change implements the spec as designed, replacing the placeholder with Google's ADK (`google.golang.org/adk`) to provide real LLM agent capabilities with tool calling, streaming, and multi-model support.

## What Changes

- **Replace placeholder agent runtime**: Integrate ADK framework to replace the current placeholder implementation in `internal/agent/runtime.go`
- **Add ADK dependency**: Add `google.golang.org/adk` to go.mod
- **Implement with LlmAgent**: Use ADK's `llmagent.New()` to create LLM-driven agents with model configuration
- **Support multiple models**: Enable Gemini, Claude (via Vertex AI), and custom model providers through ADK's model interface
- **Wire up tool system**: Connect existing tool infrastructure (`Tool`, `ToolHandler`) to ADK's `tool.Tool` interface
- **Enable streaming responses**: Implement proper streaming using ADK's agent execution and response handling
- **Maintain session integration**: Ensure ADK agents work with existing session store

## Capabilities

### New Capabilities

None (implementing existing capability)

### Modified Capabilities

- `agent-runtime`: Update implementation to use ADK framework instead of placeholder. The spec requirements remain unchanged - this is purely an implementation change to fulfill the existing spec.

## Impact

- **Dependencies**: Adds ADK (`google.golang.org/adk@v0.4.0`) and Gemini SDK (`google.golang.org/genai`)
- **Code changes**: Major refactor of `internal/agent/runtime.go` to use ADK's llmagent APIs
- **Breaking changes**: None - external interfaces remain the same
- **Build**: No additional build requirements beyond existing Go 1.22+ and CGO
- **Testing**: Need integration tests with real LLM providers (requires API keys)
- **Model support**: Initial focus on Gemini (built-in), with extensibility for other providers via ADK's model interface
