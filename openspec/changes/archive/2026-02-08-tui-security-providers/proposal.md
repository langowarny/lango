# TUI Security & Providers Configuration

## Goal
Expand the TUI configuration capabilities to include full support for the `security` block (Interceptor, Signer, Passphrase) and the `providers` map, which are currently present in `lango.json` but missing from the onboarding interface.

## Background
The user noticed discrepancies between `lango.json` capabilities and the TUI:
1.  **Security**: The current TUI only configures Session settings in the "Security" menu. It lacks controls for:
    -   AI Privacy Interceptor (PII redaction, approval workflows)
    -   Secure Signer (RPC provider, Key ID)
    -   Local Passphrase encryption
2.  **Providers**: The `config.Config` struct has a `Providers` map (`map[string]ProviderConfig`) for multi-provider support, but the TUI only configures the single `Agent` block.

## Proposed Changes

### Security Form Expansion
Update `SecurityForm` to include:
-   **Interceptor**: Toggle Enabled, Redact PII, Approval Required.
-   **Signer**: Select Provider (Local, RPC), Input RPC URL/Key ID.
-   **Passphrase**: input field for local encryption passphrase.

### New Providers Form
Create a new `ProvidersForm` to manage the `providers` map:
-   List existing providers.
-   Add/Edit/Delete provider configurations.
-   Configure Type, API Key, and Base URL for each.

### Capabilities
-   `cli-onboard`: Enhance `security` configuration and add `providers` management.

## Impact
-   **Configuration**: Users can fully utilize security features without manual JSON editing.
-   **Extensibility**: Prepares the TUI for multi-provider workflows.
