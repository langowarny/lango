## Context

The `lango onboard` command's `printNextSteps()` function was written during the plaintext JSON config era. It instructs users to set environment variables for API keys and generates a `.lango.env.example` file. Since the migration to encrypted SQLite profiles, API keys stored directly in profiles are already secure. The legacy messaging creates confusion and the example file exposes key names in plaintext.

## Goals / Non-Goals

**Goals:**
- Remove all environment variable export guidance from post-save output
- Remove `.lango.env.example` file generation
- Simplify next-steps messaging to essential actions only
- Update doctor's API key security check to not alarm users who store keys in encrypted profiles
- Update onboard Long description to reflect current config categories

**Non-Goals:**
- Changing the actual config storage mechanism
- Removing environment variable reference support from the config system (still valid for portability)
- Adding new post-save features

## Decisions

1. **Delete `generateEnvExample()` entirely** rather than conditionally skipping it. The function has no callers outside `printNextSteps()` and the concept is incompatible with encrypted storage.

2. **Change `printNextSteps()` signature to `_ *config.Config`** since the config parameter is no longer inspected (no channel-specific env var hints). Kept the parameter for interface stability.

3. **Doctor check uses "Inline API keys" wording** instead of "Plaintext API keys detected". The term "plaintext" implies insecurity, but keys in encrypted profiles are safe. "Inline" accurately describes keys that aren't using `${ENV_VAR}` references without implying a security problem.

4. **Doctor check Details field explains both options** — encrypted profiles are safe, `${ENV_VAR}` references are an alternative for portability. This is informational rather than prescriptive.

## Risks / Trade-offs

- [Users who relied on `.lango.env.example`] → Minimal risk; the file was a convenience artifact and users with encrypted profiles don't need it.
- [Doctor check severity remains `StatusWarn`] → Kept as warning since inline keys do have a portability trade-off. The message is now accurate rather than alarming.
