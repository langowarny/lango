## 1. OIDC Provider List Component

- [x] 1.1 Create `internal/cli/onboard/auth_providers_list.go` with `AuthProviderItem` and `AuthProvidersListModel` types
- [x] 1.2 Implement `NewAuthProvidersListModel()` to populate list from `config.Auth.Providers`
- [x] 1.3 Implement `Update()` with key bindings: up/down navigation, enter select, d delete, esc back
- [x] 1.4 Implement `View()` with "Manage OIDC Providers" title, provider list (ID + IssuerURL), and "+ Add New OIDC Provider" option

## 2. OIDC Provider Form

- [x] 2.1 Add `NewOIDCProviderForm()` to `internal/cli/onboard/forms_impl.go`
- [x] 2.2 Include fields: oidc_id (new only), oidc_issuer, oidc_client_id (password), oidc_client_secret (password), oidc_redirect, oidc_scopes
- [x] 2.3 Set form title to "Add New OIDC Provider" or "Edit OIDC Provider: <id>"

## 3. State Update Logic

- [x] 3.1 Add `UpdateAuthProviderFromForm()` to `internal/cli/onboard/state_update.go`
- [x] 3.2 Initialize `Auth.Providers` map if nil, extract ID from form if empty
- [x] 3.3 Map form fields to `OIDCProviderConfig` fields with comma-split for scopes
- [x] 3.4 Call `MarkDirty("auth")` after update

## 4. Menu Integration

- [x] 4.1 Add `{"auth", "🔑 Auth", "OIDC provider configuration"}` to `NewMenuModel()` categories after Security

## 5. Wizard Flow Integration

- [x] 5.1 Add `StepAuthProvidersList` constant to `WizardStep` enum
- [x] 5.2 Add `authProvidersList` and `activeAuthProviderID` fields to `Wizard` struct
- [x] 5.3 Add `case "auth"` to `handleMenuSelection()` initializing auth providers list
- [x] 5.4 Add `StepAuthProvidersList` case to `Update()` routing with delete/select/exit handling
- [x] 5.5 Add ESC handling for `StepAuthProvidersList` (return to menu) and auth provider form (save and return to list)
- [x] 5.6 Add `isAuthProviderForm()` helper and update `isProviderForm()` to exclude OIDC titles
- [x] 5.7 Add `StepAuthProvidersList` case to `View()`

## 6. README Update

- [x] 6.1 Update Authentication section text to reference `lango onboard` > Auth menu and `lango config import`

## 7. Verification

- [x] 7.1 Run `go build ./...` and confirm no errors
- [x] 7.2 Run `go test ./internal/cli/onboard/...` and confirm all tests pass
- [x] 7.3 Run `go vet ./...` and confirm no issues
