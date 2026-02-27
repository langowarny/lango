## Why

The Settings TUI was functional but lacked visual polish. Users had no spatial context when navigating deeply nested menus, key bindings were inconsistent across views, and the interface had no unifying design language. This made the settings editor feel disconnected from the rest of the Lango CLI.

## What Changes

- Add `tui.Breadcrumb()` function for hierarchical navigation context (e.g., "Settings > Agent Configuration") displayed dynamically per editor step
- Wrap menu body, welcome screen, and provider/auth lists in `lipgloss.RoundedBorder()` styled containers
- Introduce a design system token layer in `internal/cli/tui/styles.go`: color constants (`Primary`, `Muted`, `Foreground`, `Accent`, `Dim`, `Warning`) and reusable styles (`SectionHeaderStyle`, `SeparatorLineStyle`, `CursorStyle`, `SearchBarStyle`)
- Add `tui.HelpBar()` and `tui.HelpEntry()` for consistent keyboard shortcut legends across all views (welcome, menu, providers list, auth providers list)

## Capabilities

### New Capabilities

- `tui-design-tokens`: Centralized color palette and reusable styles in `internal/cli/tui/styles.go`
- `tui-breadcrumbs`: Hierarchical navigation breadcrumbs for spatial orientation
- `tui-help-bars`: Consistent key legend bars using badge + label pattern

### Modified Capabilities

- `cli-settings`: Settings editor views updated to use breadcrumbs, styled containers, and help bars

## Impact

- **Files modified**: `internal/cli/tui/styles.go` (new design tokens), `internal/cli/settings/editor.go` (breadcrumbs + welcome box), `internal/cli/settings/menu.go` (container + help bars), `internal/cli/settings/providers_list.go` (container + help bars), `internal/cli/settings/auth_providers_list.go` (container + help bars)
- **No behavioral changes**: All modifications are purely visual. No config logic, data flow, or key bindings changed.
- **No new dependencies**: Uses existing `lipgloss` package already in the dependency tree.
