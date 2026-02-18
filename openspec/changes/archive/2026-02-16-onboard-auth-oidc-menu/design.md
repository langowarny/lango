## Context

The gateway-level OIDC authentication is fully implemented in `internal/gateway/auth.go` with config types `AuthConfig` and `OIDCProviderConfig` defined in `internal/config/types.go`. However, the onboard TUI wizard has no UI for configuring OIDC providers. The README states users can configure OIDC via `lango onboard` but no such menu exists, creating a documentation-code gap.

The onboard wizard already has an established pattern for list-based management (providers list) and form-based editing that this change follows.

## Goals / Non-Goals

**Goals:**
- Add Auth menu to onboard TUI for OIDC provider CRUD operations
- Follow existing patterns: `providers_list.go` for list, `forms_impl.go` for form, `state_update.go` for persistence
- Fix README documentation to accurately reflect the TUI capability

**Non-Goals:**
- Changing the OIDC auth middleware or gateway auth logic
- Adding OIDC provider validation (e.g., testing issuer URL reachability)
- Modifying `config.OIDCProviderConfig` struct

## Decisions

1. **Pattern replication over abstraction**: Create `AuthProvidersListModel` as a separate type mirroring `ProvidersListModel` rather than generalizing. Rationale: the two lists display different fields (Type vs IssuerURL) and keeping them separate is simpler with no DRY violation risk at two instances.

2. **Form field key prefix `oidc_`**: All OIDC form fields use `oidc_` prefix to avoid collision with existing provider form keys (`type`, `id`, `apikey`). This ensures `isAuthProviderForm()` and `isProviderForm()` heuristics remain unambiguous.

3. **Title-based form type detection**: Extend existing heuristic (`isProviderForm` checks for "Provider" in title) by adding `isAuthProviderForm` checking for "OIDC" in title. Updated `isProviderForm` to exclude "OIDC" titles to prevent false positives.

4. **Menu placement**: Auth menu placed after Security, before Knowledge. This groups security-adjacent concerns together.

## Risks / Trade-offs

- [Heuristic form detection] Title-based detection (`isProviderForm`, `isAuthProviderForm`) is fragile if future form titles change. → Mitigation: Both methods are adjacent in code and documented; a more robust approach (e.g., enum field on Wizard) can be adopted if more list-backed forms are added.
- [No input validation on OIDC fields] Issuer URL and redirect URL are not validated for format. → Mitigation: Matches existing provider form behavior; validation can be added in a follow-up.
