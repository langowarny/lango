# Security Hardening Tasks

## P0 - Critical

- [x] T1: Change `wrapWithApproval` from fail-open to fail-closed
- [x] T2: Add `crypto_decrypt` to sensitive tools pattern (approval required via config)

## P1 - High

- [x] T3: Implement Proxy/Reference Token pattern for secrets (`secrets_get` returns `{{secret:name}}`)
- [x] T4: Implement Proxy/Reference Token pattern for crypto decrypt (`crypto_decrypt` returns `{{decrypt:id}}`)
- [x] T5: Add reference token resolution in exec tool (`resolveRefs` before execution)
- [x] T6: Implement output-side secret scanning (`SecretScanner`)
- [x] T7: Extend PIIRedactingModelAdapter for output scanning (tool results + model responses)

## Testing

- [x] T8: Tests for fail-closed approval (wrapWithApproval behavior change)
- [x] T9: Tests for reference token system (RefStore: 11 tests, all pass with -race)
- [x] T10: Tests for output secret scanning (SecretScanner: 9 tests, all pass with -race)
- [x] T11: Updated secrets tool tests (proxy pattern verified)
- [x] T12: Updated crypto tool tests (proxy pattern verified)
