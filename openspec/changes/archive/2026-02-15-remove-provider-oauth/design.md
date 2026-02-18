## Context

Lango currently supports two provider authentication methods: API keys and OAuth tokens. The OAuth path includes a `lango login [provider]` CLI command that performs browser-based OAuth flow, stores tokens in `~/.lango/tokens/`, and refreshes them on demand. However, AI providers like Google and GitHub may ban accounts using OAuth for programmatic API access, making this approach risky.

The Gateway OIDC authentication system (for user login to Lango itself) is completely separate and uses `OIDCProviderConfig` in `AuthConfig`. This system is unaffected.

## Goals / Non-Goals

**Goals:**
- Remove all provider OAuth authentication code paths
- Remove OAuth-related fields from `ProviderConfig`
- Add a doctor check to warn when API keys are stored as plaintext
- Maintain backward compatibility for existing configs (unknown JSON fields are silently ignored)

**Non-Goals:**
- Removing Gateway OIDC authentication (separate system, separate config types)
- Removing `golang.org/x/oauth2` from `go.mod` (still used by gateway)
- Auto-migrating existing OAuth tokens or config fields
- Implementing a secrets manager integration

## Decisions

### Decision 1: Full removal vs deprecation warning
**Choice**: Full removal of OAuth code
**Rationale**: OAuth with AI providers poses account ban risk. A deprecation period would leave users vulnerable. Clean removal is safer and simpler.
**Alternative**: Keep code but warn — rejected because the risk is ongoing, not eventual.

### Decision 2: Warning log vs error for missing API key
**Choice**: Warning log (`logger.Warnw`) when a provider has no API key
**Rationale**: Some providers (e.g., Ollama) may not require an API key. Hard-failing would break valid use cases.

### Decision 3: Plaintext detection approach
**Choice**: Simple `${...}` prefix/suffix check in doctor
**Rationale**: Environment variable references follow the `${VAR_NAME}` pattern. Any key not matching this pattern is treated as plaintext. This is simple, deterministic, and covers the common case.
**Alternative**: Regex pattern matching for known key formats (sk-*, AIza*) — rejected as fragile and provider-specific.

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| Users relying on `lango login` | Medium | Breaking change documented in spec. Migration path: use `apiKey` with `${ENV_VAR}` |
| Existing config with `clientId`/`clientSecret` | Low | Go JSON unmarshal ignores unknown fields — no errors |
| False positive plaintext detection | Low | Only warns, does not block. Users can ignore if intentional |
| Gateway OIDC accidentally broken | High | Completely separate code path (`OIDCProviderConfig` in `AuthConfig`), no shared code |
