## Why

The `lango settings` menu lists 28+ configuration categories in a flat, unsorted list. As the number of settings categories has grown (Providers, Agent, Server, Channels, ... through Security KMS), users must scroll through the entire list to find the category they want. There is no way to quickly jump to a category by name, and the lack of visual grouping makes it hard to understand which categories are related.

## What Changes

- Restructure the menu from a flat `[]Category` list into grouped `[]Section` with 6 logical headings: Core, Communication, AI & Knowledge, Infrastructure, P2P Network, Security (plus an untitled section for Save & Exit / Cancel)
- Add a keyword search feature activated by pressing `/` that filters categories in real-time by matching against title, description, and ID
- Highlight matching substrings in search results with amber/warning color
- Render section headers with visual separators between groups

## Capabilities

### New Capabilities

- `cli-settings`: Grouped section layout for the settings menu
- `cli-settings`: Keyword search with `/` activation, real-time filtering, and match highlighting

### Modified Capabilities

- `cli-settings`: Menu categories are now organized under `Section` groupings instead of a flat list

## Impact

- `internal/cli/settings/menu.go`: Major rewrite — added `Section` struct, `searchInput textinput.Model`, `filtered []Category`, `applyFilter()`, `highlightMatch()`, `renderGroupedView()`, `renderFilteredView()` methods
