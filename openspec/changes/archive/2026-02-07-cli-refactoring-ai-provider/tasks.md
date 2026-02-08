# Tasks: CLI Multi-Provider Refactoring

## 1. Provider Management Core

- [x] 1.1 Implement `GetSupportedProviders()` in `internal/provider` (or ensuring `Registry.List()` suffices with metadata).
- [x] 1.2 Create `internal/cli/common` package (or similar) to hold shared CLI provider logic (list providers, get API key env var name).

## 2. Update Onboard Command

- [x] 2.1 Refactor `onboard/wizard.go` state machine to include a "Provider Selection" step.
- [x] 2.2 Implement provider selection UI using `bubbletea` list.
- [x] 2.3 Implement dynamic API key prompt based on selected provider.
- [x] 2.4 Update `SaveConfig` in `onboard/config.go` to generate `providers` map structure + `agent.provider` reference.

## 3. Update Doctor Command

- [x] 3.1 Refactor `doctor/checks` to remove hardcoded `APIKeyCheck` for `GOOGLE_API_KEY`.
- [x] 3.2 Implement `ProvidersCheck` in `doctor/checks/providers.go` that iterates `config.Providers`.
- [x] 3.3 Verify credentials for each configured provider (OpenAI, Anthropic, Gemini).

## 4. Verification & Cleanup

- [x] 4.1 Verify `onboard` flow with all 3 providers.
- [x] 4.2 Verify `doctor` check with valid/invalid keys for all 3 providers.
- [x] 4.3 Ensure no regressions in `lango serve` startup with the new config structure.
