## Context

The settings menu in `internal/cli/settings/menu.go` originally stored all configuration categories in a flat `Categories []Category` slice on `MenuModel`. With 28+ categories, users had to scroll through the entire list without any visual grouping or search capability.

## Goals / Non-Goals

**Goals:**
- Group related categories under named sections for visual clarity
- Provide a keyboard-driven search feature to quickly find categories
- Highlight matching text in search results for discoverability
- Maintain backward compatibility — `allCategories()` still returns the full flat list for cursor navigation

**Non-Goals:**
- Fuzzy matching (exact substring matching is sufficient for ~30 items)
- Persistent search history or favorites
- Collapsible/expandable sections

## Decisions

**Decision 1: Section-based grouping with 6 named sections**

Categories are organized into: Core (4), Communication (4), AI & Knowledge (6), Infrastructure (4), P2P Network (5), Security (5), plus an untitled section for Save & Exit / Cancel. The grouping reflects the logical domain of each setting.

Rationale: Mirrors the mental model of the system architecture. Users can visually scan section headers to find the right area.

**Decision 2: `/` key activates search mode with `textinput.Model`**

Pressing `/` in normal mode focuses a `textinput.Model` at the top of the menu. Typing filters categories in real-time. Esc cancels search and restores the full grouped view. Enter selects the highlighted filtered result.

Rationale: `/` is a well-known convention (vim, GitHub, Slack). The `textinput.Model` from Bubbles provides cursor, styling, and input handling for free.

**Decision 3: Search matches against title, description, and ID**

`applyFilter()` converts both query and each field to lowercase, then uses `strings.Contains` for substring matching. Matching is case-insensitive.

Alternative considered: Fuzzy matching via a scoring library. Rejected because with ~30 categories, exact substring search is fast and predictable. Users can type any part of the category name or description.

**Decision 4: Amber/warning-colored highlight for matching substrings**

`highlightMatch()` finds the first occurrence of the query in each text field and renders it with `tui.Warning` (amber) color and bold. Selected items additionally get underline on the match. Non-matching text uses dim or accent styling based on selection state.

Rationale: Amber stands out from the default dim/accent palette without conflicting with error (red) or success (green) colors.

**Decision 5: Cursor navigates a flat list derived from sections**

`allCategories()` flattens all sections into a single `[]Category` slice. The cursor indexes into this flat list in normal mode, or into `filtered` in search mode. Section headers are not selectable — they are rendered as visual separators only.

Rationale: Keeps cursor logic simple. No need for a two-level index (section + item). The global index maps 1:1 to the rendered rows.

## Risks / Trade-offs

- [Flat cursor with sections] The cursor skips over section headers, so the visual gap between items at section boundaries may feel slightly odd. Acceptable because headers are clearly styled differently.
- [First-match-only highlighting] `highlightMatch` highlights only the first occurrence of the query in each string. Acceptable because category titles/descriptions are short.
