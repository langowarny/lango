## Why

The `lango onboard` TUI wizard has several usability bugs and missing features: an index-out-of-range panic when navigating form fields, legacy values/fields that no longer match the spec (sessions.db path, dead DB Passphrase field), no way to delete providers, hardcoded provider options in the Agent form, and session settings mixed into the Security form.

## What Changes

- Fix panic when pressing up/shift+tab at the first form field (cursor overflow to `len(Fields)`)
- Change default session database path from `sessions.db` to `data.db` (matching openspec standard)
- Remove dead `DB Passphrase` field from Security form (passphrase is acquired via keyfile/terminal prompt, not stored in config)
- Add provider delete functionality (`d` key) in the Providers list view
- Make Agent form Provider/Fallback Provider dropdowns dynamically populated from registered providers
- Split Session settings (DB Path, TTL, Max History Turns) out of Security form into a dedicated Session menu category
- Improve Provider creation UX: move Type selector before ID field, rename ID label to "Provider Name"

## Capabilities

### New Capabilities

### Modified Capabilities
- `cli-onboard`: Fix panic, add provider delete, split session form, dynamic provider options, provider ID UX improvement

## Impact

- `internal/cli/onboard/form.go` — cursor bounds fix
- `internal/cli/onboard/forms_impl.go` — session form extraction, passphrase removal, dynamic provider options, provider ID reorder
- `internal/cli/onboard/providers_list.go` — delete key binding
- `internal/cli/onboard/wizard.go` — delete handling, session menu routing
- `internal/cli/onboard/menu.go` — session category addition
- `internal/config/loader.go` — default database path change
- `internal/cli/onboard/forms_impl_test.go` — test updates for restructured forms
