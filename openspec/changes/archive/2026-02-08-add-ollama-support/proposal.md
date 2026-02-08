## Why

Users want to run LLMs locally using Ollama for privacy, cost savings, and offline capabilities. Although Ollama support was mentioned in previous design documents, the implementation is currently missing from the codebase. The `onboard` wizard does not offer Ollama as an option, and the provider registry does not recognize it, preventing users from easily configuring local models.

## What Changes

1.  **Provider Registry**: Add `ollama` to the list of supported providers in `internal/provider` and `internal/cli/common`.
2.  **Supervisor Initialization**: Update `internal/supervisor` to initialize `ollama` using the existing OpenAI-compatible provider implementation, configured with the default Ollama base URL (`http://localhost:11434/v1`).
3.  **Onboarding Wizard**: Update the `lango onboard` CLI command to include "Ollama" in the provider selection list and skip the API key step (as it's not required for local Ollama).

## Capabilities

### New Capabilities
- `provider-ollama`: Support for local Ollama models via the OpenAI-compatible API interface.

### Modified Capabilities
<!-- No existing capability specs are being modified in terms of requirements. -->

## Impact

- `internal/provider/registry.go`
- `internal/cli/common/providers.go`
- `internal/supervisor/supervisor.go`
- `internal/cli/onboard/wizard.go` & `config.go`
