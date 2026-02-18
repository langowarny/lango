## MODIFIED Requirements

### Requirement: API key security warning message
When API keys are stored inline (not using `${ENV_VAR}` references), the check SHALL report status `StatusWarn` with the message "Inline API keys for: <provider-list>". The details SHALL explain that keys in encrypted profiles are safe and suggest `${ENV_VAR}` references as an alternative for portability.

#### Scenario: Inline API key detected
- **WHEN** a provider has an API key that is not an `${ENV_VAR}` reference
- **THEN** check returns `StatusWarn` with message "Inline API keys for: <provider-id>"
- **THEN** details state that encrypted profile storage is safe and mention `${ENV_VAR}` as a portability option

### Requirement: API key security pass message
When all API keys use environment variable references, the check SHALL report status `StatusPass` with the message "All API keys secured".

#### Scenario: All keys use env var references
- **WHEN** all configured providers use `${ENV_VAR}` references for API keys
- **THEN** check returns `StatusPass` with message "All API keys secured"
