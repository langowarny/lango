## ADDED Requirements

### Requirement: Breadcrumb navigation in settings editor
The settings editor SHALL display a breadcrumb navigation header that reflects the current editor step. The breadcrumb SHALL use `tui.Breadcrumb()` with the following segments per step:
- **StepWelcome / StepMenu**: "Settings"
- **StepForm**: "Settings" > form title (from `activeForm.Title`)
- **StepProvidersList**: "Settings" > "Providers"
- **StepAuthProvidersList**: "Settings" > "Auth Providers"

The last breadcrumb segment SHALL be rendered in `Primary` color with bold weight. Preceding segments SHALL be rendered in `Muted` color. Segments SHALL be separated by " > " in `Dim` color.

#### Scenario: Breadcrumb at menu
- **WHEN** user is at StepMenu
- **THEN** the breadcrumb SHALL display "Settings" as a single segment

#### Scenario: Breadcrumb at form
- **WHEN** user is editing the Agent form (StepForm)
- **THEN** the breadcrumb SHALL display "Settings > Agent Configuration"

#### Scenario: Breadcrumb at providers list
- **WHEN** user is at StepProvidersList
- **THEN** the breadcrumb SHALL display "Settings > Providers"

### Requirement: Styled containers for menu and list views
The settings menu body, providers list body, and auth providers list body SHALL each be wrapped in a `lipgloss.RoundedBorder()` container with `tui.Muted` border color and padding `(0, 1)`. The welcome screen SHALL be wrapped in a `lipgloss.RoundedBorder()` container with `tui.Primary` border color and padding `(1, 3)`.

#### Scenario: Menu container
- **WHEN** user is at StepMenu
- **THEN** the menu items SHALL be rendered inside a rounded-border container

#### Scenario: Welcome container
- **WHEN** user is at StepWelcome
- **THEN** the welcome message SHALL be rendered inside a primary-colored rounded-border box

### Requirement: Help bars in all interactive views
Every interactive settings view SHALL display a help bar at the bottom using `tui.HelpBar()` with `tui.HelpEntry()` badges. The help bars SHALL contain:
- **Welcome**: Enter (Start), Esc (Quit)
- **Menu (normal)**: up/down (Navigate), Enter (Select), / (Search), Esc (Back)
- **Menu (searching)**: up/down (Navigate), Enter (Select), Esc (Cancel)
- **Providers list**: up/down (Navigate), Enter (Select), d (Delete), Esc (Back)
- **Auth providers list**: up/down (Navigate), Enter (Select), d (Delete), Esc (Back)

#### Scenario: Menu help bar in normal mode
- **WHEN** user is at StepMenu in normal mode (not searching)
- **THEN** the help bar SHALL show Navigate, Select, Search, and Back entries

#### Scenario: Menu help bar in search mode
- **WHEN** user is at StepMenu in search mode
- **THEN** the help bar SHALL show Navigate, Select, and Cancel entries

### Requirement: Design system tokens in tui package
The `internal/cli/tui/styles.go` file SHALL export the following design tokens:
- **Colors**: `Primary` (#7C3AED), `Success` (#10B981), `Warning` (#F59E0B), `Error` (#EF4444), `Muted` (#6B7280), `Foreground` (#F9FAFB), `Background` (#1F2937), `Highlight` (#3B82F6), `Accent` (#04B575), `Dim` (#626262), `Separator` (#374151)
- **Styles**: `TitleStyle`, `SubtitleStyle`, `SuccessStyle`, `WarningStyle`, `ErrorStyle`, `MutedStyle`, `HighlightStyle`, `BoxStyle`, `ListItemStyle`, `SelectedItemStyle`, `SectionHeaderStyle`, `SeparatorLineStyle`, `CursorStyle`, `ActiveItemStyle`, `SearchBarStyle`, `FormTitleBarStyle`, `FieldDescStyle`
- **Functions**: `Breadcrumb(segments ...string)`, `HelpEntry(key, label string)`, `HelpBar(entries ...string)`, `KeyBadge(key string)`, `FormatPass(msg)`, `FormatWarn(msg)`, `FormatFail(msg)`, `FormatMuted(msg)`

#### Scenario: Breadcrumb rendering
- **WHEN** `tui.Breadcrumb("Settings", "Agent")` is called
- **THEN** the result SHALL be "Settings" in muted color, " > " separator in dim color, and "Agent" in primary bold

#### Scenario: HelpEntry rendering
- **WHEN** `tui.HelpEntry("Esc", "Back")` is called
- **THEN** the result SHALL be a key badge with "Esc" followed by "Back" label in dim color
