# TUI Security & Providers Design

## Goal
Implement interactive configuration for Security settings and Providers management in the `lango onboard` TUI.

## Architecture

### Security Form Expansion
The existing `SecurityForm` (currently only Session settings) will be split or sectioned into:
1.  **Session Settings**: Database Path, TTL.
2.  **Privacy Interceptor**:
    -   `Enabled` (Bool)
    -   `RedactPII` (Bool)
    -   `ApprovalRequired` (Bool)
3.  **Secure Signer**:
    -   `Provider` (Select: local, rpc)
    -   `KeyID` (Text)
    -   `RPCUrl` (Text)
4.  **Local Secrets**:
    -   `Passphrase` (Text, masked input ideally, or plain text warning)

### Providers Management
A new "Providers" menu item or sub-menu.
-   **List View**: Show configured providers (by ID).
-   **Add/Edit Form**:
    -   `ID` (Text, unique)
    -   `Type` (Select: openai, anthropic, etc.)
    -   `APIKey` (Text)
    -   `BaseURL` (Text)

## UX Flow

### Security
-   User selects "Security" from Main Menu.
-   Form now lists all security options.
-   Advanced options (Interceptor, Signer) might be collapsible or just listed below Session settings.

### Providers
-   User selects "Providers" from Main Menu.
-   Shows a list of existing providers (e.g., "anthropic", "ollama").
-   Actions:
    -   [Add New]
    -   Select existing to Edit.
    -   Select existing to Delete (optional).
-   Editing opens a `ProviderForm`.

## Technical Implementation

### Components
-   `internal/cli/onboard/forms_impl.go`:
    -   Update `NewSecurityForm` to include Interceptor/Signer fields.
    -   Create `NewProviderForm`.
-   `internal/cli/onboard/menu.go`:
    -   Add "Providers" to the main menu.
-   `internal/cli/onboard/state.go`:
    -   Update `UpdateConfigFromForm` to handle nested Security structs and the Providers map.
