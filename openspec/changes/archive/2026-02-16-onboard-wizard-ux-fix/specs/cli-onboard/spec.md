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

7.  **Knowledge**:
    - Toggle Knowledge System Enabled (boolean)
    - Set Max Learnings (integer)
    - Set Max Knowledge (integer)
    - Set Max Context Per Layer (integer)
    - Toggle Auto Approve Skills (boolean)
    - Set Max Skills Per Day (integer)

8.  **Providers**:
    - Add, edit, and delete multi-provider configurations

#### Scenario: Agent provider options from registered providers
- **WHEN** user navigates to Agent settings and providers are registered in config
- **THEN** the Provider and Fallback Provider dropdowns SHALL list the registered provider IDs

#### Scenario: Agent provider options with no registered providers
- **WHEN** user navigates to Agent settings and no providers are registered
- **THEN** the Provider dropdown SHALL fall back to default options: anthropic, openai, gemini, ollama

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

### User Interface
- **Navigation**:
    - Users MUST be able to navigate between configuration categories freely.
    - Uses a menu-based system (e.g., Main Menu -> Category -> Form).
    - The menu SHALL include categories: Agent, Server, Channels, Tools, Session, Security, Knowledge, Observational Memory, Embedding & RAG, Providers, Save & Exit, Cancel.
    - Form cursor navigation SHALL NOT panic when navigating past the first or last field.

#### Scenario: Form cursor at top boundary
- **WHEN** user presses up or shift+tab while cursor is on the first field
- **THEN** the cursor SHALL remain on the first field without error

#### Scenario: Session category in menu
- **WHEN** user views the configuration menu
- **THEN** "Session" category SHALL be listed before "Security"

#### Scenario: Provider creation field order
- **WHEN** user selects "Add New Provider" from the Providers list
- **THEN** the form SHALL display Type selector before Provider Name field
- **AND** the Provider Name field label SHALL be "Provider Name"

- **Validation**:
    - Input fields MUST validate data types (int, float, bool).
    - Port numbers MUST be within valid range (1-65535).
    - Essential fields (like Provider) MUST NOT be empty.
- **Feedback**:
    - Invalid inputs MUST display an error message immediately or upon submission.
    - Changes MUST be explicitly saved or discarded.
    - Provider list help footer SHALL display available key bindings including delete.

## REMOVED Requirements

### Requirement: DB Passphrase field in Security form
**Reason**: Passphrase is acquired via keyfile or terminal prompt during bootstrap, not stored in configuration. The field was dead code with no corresponding state handler.
**Migration**: No migration needed. Passphrase acquisition continues via existing bootstrap flow.
