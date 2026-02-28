## Why

The settings TUI forms are functional but lack intelligence -- fields have no inline help, no input validation, model IDs must be typed manually, the embedding config has a redundant Provider/ProviderID split, and all fields are always visible regardless of context. Users must guess valid ranges, remember exact model names, and wade through irrelevant options.

## What Changes

Five improvements to settings form UX, all within `internal/cli/settings/` and `internal/cli/tuicore/`:

1. **Inline descriptions** -- Every form field gets a human-readable `Description` string rendered below the focused field.
2. **Field validators** -- Numeric and range-sensitive fields get `Validate` functions with clear error messages (e.g. Temperature 0.0-2.0, port 1-65535, positive integers).
3. **Auto-fetch models** -- New `model_fetcher.go` queries provider `ListModels` API (5s timeout) and converts model text fields to select dropdowns. Supports OpenAI, Anthropic, Gemini, GitHub, Ollama. Falls back to text input on failure.
4. **Unify embedding provider** -- Merge `Embedding.Provider` and `Embedding.ProviderID` into single `Provider` field, clearing the deprecated `ProviderID` on save.
5. **Conditional field visibility** -- New `VisibleWhen func() bool` on Field struct. Channel tokens show only when the channel is enabled. Security PII fields show under interceptor enabled. Presidio fields nest under both interceptor + presidio enabled. P2P container fields show when container sandbox is enabled. KMS fields show based on backend type.

## Capabilities

### New Capabilities
- `cli-settings`: Inline descriptions for ~40 fields, validators for numeric fields, auto-fetched model dropdowns (Agent, Observational Memory, Embedding, Librarian), conditional visibility for channel tokens / security PII / Presidio / P2P container / KMS backend fields
- `cli-tuicore`: `VisibleWhen` field on `Field` struct, `IsVisible()` method, `VisibleFields()` on FormModel, cursor clamping after visibility changes, description rendering in form View

### Modified Capabilities
- `cli-settings`: Embedding form uses single `Provider` field (was Provider + ProviderID)
- `cli-tuicore`: FormModel cursor navigation operates on visible fields only; form View renders description below focused field

## Impact

- **New file**: `internal/cli/settings/model_fetcher.go` -- provider instantiation + model listing
- **Modified**: `internal/cli/settings/forms_impl.go` -- descriptions, validators, fetchModelOptions calls, VisibleWhen closures on ~15 fields
- **Modified**: `internal/cli/tuicore/field.go` -- Description field, VisibleWhen field, IsVisible() method
- **Modified**: `internal/cli/tuicore/form.go` -- VisibleFields(), cursor clamping, description rendering in View
- **Modified**: `internal/cli/tuicore/state_update.go` -- `emb_provider_id` case clears deprecated ProviderID
