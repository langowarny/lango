## Why

When pressing the down arrow key rapidly in the Settings menu, the TUI exits unexpectedly. Terminal arrow keys are sent as escape sequences (`\x1b[B`), and during rapid input the `\x1b` byte can arrive separately, being interpreted as a standalone Esc keypress. Since StepMenu maps Esc directly to `tea.Quit`, this causes accidental TUI termination.

## What Changes

- Change Esc behavior at StepMenu from quitting the TUI to navigating back to the Welcome screen
- Update the help bar text at StepMenu from "Quit" to "Back" to reflect the new behavior
- Add editor navigation tests covering Esc behavior at each step and ctrl+c quit behavior

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `cli-settings`: Change Esc key at StepMenu from quit to back-navigation to Welcome screen

## Impact

- `internal/cli/settings/editor.go`: Esc at StepMenu returns to StepWelcome instead of quitting
- `internal/cli/settings/menu.go`: Help bar text updated from "Quit" to "Back"
- `internal/cli/settings/editor_test.go`: New test file with 4 test cases for editor navigation
