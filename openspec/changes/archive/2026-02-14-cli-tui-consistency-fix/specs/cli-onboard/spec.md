## MODIFIED Requirements

### Configuration Coverage
The onboarding tool MUST support editing the following configuration sections:
1.  **Agent**:
    - Select Provider (Anthropic, OpenAI, Gemini, Ollama)
    - Set Model ID
    - Set Max Tokens (integer)
    - Set Temperature (float)
    - Set System Prompt Path (file path)
    - Select Fallback Provider (empty, Anthropic, OpenAI, Gemini, Ollama)
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

5.  **Security**:
    - Set Session DB Path
    - Set Session TTL (duration)
    - Set Max History Turns (integer)
    - Toggle Privacy Interceptor Enabled
    - Toggle PII Redaction
    - Toggle Approval Requirement
    - Select Signer Provider (local, rpc, enclave)
    - Set RPC URL
    - Set Key ID
    - Configure Passphrase (for local encryption)

6.  **Knowledge**:
    - Toggle Knowledge System Enabled (boolean)
    - Set Max Learnings (integer)
    - Set Max Knowledge (integer)
    - Set Max Context Per Layer (integer)
    - Toggle Auto Approve Skills (boolean)
    - Set Max Skills Per Day (integer)

7.  **Providers**:
    - Add, edit, remove multi-provider configurations

#### Scenario: Agent fallback configuration
- **WHEN** user navigates to Agent settings
- **THEN** the form SHALL display fields for system_prompt_path, fallback_provider, and fallback_model
- **AND** fallback_provider SHALL be an InputSelect with options: empty, anthropic, openai, gemini, ollama

#### Scenario: Browser tool fields in Tools form
- **WHEN** user navigates to Tools settings
- **THEN** the form SHALL display browser_enabled toggle before browser_headless
- **AND** the form SHALL display browser_session_timeout as a duration text field after browser_headless

#### Scenario: Session max history in Security form
- **WHEN** user navigates to Security settings
- **THEN** the form SHALL display max_history_turns as an integer field after Session TTL

#### Scenario: Knowledge menu and form
- **WHEN** user views the Configuration Menu
- **THEN** a "Knowledge" category SHALL appear between Security and Providers
- **AND** selecting it SHALL display the Knowledge configuration form with 6 fields

## MODIFIED Requirements

### User Interface
- **Navigation**:
    - Users MUST be able to navigate between configuration categories freely.
    - Uses a menu-based system (e.g., Main Menu -> Category -> Form).
    - The menu SHALL include categories: Agent, Server, Channels, Tools, Security, Knowledge, Providers, Save & Exit, Cancel.

#### Scenario: Knowledge category in menu
- **WHEN** user views the configuration menu
- **THEN** "Knowledge" category SHALL be listed after "Security" and before "Providers"

## ADDED Requirements

### Requirement: Onboard command description reflects actual flow
The `lango onboard` command Long description SHALL accurately list all configurable sections.

#### Scenario: Long description content
- **WHEN** user runs `lango onboard --help`
- **THEN** the description SHALL list Agent, Server, Channels, Tools, Security, Knowledge, and Providers as configurable sections
