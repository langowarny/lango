# TUI Security & Providers Spec

## Goal
The `lango onboard` TUI MUST support configuration of advanced Security features (Interceptor, Signer, Passphrase) and management of multiple AI Providers.

## Requirements

### Security Configuration
The Security form MUST be expanded to include the following sections:

#### Privacy Interceptor
-   **Enable/Disable**: Toggle switch for the interceptor.
-   **Redact PII**: Toggle switch for automatic PII redaction.
-   **Approval Required**: Toggle switch for requiring human approval for sensitive actions.

#### Secure Signer
-   **Provider**: Select input (options: "local", "rpc").
-   **Tethering**:
    -   Host connection configuration if "rpc" is selected.
    -   Key ID specification.

#### Local Secrets
-   **Passphrase**: Masked text input for setting the local encryption passphrase.

### Providers Configuration
A new "Providers" section MUST be added to the main menu.

#### Provider List
-   Display a list of all configured providers by their ID (e.g., "anthropic", "custom-openai").
-   Allow selection of a provider to edit.
-   Allow adding a new provider.

#### Provider Form
-   **ID**: Text input (unique identifier).
-   **Type**: Select input (options: "openai", "anthropic", "gemini", "ollama").
-   **API Key**: Text input (standard environment variable handling applies).
-   **Base URL**: Text input (optional, for OpenAI-compatible providers).
-   **Model**: Text input (default model for this provider).

## Constraints
-   All new configuration MUST be saved to `lango.json` correctly mapped to the `security` block and `providers` map.
-   The "Agent" configuration block remains the *default* active provider/model configuration.
