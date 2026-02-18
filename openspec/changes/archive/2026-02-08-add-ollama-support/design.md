## Context

The Lango CLI currently supports Gemini, OpenAI, and Anthropic. Users have requested support for local LLMs via Ollama. The codebase already has an `OpenAIProvider` implementation which is compatible with Ollama's API. This design outlines how to wire Ollama into the existing system with minimal friction.

## Goals / Non-Goals

**Goals:**
- Enable "ollama" as a first-class provider in `lango.json`.
- Support standard chat completion via Ollama using the existing OpenAI compatibility layer.
- Add Ollama to the `onboard` wizard with a streamlined setup (skip API key).

**Non-Goals:**
- Auto-discovery of installed Ollama models (users must specify model name manually for now, or use a default).
- Managing the Ollama process (starting/stopping).

## Decisions

### 1. Reuse OpenAI Provider
We will reuse the existing `internal/provider/openai` package. Ollama provides an OpenAI-compatible API endpoint `/v1`.
- **Decision**: In `supervisor.go`, initializing a provider with type "ollama" will instantiate an `OpenAIProvider` with `BaseURL` set to `http://localhost:11434/v1`.
- **Rationale**: Avoids code duplication. `ollama` is just a configuration variance of `openai`.

### 2. Explicit Provider Type "ollama"
We will introduce a distinct provider type "ollama" in the configuration and registry, rather than asking users to configure "openai" with a custom base URL.
- **Decision**: `lango.json` will have a specific `ollama` section.
- **Rationale**: Better UX. Users think of "Ollama" as a distinct tool, not just "OpenAI with a different URL".

### 3. Streamlined Onboarding
The `onboard` wizard will be updated to handle "ollama" specially.
- **Decision**: When "Ollama" is selected, the API Key step will be skipped or auto-filled with a placeholder.
- **Rationale**: Ollama does not require an API key by default. Asking for one is confusing.

## Risks / Trade-offs

- **Risk**: Hardcoded default URL (`http://localhost:11434/v1`) might not work for all setups (e.g., Docker, remote Ollama).
    - **Mitigation**: Users can override `baseUrl` in `lango.json`.
- **Risk**: Default model might not be pulled in Ollama.
    - **Mitigation**: We will default to a common model like `llama3`, but if it's missing, the error from Ollama will be propagated. We rely on the user to `ollama pull` the model.

## Migration Plan

No migration needed for existing users. New users can select Ollama during onboarding. Existing users can manually add the `ollama` block to `lango.json`.
