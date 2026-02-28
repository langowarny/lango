## MODIFIED Requirements

### Requirement: User Interface
The settings editor SHALL provide menu-based navigation with categories, free navigation between categories, and shared `tuicore.FormModel` for all forms. Provider and OIDC provider list views SHALL support managing collections. Pressing Esc at StepMenu SHALL navigate back to StepWelcome instead of quitting the TUI. The help bar at StepMenu SHALL display "Back" for the Esc key.

#### Scenario: Launch settings
- **WHEN** user runs `lango settings`
- **THEN** the editor SHALL display a welcome screen followed by the configuration menu

#### Scenario: Save from settings
- **WHEN** user selects "Save & Exit" from the menu
- **THEN** the configuration SHALL be saved as an encrypted profile

#### Scenario: Esc at Welcome screen quits
- **WHEN** user presses Esc at the Welcome screen (StepWelcome)
- **THEN** the TUI SHALL quit

#### Scenario: Esc at Menu navigates back to Welcome
- **WHEN** user presses Esc at the settings menu (StepMenu) while not in search mode
- **THEN** the editor SHALL navigate back to StepWelcome without quitting

#### Scenario: Esc at Menu during search cancels search
- **WHEN** user presses Esc at the settings menu while search mode is active
- **THEN** the search SHALL be cancelled and the menu SHALL remain at StepMenu

#### Scenario: Ctrl+C always quits
- **WHEN** user presses Ctrl+C at any step
- **THEN** the TUI SHALL quit immediately with Cancelled flag set

#### Scenario: Menu help bar shows Back for Esc
- **WHEN** the settings menu is displayed in normal mode (not searching)
- **THEN** the help bar SHALL display "Back" as the label for the Esc key
