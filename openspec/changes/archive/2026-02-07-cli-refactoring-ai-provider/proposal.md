# Proposal: CLI Multi-Provider Refactoring

## Problem Description

The `lango` core application (`internal/app`) and provider system (`internal/provider`) have evolved to support multiple AI providers (OpenAI, Anthropic, Gemini) via a registry pattern and flexible configuration.

However, the CLI tools—specifically `lango onboard` and `lango doctor`—are stuck in the past. They:
1.  **Hardcode Gemini**: `onboard` only offers Gemini models.
2.  **Use Legacy Config**: `onboard` writes to and `doctor` reads from the deprecated `agent.provider` and `agent.apiKey` fields, ignoring the new `providers` map.
3.  **Ignore Other Providers**: There is no way to configure or verify OpenAI or Anthropic setups via the CLI.

This "generation gap" confuses users who want to use non-Gemini providers and creates technical debt where the CLI and App use different configuration sources.

## Proposed Solution

Refactor the CLI tools to align with the modern multi-provider architecture.

1.  **Update `onboard`**:
    *   Allow users to select their preferred provider (OpenAI, Anthropic, Gemini).
    *   Prompt for the appropriate API key based on the selection.
    *   Generate a `lango.json` that uses the `providers` map configuration structure.

2.  **Update `doctor`**:
    *   Verify credentials for *all* configured providers in the `providers` map.
    *   Deprecate checks for the legacy `GOOGLE_API_KEY` environment variable in favor of provider-specific checks.

3.  **Configuration Migration**:
    *   Ensure the generated config is compatible with the current `internal/app` logic.

## What Changes

### Capabilities

#### New Capabilities
- `cli-provider-management`: CLI support for configuring and verifying multiple AI providers.

### Modified Capabilities
- `config-system`: Update configuration validation to prioritize the `providers` map over legacy fields.

## Impact

- **CLI**: `internal/cli/onboard`, `internal/cli/doctor`
- **Config**: `internal/config` (potentially updating validation logic)
- **User Experience**: Users will see new options in the onboarding wizard and more comprehensive health checks.
