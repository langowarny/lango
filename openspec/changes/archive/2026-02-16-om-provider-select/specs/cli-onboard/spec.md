## MODIFIED Requirements

### Requirement: Observational Memory onboard form
The system SHALL include an Observational Memory configuration form in the onboard TUI wizard. The form SHALL have fields for: enabled (bool), provider (select, dynamically populated from registered providers with empty option for agent default), model (text), message token threshold (int, positive validation), observation token threshold (int, positive validation), and max message token budget (int, positive validation). The menu entry SHALL appear between Knowledge and Providers with the label "Observational Memory".

#### Scenario: Navigate to OM form
- **WHEN** user selects "Observational Memory" from the onboard menu
- **THEN** the wizard displays the OM configuration form with current values from config

#### Scenario: Edit OM settings
- **WHEN** user modifies fields in the OM form and presses ESC
- **THEN** the changes are saved to the in-memory config state and the wizard returns to the menu

#### Scenario: Invalid threshold value
- **WHEN** user enters a non-positive number in a threshold field
- **THEN** the form displays a validation error "must be a positive integer"

#### Scenario: OM provider field is a select dropdown
- **WHEN** user navigates to Observational Memory settings
- **THEN** the provider field SHALL be an InputSelect dropdown
- **AND** the first option SHALL be an empty string representing "use agent default"
- **AND** subsequent options SHALL be registered provider IDs from buildProviderOptions

#### Scenario: OM provider options with no registered providers
- **WHEN** user navigates to Observational Memory settings and no providers are registered
- **THEN** the provider dropdown SHALL fall back to default options: empty string, anthropic, openai, gemini, ollama
