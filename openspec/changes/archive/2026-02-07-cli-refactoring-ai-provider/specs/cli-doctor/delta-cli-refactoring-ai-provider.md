# Delta Spec: CLI Doctor - Multi-Provider

## Overview

Updates the doctor command to verify all configured providers.

## ADDED Requirements

### Requirement: Verify all providers
The command SHALL verify the status of every provider defined in the `providers` configuration map.

#### Scenario: Multiple providers configured
- **WHEN** `lango.json` contains both "openai" and "anthropic" in `providers`
- **THEN** the doctor output includes checks for both "OpenAI" and "Anthropic"

## MODIFIED Requirements

### Requirement: Legacy API Key Check
**Reason**: Replaced by provider-specific verification.
**Impact**: The check for `GOOGLE_API_KEY` SHALL ONLY run if a Gemini provider is configured or no providers are configured at all (fallback behavior).
