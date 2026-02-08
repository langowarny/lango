## ADDED Requirements

### Requirement: Onboard Command Entry Point
The system SHALL provide a `lango onboard` command that guides users through initial setup.

#### Scenario: Running onboard command
- **WHEN** user executes `lango onboard`
- **THEN** system displays interactive TUI wizard

#### Scenario: Running onboard with existing config
- **WHEN** lango.json already exists
- **THEN** system asks whether to keep, modify, or reset existing configuration

### Requirement: Welcome Screen
The system SHALL display a welcome screen with mode selection.

#### Scenario: Mode selection
- **WHEN** onboard wizard starts
- **THEN** user can choose between "QuickStart" and "Advanced" modes

### Requirement: API Key Configuration Step
The system SHALL collect and validate the AI provider API key.

#### Scenario: API key input
- **WHEN** user is on API key step
- **THEN** system displays masked input field for API key

#### Scenario: API key validation
- **WHEN** user submits API key
- **THEN** system validates the key with the provider if network is available

#### Scenario: Skip validation offline
- **WHEN** network is unavailable during API key step
- **THEN** system skips validation with warning and proceeds

### Requirement: Model Selection Step
The system SHALL allow users to select their default AI model.

#### Scenario: Model dropdown
- **WHEN** user is on model selection step
- **THEN** system displays dropdown with available models for the configured provider

#### Scenario: Gemini models available
- **WHEN** API key is for Gemini
- **THEN** available models include "gemini-2.0-flash-exp", "gemini-1.5-pro", "gemini-1.5-flash"

### Requirement: Channel Setup Step
The system SHALL guide users through enabling one messaging channel.

#### Scenario: Channel selection
- **WHEN** user is on channel setup step
- **THEN** system displays options for Telegram, Discord, Slack, or "Skip for now"

#### Scenario: Telegram setup
- **WHEN** user selects Telegram
- **THEN** system prompts for bot token and provides link to @BotFather

#### Scenario: Discord setup
- **WHEN** user selects Discord
- **THEN** system prompts for bot token and provides Discord Developer Portal instructions

#### Scenario: Slack setup
- **WHEN** user selects Slack
- **THEN** system prompts for bot token and app token

### Requirement: Configuration Save
The system SHALL save the completed configuration to lango.json.

#### Scenario: Save new configuration
- **WHEN** user completes all steps
- **THEN** system writes lango.json with collected values

#### Scenario: Environment variable hints
- **WHEN** configuration is saved
- **THEN** system displays which environment variables to set (e.g., GOOGLE_API_KEY)

### Requirement: Post-Setup Verification
The system SHALL offer to run doctor after onboarding completes.

#### Scenario: Offer doctor run
- **WHEN** onboarding completes successfully
- **THEN** system asks "Run lango doctor to verify setup?"

### Requirement: Keyboard Navigation
The system SHALL support keyboard-only navigation through the wizard.

#### Scenario: Arrow key navigation
- **WHEN** user is in selection step
- **THEN** up/down arrows navigate options, Enter selects

#### Scenario: Cancel with Escape
- **WHEN** user presses Escape during any step
- **THEN** system confirms exit and allows cancellation
