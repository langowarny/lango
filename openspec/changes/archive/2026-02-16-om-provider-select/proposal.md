## Why

The Observational Memory form's Provider field uses free-text `InputText`, requiring users to manually type provider names. The Agent form already uses `buildProviderOptions(cfg)` with `InputSelect` for provider selection. This inconsistency creates a poor UX where users must remember exact provider names for OM configuration.

## What Changes

- Change `om_provider` field from `InputText` to `InputSelect` in the Observational Memory onboard form
- Reuse `buildProviderOptions(cfg)` to dynamically populate registered provider options
- Add empty string (`""`) as the first option to represent "use agent default"
- Add test coverage verifying the field type and options structure

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `cli-onboard`: The OM provider field changes from free-text input to a select dropdown populated from registered providers

## Impact

- `internal/cli/onboard/forms_impl.go`: OM form provider field type and options
- `internal/cli/onboard/forms_impl_test.go`: New test for OM form field types
- No API or dependency changes; purely UI/UX improvement within the onboard TUI
