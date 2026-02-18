# Upgrade Onboard TUI Design

## Goal
Implement a comprehensive TUI for `lango.json` configuration, replacing the linear wizard with a category-based editor.

## Architecture

### State Management
- **Refactor `WizardConfig`**: Instead of a simplified struct, the wizard will operate directly on a representation of `config.Config`.
- **Dirty State**: Track modified fields to highlight changes before saving.

### UI Structure (BubbleTea)
We will use a "Master-Detail" or "Menu-Form" structure.

1.  **Main Menu (Categories)**
    - Agent
    - Server
    - Tools
    - Channels
    - Security
    - *Save & Exit*
    - *Cancel*

2.  **Detail Views (Forms)**
    - Each category opens a form with relevant fields.
    - Fields support:
        - **TextInput**: Strings, numbers (with validation).
        - **Select/Cursor**: Choices (e.g., Provider model, Log level).
        - **Toggle**: Booleans (HTTP enabled, Headless mode).
        - **Multi-Select**: Channels list.

### Component Design
- **`Wizard` Model**: The central coordinator. Holds the `Config` state and current `View`.
- **`CategoryList`**: Component for the main menu.
- **`Forms`**: Reusable form components for editing config sections.
    - `AgentForm`
    - `ServerForm`
    - `ToolsForm`
    - ...

## UX Flow
1.  **Welcome Screen**:
    - "QuickStart" (Existing logic: defaults + API key only)
    - "Advanced Configuration" (Enters new Editor mode)

2.  **Editor Mode**:
    - User sees categories.
    - Selects "Agent".
    - Edits Provider (Select), Model (Text/Select), MaxTokens (Int).
    - Esc/Enter returns to Menu.
    - Selects "Save & Exit".
    - Validates all config.
    - Saves to `lango.json`.

## Technical Implementation

### File Structure
- `internal/cli/onboard/`
    - `wizard.go`: Main model and update loop.
    - `views.go`: View rendering logic.
    - `forms.go`: Form components for specific config sections.
    - `state.go`: Configuration state management.

### Validation Logic
- Re-use `internal/config/loader.go` validation logic where possible, but add granular field-level validation in the UI (e.g., immediate feedback on invalid port).

## Open Questions
- **Secret Handling**: Should we display existing secrets (masked)?
    - *Decision*: Display as `*****`. If edited, starts blank.

## Milestones
1.  Refactor Wizard to support new Navigation state.
2.  Implement Category Menu.
3.  Implement Forms for each section.
4.  Wire up Save logic.
