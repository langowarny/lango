## ADDED Requirements

### Requirement: Browser config fields exposed in TUI
The Onboard TUI Tools form SHALL expose the `enabled` and `sessionTimeout` fields for browser tool configuration.

#### Scenario: Browser enabled toggle in TUI
- **WHEN** user navigates to Tools configuration in the onboard wizard
- **THEN** a "Browser Enabled" boolean toggle SHALL be displayed before the "Browser Headless" toggle

#### Scenario: Browser session timeout in TUI
- **WHEN** user navigates to Tools configuration in the onboard wizard
- **THEN** a "Browser Session Timeout" duration text field SHALL be displayed after the "Browser Headless" toggle
- **AND** the field SHALL accept Go duration strings (e.g., "5m", "10m")
