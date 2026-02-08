# Design: CLI Multi-Provider Refactoring

## Context

The `lango` CLI tools (`onboard`, `doctor`) currently assume a single-provider world (Gemini), hardcoding API key environment variables (`GOOGLE_API_KEY`) and configuration paths (`agent.apiKey`).

The core application (`internal/app`) has already migrated to a multi-provider architecture using a `Registry` pattern and a `providers` map in the configuration. This discrepancy causes friction for users wanting to use OpenAI or Anthropic and represents significant technical debt.

## Goals / Non-Goals

**Goals:**
1.  **Multi-Provider Onboarding**: Update `lango onboard` to support selecting and configuring OpenAI, Anthropic, and Gemini.
2.  **Comprehensive Health Checks**: Update `lango doctor` to verify credentials for all configured providers, not just Gemini.
3.  **Unified Configuration**: Ensure CLI tools read from and write to the modern `providers` configuration structure, aligning with `internal/app`.
4.  **Backward Compatibility**: `doctor` should still recognize legacy configurations but potentially warn or encourage migration.

**Non-Goals:**
1.  **Refactoring existing functionality**: We are not rewriting the entire CLI or `internal/app` logic beyond what is necessary for provider integration.
2.  **Adding new providers**: We are only exposing the providers *already supported* by the core (OpenAI, Anthropic, Gemini).

## Decisions

### 1. Onboarding Flow Update
The `onboard` wizard will be updated to:
- **Step 1: Select Provider**: Display a list (Gemini, OpenAI, Anthropic).
- **Step 2: API Key**: Prompt for the specific environment variable for the selected provider (e.g., `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`).
- **Step 3: Save Config**: Write to the `providers` map in `lango.json` instead of the legacy `agent` fields.
  ```json
  "providers": {
    "openai": {
      "type": "openai",
      "apiKey": "${OPENAI_API_KEY}"
    }
  },
  "agent": {
    "provider": "openai",
    "model": "gpt-4o"
  }
  ```

### 2. Doctor Checks Update
The `doctor` command will:
- Iterate through all entries in `config.Providers`.
- For each provider, validate that the referenced environment variable is set or the key is present.
- Remove the hardcoded `APIKeyCheck` that looks specifically for `GOOGLE_API_KEY`.
- Add a new `ProvidersCheck` that reports the status of all configured providers.

### 3. Config Structure
We will prioritize the `providers` map. The `agent` section will still be used to define the *default* provider and model, but the credentials and provider-specific settings will live in `providers`.

## Risks / Trade-offs

- **User Confusion**: Existing users with legacy configs might see warnings in `doctor`. We explicitly accept this to encourage migration to the more flexible system.
- **Env Var proliferation**: Users might need to manage multiple API keys if they configure multiple providers. This is acceptable as it matches the flexibility of the system.
