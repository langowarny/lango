## Why

The Context Retriever defines 8 context layers but only 4 are implemented. Tool Registry and Runtime Context layers are skipped with `continue // handled elsewhere` but are never actually handled anywhere. This means the LLM cannot see available tools or current session state in its prompt context, limiting its situational awareness.

## What Changes

- Add `ToolDescriptor` and `RuntimeContext` domain types to the knowledge package
- Define `ToolRegistryProvider` and `RuntimeContextProvider` interfaces on the retriever
- Implement `retrieveTools()` and `retrieveRuntimeContext()` methods, replacing the `continue` skip
- Add `ToolRegistryAdapter` (adapts `[]*agent.Tool`) and `RuntimeContextAdapter` (session/system state with mutex-protected updates) in the adk package
- Update `AssemblePrompt()` to emit "Runtime Context" and "Available Tools" sections before existing sections
- Update `ContextAwareModelAdapter` to request all 6 layers and set session state before retrieval
- Wire adapters in `initAgent()` when knowledge system is enabled

## Capabilities

### New Capabilities

### Modified Capabilities
- `context-retriever`: Tool Registry and Runtime Context layers are now fully implemented with retrieval, prompt assembly, and adapter wiring

## Impact

- `internal/knowledge/types.go`: New types added
- `internal/knowledge/retriever.go`: New interfaces, builder methods, retrieval methods, updated switch and prompt assembly
- `internal/adk/context_providers.go`: New file with adapter implementations
- `internal/adk/context_model.go`: Extended layer requests, runtime adapter integration
- `internal/app/wiring.go`: Adapter creation and wiring
- No breaking changes — default layer list (nil) remains the original 4 layers
