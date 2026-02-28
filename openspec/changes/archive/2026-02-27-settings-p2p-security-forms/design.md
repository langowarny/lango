## Context

The `lango settings` TUI editor follows a Bubble Tea pattern: menu categories -> form builders (`NewXForm()`) -> centralized config write-back (`UpdateConfigFromForm()` switch). All P2P and advanced security config types already exist in `internal/config/types.go` and are consumed by `internal/app/wiring.go`, but lacked TUI exposure.

## Goals / Non-Goals

**Goals:**
- Expose all P2P and advanced security settings through the existing TUI settings editor
- Follow established patterns for form builders, state update switch, and menu categories
- Use conditional field visibility for nested sub-sections (container sandbox, KMS backend-specific fields)
- Handle `*bool` config fields correctly with helper functions

**Non-Goals:**
- Changing config types or initialization logic
- Adding list management UIs for complex struct arrays (e.g., FirewallRules)
- Backend implementation of P2P or KMS features (already done)

## Decisions

### Split P2P into 5 sub-categories
P2PConfig has 6 nested sub-domains with 30+ fields. A single form would be unwieldy. Split into P2P Network (14 fields), ZKP (5), Pricing (3), Owner Protection (5), and Sandbox (11).

### Separate security sub-categories
The existing Security form already has 15 fields. Adding Keyring (1), DB Encryption (2), and KMS (12) would create a 30-field form. Separate menu entries keep each form focused.

### Conditional field visibility for container and KMS fields
Container sandbox fields (runtime, image, network, rootfs, CPU, pool) are only visible when Container Sandbox is enabled. KMS fields are conditionally visible based on selected backend (cloud vs PKCS#11 vs local).

### Expand signer provider options in existing Security form
Adding `aws-kms`, `gcp-kms`, `azure-kv`, `pkcs11` to the existing signer provider dropdown avoids a redundant form. The KMS form's backend selector mirrors this for consistency.

### Reuse existing type-mapping patterns
- `[]string` -> comma-separated text with `splitCSV()` (same as RAG Collections)
- `map[string]string` -> `key:value` comma-separated with `parseCustomPatterns()` (same as PII)
- `*bool` -> `derefBool(ptr, defaultVal)` / `boolPtr(val)` helpers (new, minimal)

## Risks / Trade-offs

- Menu length increases from 21 to 29 items -- acceptable with `/` search and j/k scrolling
- FirewallRules not editable in TUI -- complex struct arrays deferred to future work
- `*bool` is new to the form system -- contained to 2 fields, well-tested
