## 1. Menu Restructuring

- [x] 1.1 Add `Section` struct with `Title string` and `Categories []Category`
- [x] 1.2 Change `MenuModel.Sections` from `[]Category` to `[]Section`
- [x] 1.3 Populate 6 named sections (Core, Communication, AI & Knowledge, Infrastructure, P2P Network, Security) plus untitled Save/Cancel section in `NewMenuModel()`
- [x] 1.4 Add `allCategories()` method to flatten sections into a single `[]Category` slice
- [x] 1.5 Add `renderGroupedView()` to render section headers with separator lines

## 2. Search Feature

- [x] 2.1 Add `searching bool`, `searchInput textinput.Model`, `filtered []Category` fields to `MenuModel`
- [x] 2.2 Initialize `textinput.Model` in `NewMenuModel()` with `/ ` prompt, placeholder, and styling
- [x] 2.3 Add `/` key handler in normal mode to activate search (focus input, reset cursor)
- [x] 2.4 Add search-mode key handling: Esc cancels, Enter selects, up/down navigates, default forwards to text input
- [x] 2.5 Implement `applyFilter()` — case-insensitive substring match on title, desc, and ID
- [x] 2.6 Add `selectableItems()` helper to return filtered or full list based on mode
- [x] 2.7 Implement `renderFilteredView()` with "No matching items" empty state

## 3. Search Highlighting

- [x] 3.1 Implement `highlightMatch()` — finds first match, renders with `tui.Warning` bold, underline if selected
- [x] 3.2 Integrate highlighting into `renderItem()` when search query is active

## 4. Help Bar and View

- [x] 4.1 Update `View()` to show search bar (active input or dim hint) at top
- [x] 4.2 Update help footer to show Search(`/`) in normal mode, Cancel(`Esc`) in search mode
- [x] 4.3 Add `IsSearching()` and `AllCategories()` public accessors

## 5. Verification

- [x] 5.1 Run `go build ./...` — no compilation errors
- [x] 5.2 Run `go test ./internal/cli/settings/...` — all tests pass
