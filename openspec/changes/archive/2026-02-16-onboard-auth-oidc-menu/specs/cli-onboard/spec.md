## ADDED Requirements

### Requirement: Auth OIDC provider management in onboard TUI
The onboard wizard SHALL include an "Auth" menu category for managing OIDC provider configurations. The Auth menu SHALL open an OIDC provider list view that allows users to add, edit, and delete OIDC providers. Each OIDC provider form SHALL include fields for: Provider Name (text, new providers only), Issuer URL (text), Client ID (password), Client Secret (password), Redirect URL (text), and Scopes (text, comma-separated).

#### Scenario: Auth category in menu
- **WHEN** user views the configuration menu
- **THEN** an "Auth" category SHALL appear after "Security" with description "OIDC provider configuration"

#### Scenario: Navigate to OIDC provider list
- **WHEN** user selects "Auth" from the configuration menu
- **THEN** the wizard SHALL display the "Manage OIDC Providers" list view
- **AND** the list SHALL show existing OIDC providers with their ID and Issuer URL

#### Scenario: Add new OIDC provider
- **WHEN** user selects "+ Add New OIDC Provider" from the auth providers list
- **THEN** the wizard SHALL display a form titled "Add New OIDC Provider"
- **AND** the form SHALL include a "Provider Name" text field

#### Scenario: Edit existing OIDC provider
- **WHEN** user selects an existing OIDC provider from the auth providers list
- **THEN** the wizard SHALL display a form titled "Edit OIDC Provider: <id>"
- **AND** the form SHALL be pre-populated with the provider's current values
- **AND** the form SHALL NOT include a "Provider Name" field

#### Scenario: Delete OIDC provider
- **WHEN** user presses `d` on an OIDC provider in the list
- **THEN** the provider SHALL be removed from the in-memory auth config
- **AND** the list SHALL refresh to reflect the deletion

#### Scenario: Cannot delete Add New option
- **WHEN** user presses `d` while the cursor is on "+ Add New OIDC Provider"
- **THEN** no deletion SHALL occur

#### Scenario: OIDC scopes stored as string slice
- **WHEN** user enters scopes as "openid,email,profile" in the OIDC provider form
- **THEN** the scopes SHALL be stored as `["openid", "email", "profile"]` in config

#### Scenario: ESC from OIDC form saves and returns to list
- **WHEN** user presses ESC while in an OIDC provider form
- **THEN** the form values SHALL be saved to the in-memory config
- **AND** the wizard SHALL return to the OIDC provider list view

#### Scenario: ESC from OIDC provider list returns to menu
- **WHEN** user presses ESC while in the OIDC provider list
- **THEN** the wizard SHALL return to the configuration menu

## MODIFIED Requirements

### Configuration Coverage
The onboarding tool MUST support editing the following configuration sections:
1.  **Agent**:
    - Select Provider (dynamically populated from registered providers; falls back to Anthropic, OpenAI, Gemini, Ollama when no providers are registered)
    - Set Model ID
    - Set Max Tokens (integer)
    - Set Temperature (float)
    - Set System Prompt Path (file path)
    - Select Fallback Provider (empty + same dynamic provider list as Provider)
    - Set Fallback Model ID

2.  **Server**:
    - Set Host (default: localhost)
    - Set Port (integer, 1-65535)
    - Toggle HTTP Enabled (boolean)
    - Toggle WebSocket Enabled (boolean)

3.  **Channels**:
    - Enable/Disable each supported channel (Telegram, Discord, Slack)
    - Set Bot Tokens for enabled channels
    - Set App Token/Signing Secret for Slack if enabled

4.  **Tools**:
    - Configure Exec Tool: Default Timeout, Allow Background
    - Toggle Browser Enabled (boolean)
    - Toggle Browser Headless mode (boolean)
    - Set Browser Session Timeout (duration)
    - Configure Filesystem Tool: Max Read Size

5.  **Session**:
    - Set Database Path (default: `~/.lango/data.db`)
    - Set Session TTL (duration)
    - Set Max History Turns (integer)

6.  **Security**:
    - Toggle Privacy Interceptor Enabled
    - Toggle PII Redaction
    - Toggle Approval Requirement
    - Set Approval Timeout (integer, non-negative)
    - Select Notify Channel (empty, telegram, discord, slack)
    - Set Sensitive Tools (comma-separated)
    - Select Signer Provider (local, rpc, enclave)
    - Set RPC URL
    - Set Key ID

