# Delta Spec: CLI Onboard - Multi-Provider

## Overview

Updates the onboarding wizard to support selecting from multiple AI providers.

## ADDED Requirements

### Requirement: Provider Selection
The wizard SHALL allow the user to select their desired AI provider from a list.

#### Scenario: Selecting a provider
- **WHEN** the user reaches the "Provider Selection" step
- **THEN** a list of supported providers (Gemini, OpenAI, Anthropic) is displayed
- **AND** the user can select one

### Requirement: Dynamic API Key Prompt
The wizard SHALL prompt for the correct API key environment variable based on the selected provider.

#### Scenario: Prompt for OpenAI key
- **WHEN** "OpenAI" is selected
- **THEN** the wizard prompts the user to set `OPENAI_API_KEY`
- **AND** displays instructions relevant to getting an OpenAI key

### Requirement: Dynamic Model Selection
The wizard SHALL display models relevant to the selected provider.

#### Scenario: OpenAI models
- **WHEN** "OpenAI" is selected
- **THEN** the model list shows "gpt-4o", "gpt-4-turbo", etc.

## MODIFIED Requirements

### Requirement: Config Generation
**Reason**: To support multi-provider architecture.
**Impact**: The generated `lango.json` MUST use the `providers` map instead of legacy `agent` fields for credentials.
