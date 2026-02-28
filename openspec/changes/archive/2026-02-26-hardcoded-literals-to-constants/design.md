## Context

The P2P, payment, and security subsystems grew rapidly with hardcoded string and integer literals scattered across ~25 files. These literals represented protocol statuses, firewall actions, chain IDs, proving schemes, KMS provider names, and display constants. The lack of typed constants meant typos compiled silently, refactoring required grep-and-pray, and the codebase diverged from Go style conventions.

## Goals / Non-Goals

**Goals:**
- Replace all hardcoded literals with typed constants, enums, or sentinel errors
- Maintain full JSON wire-format compatibility (underlying `string` types)
- Improve compile-time safety by using typed enums for domain values
- Add `Valid()` methods to enum types for runtime validation
- Follow Go naming conventions: `ErrX` for sentinel errors, `TypeValue` for enum constants

**Non-Goals:**
- Changing any wire protocol or API behavior
- Adding new validation logic beyond the `Valid()` methods
- Refactoring business logic or control flow
- Modifying test behavior (tests should pass unchanged except for literal references)

## Decisions

1. **`type X string` enums over `iota` ints**: String-based enums serialize naturally to/from JSON without custom marshalers. Since all target values are already strings in the wire protocol, this preserves backward compatibility with zero migration cost.

2. **Sentinel errors over typed errors**: The error messages are static and callers primarily need `errors.Is()` matching, not field extraction. Sentinel errors with `errors.New` are simpler and sufficient. Wrapped with `%w` for context.

3. **Phased leaf-first ordering**: Wallet/payment constants defined first since they are imported by 10+ other packages (e.g., `wallet.CurrencyUSDC`). This avoided circular dependencies and allowed parallel work on independent packages.

4. **Type cast at boundaries**: Where config strings flow into typed parameters (e.g., `KMSProviderName(cfg.Security.Signer.Provider)`), explicit casts are used at the boundary rather than changing config struct types. This keeps config deserialization simple.

5. **Local paygate status constants in handler.go**: The protocol handler has its own `PayGateResult` struct separate from `paygate.Result`. Rather than importing paygate types into the protocol package, local unexported constants mirror the paygate status values, maintaining the package boundary.

## Risks / Trade-offs

- **Type casts at config boundaries** → Mitigated by `Valid()` methods available for runtime validation if needed later
- **Parallel agent edits** → Mitigated by assigning non-overlapping file sets to each agent, with integration build verification after merge
- **Breaking callers that compare `Response.Status` to raw strings** → No external callers exist; all comparisons are internal and were updated
