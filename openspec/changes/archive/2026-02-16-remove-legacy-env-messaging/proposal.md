## Why

The onboarding post-save flow still outputs legacy guidance telling users to set environment variables and generates a `.lango.env.example` file. This is a leftover from the plaintext JSON config era. The current architecture uses encrypted SQLite profiles, making direct API key storage safe and the env example file a security risk (exposes key names in plaintext).

## What Changes

- Remove environment variable export instructions from `printNextSteps()` in the onboard command
- Remove `.lango.env.example` file generation (`generateEnvExample()` function)
- Simplify post-save messaging to show only essential next steps (serve, doctor, profile management)
- Update `Long` description to include Auth and Session categories, fix Security category scope
- Update doctor's API key security check to recognize encrypted profile keys as safe rather than flagging them as "plaintext"

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `cli-onboard`: Remove legacy env var messaging and .lango.env.example generation from post-save flow; update Long description categories
- `apikey-security-check`: Change warning/pass messages to acknowledge encrypted profiles as a safe storage method

## Impact

- `internal/cli/onboard/onboard.go` — `printNextSteps()` simplified, `generateEnvExample()` deleted, `Long` description updated, unused imports removed
- `internal/cli/doctor/checks/apikey_security.go` — Warning and pass messages updated for accuracy
- No breaking API changes; purely UI/messaging improvements