7.  **Auth**:
    - Add, edit, and delete OIDC provider configurations
    - Each OIDC provider: Issuer URL, Client ID, Client Secret, Redirect URL, Scopes

8.  **Knowledge**:
    - Toggle Knowledge System Enabled (boolean)
    - Set Max Learnings (integer)
    - Set Max Knowledge (integer)
    - Set Max Context Per Layer (integer)
    - Toggle Auto Approve Skills (boolean)
    - Set Max Skills Per Day (integer)

9.  **Providers**:
    - Add, edit, and delete multi-provider configurations

#### Scenario: Agent fallback configuration
- **WHEN** user navigates to Agent settings
- **THEN** the form SHALL display fields for system_prompt_path, fallback_provider, and fallback_model
- **AND** fallback_provider SHALL be an InputSelect with options: empty + registered provider IDs

#### Scenario: Agent provider options from registered providers
- **WHEN** user navigates to Agent settings and providers are registered in config
- **THEN** the Provider and Fallback Provider dropdowns SHALL list the registered provider IDs

#### Scenario: Agent provider options with no registered providers
- **WHEN** user navigates to Agent settings and no providers are registered
- **THEN** the Provider dropdown SHALL fall back to default options: anthropic, openai, gemini, ollama

#### Scenario: Browser tool fields in Tools form
- **WHEN** user navigates to Tools settings
- **THEN** the form SHALL display browser_enabled toggle before browser_headless
- **AND** the form SHALL display browser_session_timeout as a duration text field after browser_headless

#### Scenario: Session settings in dedicated form
- **WHEN** user selects "Session" from the configuration menu
- **THEN** the wizard SHALL display a Session form with Database Path, Session TTL, and Max History Turns fields

#### Scenario: Security form without session fields
- **WHEN** user selects "Security" from the configuration menu
- **THEN** the form SHALL NOT include Database Path, Session TTL, Max History Turns, or DB Passphrase fields

#### Scenario: Delete provider from list
- **WHEN** user presses `d` on a provider in the Providers list view
- **THEN** the provider SHALL be removed from the in-memory config state
- **AND** the list SHALL refresh to reflect the deletion

#### Scenario: Cannot delete Add New Provider option
- **WHEN** user presses `d` while the cursor is on "Add New Provider"
- **THEN** no deletion SHALL occur

#### Scenario: Knowledge menu and form
- **WHEN** user views the Configuration Menu
- **THEN** a "Knowledge" category SHALL appear between Security and Observational Memory
- **AND** selecting it SHALL display the Knowledge configuration form with 6 fields

### User Interface
- **Navigation**:
    - Users MUST be able to navigate between configuration categories freely.
    - Uses a menu-based system (e.g., Main Menu -> Category -> Form).
    - The menu SHALL include categories in this order: Providers, Agent, Server, Channels, Tools, Session, Security, Auth, Knowledge, Observational Memory, Embedding & RAG, Save & Exit, Cancel.
    - Form cursor navigation SHALL NOT panic when navigating past the first or last field.

#### Scenario: Providers category appears first
- **WHEN** user views the configuration menu
- **THEN** "Providers" SHALL be the first category in the menu, before "Agent"

#### Scenario: Form cursor at top boundary
- **WHEN** user presses up or shift+tab while cursor is on the first field
- **THEN** the cursor SHALL remain on the first field without error

#### Scenario: Session category in menu
- **WHEN** user views the configuration menu
- **THEN** "Session" category SHALL be listed before "Security"

#### Scenario: Knowledge category in menu
- **WHEN** user views the configuration menu
- **THEN** "Knowledge" category SHALL be listed after "Security" and before "Observational Memory"

#### Scenario: Auth category in menu
- **WHEN** user views the configuration menu
- **THEN** "Auth" category SHALL be listed after "Security" and before "Knowledge"

#### Scenario: Provider creation field order
- **WHEN** user selects "Add New Provider" from the Providers list
- **THEN** the form SHALL display Type selector before Provider Name field
- **AND** the Provider Name field label SHALL be "Provider Name"
