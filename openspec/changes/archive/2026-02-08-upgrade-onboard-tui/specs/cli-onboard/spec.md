# CLI Onboard Spec

## Goal
The `lango onboard` command must provide a comprehensive, interactive configuration editor that allows users to modify all aspects of their `lango.json` configuration file without manual editing.

## Requirements

### Configuration Coverage
The onboarding tool MUST support editing the following configuration sections:
1.  **Agent**:
    - Select Provider (Anthropic, OpenAI, Gemini, Ollama)
    - Set Model ID
    - Set Max Tokens (integer)
    - Set Temperature (float)
    - Set System Prompt Path (file path)

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
    - Configure Exec Tool: Default Timeout, Allow Background, Work Dir
    - Configure Filesystem Tool: Max Read Size
    - Configure Browser Tool: Headless mode, Session Timeout

5.  **Security**:
    - Toggle PII Redaction
    - Toggle Approval Requirement
    - Configure Passphrase (for local encryption)

### User Interface
- **Navigation**:
    - Users MUST be able to navigate between configuration categories freely.
    - Uses a menu-based system (e.g., Main Menu -> Category -> Form).
- **Validation**:
    - Input fields MUST validate data types (int, float, bool).
    - Port numbers MUST be within valid range (1-65535).
    - Essential fields (like Provider) MUST NOT be empty.
- **Feedback**:
    - Invalid inputs MUST display an error message immediately or upon submission.
    - Changes MUST be explicitly saved or discarded.

### Persistence
- Configuration MUST be saved to `lango.json`.
- Passwords/Secrets (API Keys, Tokens) MUST be handled securely (though typically stored in env vars, the config references them).
- The tool should generate a `.lango.env` template if new env vars are required.

## Success Criteria
1.  User can launch `lango onboard`, navigate to "Server" settings, change the port, save, and verify `lango.json` is updated.
2.  User can navigate to "Agent" settings, switch provider to Ollama, and save.
3.  Invalid inputs (e.g., Port 99999) are rejected by the UI.
