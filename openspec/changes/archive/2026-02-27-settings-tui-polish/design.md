## Context

The Settings TUI (`lango settings`) is a multi-step bubbletea application with five editor steps: Welcome, Menu, Form, ProvidersList, and AuthProvidersList. Before this change, each view rendered its own ad-hoc styles with no shared design language, no navigation breadcrumbs, and inconsistent help text.

## Goals / Non-Goals

**Goals:**
- Establish a centralized design token system (colors + reusable styles) in `internal/cli/tui/styles.go`
- Add breadcrumb navigation to all editor steps for spatial context
- Wrap list/menu views in bordered containers for visual grouping
- Provide a consistent help bar pattern across all interactive views

**Non-Goals:**
- No changes to form rendering (forms use `tuicore.FormModel` which has its own styling)
- No changes to key bindings or navigation logic
- No changes to onboard wizard (separate TUI)

## Decisions

1. **Design tokens in `tui` package, not `tuicore`**: The `tui` package (`internal/cli/tui/`) is the shared visual layer; `tuicore` holds form/config state logic. Design tokens belong in `tui` because they are consumed by both settings and onboard views.

2. **Breadcrumb function, not model**: `tui.Breadcrumb()` is a pure rendering function (segments in, styled string out) rather than a bubbletea model. It has no state and is called from `editor.View()`. This keeps it simple and composable.

3. **HelpBar as composable functions**: `HelpEntry(key, label)` renders a single badge + label, `HelpBar(entries...)` joins them. This avoids a struct-based approach and lets each view compose its own relevant entries.

4. **RoundedBorder for containers**: `lipgloss.RoundedBorder()` was chosen over `NormalBorder()` for a softer visual appearance consistent with the search bar style. Border color uses `tui.Muted` for non-focus containers, `tui.Primary` for the welcome box and search bar.

5. **Color palette values**: Colors follow a Tailwind-inspired palette for familiarity. `Primary` (#7C3AED, purple) is the brand color. `Accent` (#04B575, green) is for selection/focus states. `Dim` (#626262) is for secondary text. This avoids the common terminal pitfall of relying on ANSI colors that vary wildly across terminal emulators.

## Risks / Trade-offs

- [Terminal compatibility] Hex colors may not render correctly on terminals with fewer than 256 colors. Mitigation: lipgloss automatically degrades to the closest available color.
- [Style consistency] Other TUI views (onboard, doctor) do not yet use the new design tokens. Mitigation: The tokens are available for incremental adoption; no breaking change required.
