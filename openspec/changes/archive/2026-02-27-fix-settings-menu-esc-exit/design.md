## Context

The Settings TUI editor (`internal/cli/settings/editor.go`) uses Bubbletea for keyboard-driven navigation. Terminal arrow keys are transmitted as ANSI escape sequences (e.g., down arrow = `\x1b[B`). During rapid keystrokes, the escape byte `\x1b` can arrive in a separate read from `[B`, causing Bubbletea to interpret it as a standalone Esc keypress. At StepMenu, Esc was mapped directly to `tea.Quit`, resulting in accidental TUI exits during fast arrow key navigation.

## Goals / Non-Goals

**Goals:**
- Prevent accidental TUI exit when pressing arrow keys rapidly in StepMenu
- Maintain clean exit paths via ctrl+c and explicit menu actions (Save & Exit, Cancel)
- Keep the Esc key functional for intuitive back-navigation

**Non-Goals:**
- Modifying Bubbletea's escape sequence parsing (upstream concern)
- Adding debounce or timing-based key disambiguation
- Changing Esc behavior at StepWelcome (no arrow navigation, so no split-sequence risk)

## Decisions

**Decision 1: Esc at StepMenu navigates back to StepWelcome instead of quitting**

Rationale: StepWelcome has no arrow-key navigation, so the split-sequence bug cannot occur there. Users who intentionally press Esc can press it once more at StepWelcome to quit. This adds one extra keystroke for intentional exit but eliminates accidental exits entirely.

Alternative considered: Adding a confirmation dialog on Esc at StepMenu. Rejected because it adds UI complexity for a simple navigation fix.

Alternative considered: Debounce/timer to distinguish real Esc from split sequences. Rejected because it adds latency to all Esc handling and depends on timing heuristics.

**Decision 2: Update help bar text from "Quit" to "Back"**

Rationale: The help bar must accurately reflect the Esc key's behavior. Since Esc now navigates back rather than quitting, the label must change accordingly.

## Risks / Trade-offs

- [Extra keystroke for intentional quit] → Acceptable: Esc→Welcome→Esc→Quit is intuitive and standard for hierarchical menus
- [Users expecting Esc to quit from menu] → Mitigated by updated help bar showing "Back" and the Welcome screen still supporting Esc→Quit
