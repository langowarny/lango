## MODIFIED Requirements

### Requirement: Provider Configuration
The system SHALL allow configuring AI providers with either an API key or OAuth credentials.

#### Scenario: Provider with OAuth
- **WHEN** `lango.json` includes a provider with `clientId` and `clientSecret`
- **AND** a valid token exists in `~/.lango/tokens/<provider>.json`
- **THEN** system initializes the provider using the stored access token
- **AND** ignores any `apiKey` if present (or falls back if token invalid/missing)

#### Scenario: Provider with API Key (Legacy)
- **WHEN** `lango.json` includes a provider with `apiKey` only
- **THEN** system initializes the provider using the API key directly

#### Scenario: Provider without credentials
- **WHEN** neither `apiKey` nor valid OAuth token is available
- **THEN** system logs a warning or error during initialization
