## Context

The settings TUI editor renders forms via `tuicore.FormModel` containing `[]*Field`. Prior to this change, fields had no help text, no validation, model IDs were free-text only, embedding config carried a redundant field, and every field was unconditionally visible.

## Goals / Non-Goals

**Goals:**
- Give users inline guidance for every field without leaving the form
- Prevent invalid input at entry time with clear error messages
- Auto-populate model selection from live provider APIs where possible
- Simplify the embedding config by removing the redundant ProviderID field
- Reduce visual noise by hiding irrelevant fields based on current state

**Non-Goals:**
- Changing config struct shapes (only the embedding ProviderID deprecation)
- Adding new settings categories or menu items
- Real-time re-fetching of models when provider selection changes (fetch happens at form creation time)

## Decisions

### 1. Description rendered only for focused field
**Decision**: Show the `Description` string below the currently focused field only, prefixed with an info icon.

**Rationale**: Showing all descriptions at once would make forms too tall. Focused-only display keeps the form compact while providing help exactly when needed.

### 2. VisibleWhen as closure on Field
**Decision**: Add `VisibleWhen func() bool` to `Field`. When non-nil, the field is shown only when the function returns true. The closure captures a pointer to the controlling field (e.g. `telegramEnabled`), so toggling the parent immediately hides/shows dependent fields.

**Rationale**: Closures avoid the need for a declarative dependency graph or string-based key references. Since the controlling field is defined in the same function scope, type safety is preserved.

### 3. Cursor operates on visible fields only
**Decision**: `FormModel.Update()` and `View()` call `VisibleFields()` to get the filtered slice. The cursor indexes into this slice. After any toggle, the cursor is clamped to `len(visible)-1` to prevent out-of-bounds.

**Rationale**: Navigating to hidden fields would be confusing. Re-evaluating visibility after every toggle ensures the cursor stays valid even when a toggle hides the currently focused field.

### 4. Model fetching at form creation time with timeout
**Decision**: `fetchModelOptions()` runs synchronously during `NewXxxForm()` with a 5-second context timeout. If it fails or returns empty, the field remains a text input.

**Rationale**: Synchronous fetch is simpler than async (no loading state needed in the TUI). The 5s timeout prevents blocking the UI. Graceful fallback to text input means the form always works even without network.

### 5. Provider instantiation without full wiring
**Decision**: `newProviderFromConfig()` creates a lightweight provider instance from config alone (API key + base URL), without the full application bootstrap. It returns nil if the provider cannot be created.

**Rationale**: The settings editor runs before the application is fully wired. We only need `ListModels()`, which requires minimal provider state.

### 6. Embedding Provider field unification
**Decision**: The embedding form uses field key `emb_provider_id` mapped to `cfg.Embedding.Provider`. The state update handler also clears `cfg.Embedding.ProviderID` to empty.

**Rationale**: The config previously had both `Provider` (display name) and `ProviderID` (used for lookups). They were always set to the same value, causing confusion. Unifying into `Provider` and clearing `ProviderID` on save is backward-compatible.

## Risks / Trade-offs

- **Synchronous model fetch adds latency**: Up to 5s delay when opening forms with model fields. Acceptable because it only happens at form creation, not during navigation, and the fallback is text input.
- **Closure memory**: VisibleWhen closures capture field pointers and stay in memory for the form lifetime. Negligible because form lifetimes are short.
- **No cross-field re-fetch**: Changing the provider dropdown does not re-fetch models for the model field. The user must exit and re-enter the form. Acceptable for v1; async re-fetch can be added later.
