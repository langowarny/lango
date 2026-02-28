## Context

macOS Data Protection Keychain with biometric ACL (`kSecAccessControlBiometryAny`) requires proper Apple Developer code signing with the `keychain-access-groups` entitlement. Ad-hoc signed binaries (produced by `go build`) cannot access this Keychain, resulting in OSStatus `-34018` (`errSecMissingEntitlement`). Both `BiometricProvider` and `OSProvider` use the same macOS Keychain with identical service/account keys (`lango` / `master-passphrase`), but `OSProvider` stores items without biometric ACL, so no entitlement is needed.

## Goals / Non-Goals

**Goals:**
- Detect `-34018` entitlement errors as a typed sentinel (`ErrEntitlement`) for programmatic matching
- Automatically fall back to plain OS Keychain when biometric storage fails due to missing entitlements
- Provide clear user messaging explaining why biometric storage is unavailable and how to fix it
- Support both read-path fallback (passphrase acquisition) and write-path fallback (bootstrap + CLI store)
- Add codesign infrastructure for release builds that need biometric protection

**Non-Goals:**
- Changing the security tier detection logic — `DetectSecureProvider` still returns biometric tier on macOS with Touch ID hardware
- Removing or weakening the security tier model — plain OS keyring is a graceful degradation, not a replacement
- Auto-detecting code signing status at startup

## Decisions

**1. Sentinel error via `errors.New` + `fmt.Errorf %w` wrapping**
- Rationale: Callers use `errors.Is(err, keyring.ErrEntitlement)` without type-asserting. Follows project error conventions (go-errors.md). Each call site wraps with its own context prefix (`keychain biometric get:`, etc.)
- Alternative: Custom error type with OSStatus field — rejected as over-engineered for a single status code.

**2. FallbackProvider as explicit Options field, not automatic chain**
- Rationale: Keeps the priority chain transparent: secure provider → fallback provider → keyfile → interactive → stdin. The caller (bootstrap) decides whether to wire a fallback, preserving the principle that TierNone systems should NOT auto-use plain OS keyring.
- Alternative: Chain internally in Acquire() by detecting `ErrEntitlement` — rejected because it would require Acquire to know about OS keyring construction, violating separation of concerns.

**3. macOS-only fallback guard (`runtime.GOOS == "darwin"`)**
- Rationale: Only macOS has the shared-Keychain property where `BiometricProvider` and `OSProvider` hit the same backend. On Linux, biometric and TPM are distinct backends, so fallback to OS keyring would not help.

**4. Entitlements plist with `$(AppIdentifierPrefix)` variable**
- Rationale: Uses Apple's build-time variable expansion so the entitlement works with any team ID. Standard pattern for Keychain access groups.

## Risks / Trade-offs

- [Reduced security in fallback] Plain OS keyring items are readable by any process running as the same UID → Mitigation: Clear warning messages; biometric protection available via `make codesign`. The fallback is strictly better than no persistence (which forces keyfile or repeated interactive prompts).
- [Same service/account key collision] If a user switches between ad-hoc and codesigned binaries, the Keychain item may exist with or without biometric ACL → Mitigation: `BiometricProvider` always deletes existing items before writing (`SecItemDelete` in `keychain_set_biometric`), so the ACL is reset correctly on each store.
- [Codesign requires Apple Developer identity] `make codesign` needs `APPLE_IDENTITY` → Mitigation: Self-signed certificates can work for local use; documented in the make target help text.
