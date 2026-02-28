## ADDED Requirements

### Requirement: Inline field descriptions
All settings form fields SHALL include a `Description` string providing human-readable guidance. The description SHALL be shown only when the field is focused.

#### Scenario: Description displayed on focus
- **WHEN** the user navigates to a field with a Description
- **THEN** the form SHALL render the description text below that field

#### Scenario: Description hidden when not focused
- **WHEN** the user moves focus away from a field
- **THEN** the description for that field SHALL no longer be rendered

### Requirement: Field input validation
Numeric and range-sensitive fields SHALL have `Validate` functions that return clear error messages.

#### Scenario: Temperature validation
- **WHEN** the user enters a value outside 0.0-2.0 for the Temperature field
- **THEN** the validator SHALL return "must be between 0.0 and 2.0"

#### Scenario: Port validation
- **WHEN** the user enters a value outside 1-65535 for the Port field
- **THEN** the validator SHALL return "port out of range"

#### Scenario: Positive integer validation
- **WHEN** the user enters a non-positive value for fields requiring positive integers (Max Read Size, Max History Turns, Knowledge Max Context, Max Concurrent Jobs, Max Concurrent Tasks, Max Concurrent Steps, Max Peers, Observation Threshold, Max Bulk Import, Import Concurrency)
- **THEN** the validator SHALL return "must be a positive integer"

#### Scenario: Non-negative integer validation
- **WHEN** the user enters a negative value for fields allowing zero (Yield Time, Max Reflections in Context, Max Observations in Context, Inquiry Cooldown, Max Pending Inquiries, Approval Timeout, Embedding Dimensions, RAG Max Results)
- **THEN** the validator SHALL return "must be a non-negative integer" (with optional "(0 = unlimited)" suffix where applicable)

#### Scenario: Float range validation
- **WHEN** the user enters a value outside 0.0-1.0 for Min Trust Score
- **THEN** the validator SHALL return "must be between 0.0 and 1.0"

### Requirement: Auto-fetch model options from provider API
Form builders for Agent, Observational Memory, Embedding, and Librarian SHALL attempt to fetch available models from the configured provider API at form creation time.

#### Scenario: Successful model fetch
- **WHEN** the provider API returns a list of models within the 5-second timeout
- **THEN** the model field SHALL be converted from InputText to InputSelect with the fetched models as options, and the current model SHALL always be included

#### Scenario: Failed model fetch
- **WHEN** the provider API fails, times out, or returns empty
- **THEN** the model field SHALL remain as InputText with placeholder text

#### Scenario: Agent form model fetch
- **WHEN** the Agent form is created and the configured provider has a valid API key
- **THEN** the Model ID field SHALL be populated with models from `fetchModelOptions(cfg.Agent.Provider, ...)`

#### Scenario: Observational Memory model fetch with provider inheritance
- **WHEN** the Observational Memory form is created with an empty provider
- **THEN** the model fetch SHALL use the Agent provider as fallback

#### Scenario: Librarian model fetch with provider inheritance
- **WHEN** the Librarian form is created with an empty provider
- **THEN** the model fetch SHALL use the Agent provider as fallback

#### Scenario: Embedding model fetch
- **WHEN** the Embedding form is created with a non-empty provider
- **THEN** the Model field SHALL attempt to fetch models from the embedding provider

### Requirement: Unified embedding provider field
The Embedding & RAG form SHALL use a single "Provider" field (key `emb_provider_id`) mapped to `cfg.Embedding.Provider`. The state update handler SHALL clear the deprecated `cfg.Embedding.ProviderID` field when saving.

#### Scenario: Embedding form shows single provider field
- **WHEN** the user opens the Embedding & RAG form
- **THEN** the form SHALL display one "Provider" select field, not separate Provider and ProviderID fields

#### Scenario: State update clears deprecated ProviderID
- **WHEN** the `emb_provider_id` field is saved via UpdateConfigFromForm
- **THEN** `cfg.Embedding.Provider` SHALL be set to the value AND `cfg.Embedding.ProviderID` SHALL be set to empty string

