
## Why

The current custom agent runtime in `internal/agent` has proven brittle, particularly with the Gemini provider (e.g., protocol errors with empty text parts, incorrect role mapping for tools). Maintaining a custom runtime requires continuous updates to match evolving LLM API protocols and best practices.

Migrating to the official **Google Cloud Agent Development Kit (ADK) for Go** will provide a robust, production-ready foundation. ADK handles protocol nuances, tool execution loops, context management, and provider abstractions out of the box, allowing us to focus on building capabilities rather than debugging plumbing.

## What Changes

- **Runtime Replacement**: Replace the custom `internal/agent` package with an implementation based on `google/adk-go`.
- **Tool Adaptation**: Refactor or wrap existing tools (`filesystem`, `browser`, `secrets`, etc.) to implement ADK's `tools.Tool` interface.
- **Provider Integration**: Update `internal/provider` to leverage ADK's model abstractions or implement ADK-compatible model interfaces for our supported providers (Gemini, OpenAI, Anthropic).
- **State Management**: Adapt the ADK's state/session management to work with our existing `internal/session` (Ent-based) store.
- **Supervisor Update**: Update `internal/supervisor` to initialize and manage the ADK-based agent.

## Capabilities

### New Capabilities
- `adk-architecture`: Defines the core architecture for the ADK integration, including how the ADK agent is bootstrapped, configured, and exposed to the application.

### Modified Capabilities
- `agent-runtime`: Updates requirements for the agent runtime to align with ADK patterns (e.g., middleware configuration, tool execution flow, error handling).
- `provider-registry`: modifications to how providers are registered and resolved to support ADK's model interface.

## Impact

- **Codebase**: High impact on `internal/agent`, `internal/provider`, and `internal/supervisor`.
- **Dependencies**: Adds `github.com/google/adk-go`.
- **Configuration**: Potential changes to `lango.json` to support ADK-specific settings (though we will aim to maintain backward compatibility).
- **Behavior**: improved reliability and correctness of tool execution and multi-turn conversations.
