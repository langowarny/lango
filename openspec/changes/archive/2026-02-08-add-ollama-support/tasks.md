## 1. Registry & Config

- [x] 1.1 Add "ollama" to `internal/provider/registry.go` in `GetSupportedProviders` and `normalizeID`.
- [x] 1.2 Add "ollama" metadata to `internal/cli/common/providers.go` in `GetProviderMetadata`.

## 2. Supervisor Integration

- [x] 2.1 Update `internal/supervisor/supervisor.go` in `initializeProviders` to handle `type: "ollama"`. It should instantiate `openai.NewProvider` with the default base URL `http://localhost:11434/v1` (unless overridden by config).

## 3. Onboard Wizard

- [x] 3.1 Update `internal/cli/onboard/wizard.go` method `getProviders` to include "ollama" (or ensure it picks it up from common).
- [x] 3.2 Update `internal/cli/onboard/wizard.go` method `viewAPIKey` to show a "No API key needed" message for Ollama.
- [x] 3.3 Update `internal/cli/onboard/wizard.go` method `handleEnter` (StepAPIKey) to skip input validation for Ollama.
- [x] 3.4 Update `internal/cli/onboard/config.go` method `SaveConfig` to write the correct `ollama` provider configuration (no apiKey env var).

## 4. Verification

- [x] 4.1 Manual Verification: Run `go run ./cmd/lango onboard` and select "Ollama". Verify `lango.json` is created with correct `ollama` provider block and `http://localhost:11434/v1` base URL.
