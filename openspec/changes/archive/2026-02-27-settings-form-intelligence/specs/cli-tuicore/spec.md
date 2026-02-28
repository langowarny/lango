## ADDED Requirements

### Requirement: Field Description property
The `Field` struct SHALL include a `Description string` property for inline help text.

#### Scenario: Description stored on field
- **WHEN** a Field is created with a Description value
- **THEN** the Description SHALL be accessible on the field instance

### Requirement: VisibleWhen conditional visibility
The `Field` struct SHALL include a `VisibleWhen func() bool` property. When non-nil, the field is shown only when the function returns true. When nil, the field is always visible.

#### Scenario: VisibleWhen nil means always visible
- **WHEN** a Field has `VisibleWhen` set to nil
- **THEN** `IsVisible()` SHALL return true

#### Scenario: VisibleWhen returns false hides field
- **WHEN** a Field has `VisibleWhen` returning false
- **THEN** `IsVisible()` SHALL return false and the field SHALL not appear in `VisibleFields()`

#### Scenario: VisibleWhen dynamically responds to state
- **WHEN** a VisibleWhen closure captures a pointer to a parent field's Checked state
- **THEN** toggling the parent field SHALL immediately affect the child field's visibility on next `VisibleFields()` call

### Requirement: IsVisible method on Field
The `Field` struct SHALL expose an `IsVisible() bool` method that returns true when `VisibleWhen` is nil, and the result of `VisibleWhen()` otherwise.

### Requirement: VisibleFields on FormModel
`FormModel` SHALL expose a `VisibleFields() []*Field` method that returns only fields where `IsVisible()` returns true.

#### Scenario: VisibleFields filters hidden fields
- **WHEN** a form has 5 fields and 2 have VisibleWhen returning false
- **THEN** VisibleFields() SHALL return 3 fields

## MODIFIED Requirements

### Requirement: FormModel cursor navigation (MODIFIED)
The form cursor SHALL index into `VisibleFields()` instead of the full `Fields` slice. After any input event (including bool toggles that may change visibility), the cursor SHALL be clamped to `[0, len(visible)-1]`.

#### Scenario: Cursor clamp after visibility change
- **WHEN** the user is on the last visible field and toggles a bool that hides fields below
- **THEN** the cursor SHALL be clamped so it does not exceed the new visible field count

#### Scenario: Cursor re-evaluated after toggle
- **WHEN** the user toggles a bool field (space key)
- **THEN** the form SHALL re-evaluate `VisibleFields()` and clamp the cursor before processing further input

### Requirement: FormModel View renders description (MODIFIED)
The form View SHALL render the `Description` of the currently focused field below that field's input widget, styled with `tui.FieldDescStyle`.

#### Scenario: Focused field description displayed
- **WHEN** the form View is rendered and field at cursor has a non-empty Description
- **THEN** the view SHALL include a line with the description text below that field

#### Scenario: No description for unfocused fields
- **WHEN** a field is not focused
- **THEN** its Description SHALL not be rendered in the View output

### Requirement: Embedding ProviderID deprecation in state update (MODIFIED)
The `UpdateConfigFromForm` case for `emb_provider_id` SHALL set `cfg.Embedding.Provider` to the value AND clear `cfg.Embedding.ProviderID` to empty string.

#### Scenario: emb_provider_id clears deprecated field
- **WHEN** UpdateConfigFromForm processes key "emb_provider_id" with value "openai"
- **THEN** `cfg.Embedding.Provider` SHALL be "openai" AND `cfg.Embedding.ProviderID` SHALL be ""
