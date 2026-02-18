## Why

Gateway-level OIDC authentication is fully implemented (`internal/gateway/auth.go`) and config types (`AuthConfig`, `OIDCProviderConfig`) are already defined, but there is no TUI/CLI for users to configure OIDC providers. The README claims "Configure OIDC providers via `lango onboard`" but the onboard wizard has no Auth menu, creating a documentation-code mismatch.

## What Changes

- Add a new "Auth" menu item to the onboard TUI wizard for OIDC provider management
- Create an OIDC provider list component (add/edit/delete) following the existing providers list pattern
- Create an OIDC provider form with fields: Provider Name, Issuer URL, Client ID, Client Secret, Redirect URL, Scopes
- Add state update logic to persist OIDC provider configuration changes
- Update README.md Authentication section to accurately reflect the TUI onboard menu

## Capabilities

### New Capabilities

### Modified Capabilities

- `cli-onboard`: Add Auth/OIDC provider management menu to the onboard wizard TUI
- `oidc-auth`: No spec-level requirement change, only UI wiring (no delta spec needed)

## Impact

- `internal/cli/onboard/auth_providers_list.go` — new file for OIDC provider list UI
- `internal/cli/onboard/forms_impl.go` — new `NewOIDCProviderForm()` function
- `internal/cli/onboard/state_update.go` — new `UpdateAuthProviderFromForm()` function
- `internal/cli/onboard/menu.go` — new Auth category in menu
- `internal/cli/onboard/wizard.go` — new step, routing, ESC handling, view for auth flow
- `README.md` — Authentication section text update
