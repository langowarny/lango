## Context

The `lango onboard` TUI wizard has accumulated several UX issues: a panic when navigating to the top of a form, legacy fields that no longer serve any purpose, hardcoded provider lists, and a Security form overloaded with unrelated session settings. These issues were discovered during usability testing and code review.

All changes are confined to the `internal/cli/onboard` package and `internal/config/loader.go`, with no architectural changes needed.

## Goals / Non-Goals

**Goals:**
- Fix the index-out-of-range panic in form cursor navigation
- Remove dead DB Passphrase field from Security form
- Correct the default session database path to `data.db`
- Add provider deletion capability in the Providers list view
- Make Agent form provider dropdowns dynamically reflect registered providers
- Separate Session settings into their own menu category and form
- Improve Provider creation UX by reordering Type before ID

**Non-Goals:**
- Changing the underlying config data model or storage format
- Adding new configuration fields beyond what already exists
- Modifying the bootstrap or passphrase acquisition flow
- Adding confirmation dialogs for provider deletion (keep it simple for now)

## Decisions

### 1. Cursor bounds fix: stay at top vs wrap-around
**Decision**: Stay at position 0 when already at the top (do nothing on up/shift+tab).
**Rationale**: Wrap-around to the last field is disorienting in forms. Staying in place is the standard behavior users expect from form navigation.

### 2. Session form separation approach
**Decision**: Create a new `NewSessionForm()` function and "Session" menu category, rather than using sub-menus or tabs within Security.
**Rationale**: Follows the existing pattern of one form per menu category. Session and Security are conceptually distinct (data storage vs. access control). The `state_update.go` already handles session keys (`db_path`, `ttl`, `max_history_turns`) without any changes needed.

### 3. Provider deletion: immediate vs confirmation dialog
**Decision**: Immediate deletion on `d` key press with list refresh.
**Rationale**: The wizard operates on in-memory state that isn't persisted until "Save & Exit". Users can re-add a provider if deleted accidentally. Adding a confirmation dialog would add complexity without significant benefit at this stage.

### 4. Dynamic provider options: shared helper function
**Decision**: Extract `buildProviderOptions()` helper used by both Provider and Fallback Provider fields.
**Rationale**: Avoids duplicating the map iteration logic. Falls back to hardcoded defaults when no providers are registered, ensuring the form is always usable.

### 5. Provider form field ordering
**Decision**: Move Type selector before ID field for new providers only.
**Rationale**: Users typically decide the provider type first, then name it. The Type value can inform the ID choice. Editing existing providers doesn't show the ID field at all, so this only affects creation.

## Risks / Trade-offs

- **[Risk] Provider deletion has no undo** → Mitigation: Changes are in-memory only until Save & Exit; user can cancel the wizard to discard all changes.
- **[Risk] Dynamic provider list may be empty** → Mitigation: Falls back to hardcoded defaults when no providers are registered.
- **[Risk] Existing tests expect old Security form field count** → Mitigation: Update `TestNewSecurityForm_AllFields` and add `TestNewSessionForm_AllFields`.