### Requirement: Conditional field visibility in channel forms
Channel token fields SHALL be visible only when the parent channel is enabled.

#### Scenario: Telegram token hidden when disabled
- **WHEN** the Telegram Enabled toggle is unchecked
- **THEN** the Telegram Bot Token field SHALL be hidden

#### Scenario: Telegram token shown when enabled
- **WHEN** the user checks the Telegram Enabled toggle
- **THEN** the Telegram Bot Token field SHALL become visible

#### Scenario: Discord token visibility
- **WHEN** the Discord Enabled toggle is toggled
- **THEN** the Discord Bot Token field visibility SHALL match the toggle state

#### Scenario: Slack token visibility
- **WHEN** the Slack Enabled toggle is toggled
- **THEN** the Slack Bot Token and App Token fields visibility SHALL match the toggle state

### Requirement: Conditional visibility in security form
Security sub-fields SHALL be visible only when their parent toggle is enabled.

#### Scenario: PII fields hidden when interceptor disabled
- **WHEN** the Privacy Interceptor toggle is unchecked
- **THEN** all interceptor sub-fields (Redact PII, Approval Policy, Timeout, Notify Channel, Sensitive Tools, Exempt Tools, Disabled PII Patterns, Custom PII Patterns, Presidio) SHALL be hidden

#### Scenario: Presidio detail fields nested under both interceptor and presidio
- **WHEN** the interceptor is enabled but Presidio is disabled
- **THEN** the Presidio URL and Presidio Language fields SHALL be hidden

#### Scenario: Presidio fields visible when both enabled
- **WHEN** both the Privacy Interceptor and Presidio toggles are checked
- **THEN** the Presidio URL and Presidio Language fields SHALL be visible

#### Scenario: Signer Key ID visibility based on provider
- **WHEN** the signer provider is "local" or "enclave"
- **THEN** the Key ID field SHALL be hidden

#### Scenario: Signer RPC URL visibility
- **WHEN** the signer provider is "rpc"
- **THEN** the RPC URL field SHALL be visible

### Requirement: Conditional visibility in P2P sandbox form
P2P container sandbox fields SHALL be visible only when the container sandbox is enabled.

#### Scenario: Container fields hidden when container disabled
- **WHEN** the Container Sandbox Enabled toggle is unchecked
- **THEN** container-specific fields (Runtime, Image, Network Mode, Read-Only RootFS, CPU Quota, Pool Size, Pool Idle Timeout) SHALL be hidden

### Requirement: Conditional visibility in KMS form
KMS backend-specific fields SHALL be visible based on the selected backend type.

#### Scenario: Azure fields visible for azure-kv backend
- **WHEN** the KMS backend is "azure-kv"
- **THEN** the Azure Vault URL and Azure Key Version fields SHALL be visible

#### Scenario: PKCS11 fields visible for pkcs11 backend
- **WHEN** the KMS backend is "pkcs11"
- **THEN** the PKCS11 Module Path, Slot ID, PIN, and Key Label fields SHALL be visible

### Requirement: Model fetcher provider support
The `newProviderFromConfig` function SHALL support creating lightweight provider instances for: OpenAI, Anthropic, Gemini/Google, Ollama (via OpenAI-compatible endpoint), and GitHub (via OpenAI-compatible endpoint).

#### Scenario: Ollama default base URL
- **WHEN** creating an Ollama provider with empty BaseURL
- **THEN** the base URL SHALL default to "http://localhost:11434/v1"

#### Scenario: GitHub default base URL
- **WHEN** creating a GitHub provider with empty BaseURL
- **THEN** the base URL SHALL default to "https://models.inference.ai.azure.com"

#### Scenario: Provider without API key
- **WHEN** creating a non-Ollama provider with empty API key
- **THEN** `newProviderFromConfig` SHALL return nil
