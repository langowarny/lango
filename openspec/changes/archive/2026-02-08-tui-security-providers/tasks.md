# TUI Security & Providers Tasks

## Security Expansion
-   [ ] Implement `NewInterceptorForm` fields (Enabled, RedactPII, Approval) in `forms_impl.go`
-   [ ] Implement `NewSignerForm` fields (Provider, RPCUrl, KeyID) in `forms_impl.go`
-   [ ] Add `Passphrase` field to `SecurityForm`
-   [ ] Update `state_update.go` to handle nested security struct updates

## Providers Management
-   [ ] Create `ProvidersList` model in `tui/providers_list.go` (new file or component)
-   [ ] Implement "Providers" menu item in `menu.go`
-   [ ] Implement `NewProviderForm` in `forms_impl.go` (ID, Type, APIKey, BaseURL)
-   [ ] Add logic to `wizard.go` to handle Provider List -> Provider Form navigation
-   [ ] Update `state_update.go` to save providers to `config.Providers` map

## Validation & Verification
-   [ ] Verify Security settings are saved to `lango.json` correctly
-   [ ] Verify Providers map is populated in `lango.json`
-   [ ] Manual test: Add a new Ollama provider and ensure it persists
