## Why

The `lango settings` TUI editor covers 21 configuration categories but lacks P2P networking (5 sub-domains) and advanced security (keyring, DB encryption, KMS) settings. Users must hand-edit encrypted config JSON to configure these features, breaking the consistent TUI experience.

## What Changes

- Add 8 new menu categories to the settings TUI under "P2P Network" and "Security" sections
- Add form builders for each category with field types matching their config counterparts (bool, text, int, select, password)
- Add ~53 new config write-back case entries in `UpdateConfigFromForm()`
- Extend the Security form's signer provider dropdown with `aws-kms`, `gcp-kms`, `azure-kv`, `pkcs11`
- Handle `*bool` config fields (BlockConversations, ReadOnlyRootfs) with `derefBool`/`boolPtr` helpers
- Use conditional field visibility (`VisibleWhen`) for container sandbox and KMS backend-specific fields

## Capabilities

### New Capabilities
- `settings-p2p`: TUI forms for P2P Network, ZKP, Pricing, Owner Protection, and Sandbox
- `settings-security-advanced`: TUI forms for Security Keyring, DB Encryption, and KMS

### Modified Capabilities
- `cli-settings`: Menu expanded from 21 to 29 categories; signer provider options extended

## Impact

- `internal/cli/settings/menu.go` -- 2 new sections (P2P Network, Security) with 8 categories
- `internal/cli/settings/forms_impl.go` -- 8 new form builders + `derefBool`/`formatKeyValueMap` helpers
- `internal/cli/settings/editor.go` -- 8 new `case` routes in `handleMenuSelection()`
- `internal/cli/tuicore/state_update.go` -- ~53 new case entries + `boolPtr` helper
- `internal/cli/settings/forms_impl_test.go` -- tests for all new forms and helpers
