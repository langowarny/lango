## Why

The `settings`, `doctor`, and `onboard` commands have undergone significant feature additions (28 categories with 6 group sections and `/` search in settings, 14 checks in doctor, GitHub provider and auto-fetch models in onboard), but their `--help` text still reflects the old state, providing inaccurate information to users.

## What Changes

- Replace the `settings` Long description to list all 6 group sections (Core, Communication, AI & Knowledge, Infrastructure, P2P Network, Security) with their 28 categories, and mention `/` keyword search
- Replace the `doctor` Long description to list all 14 checks and mention `--fix` / `--json` flags
- Replace the `onboard` Long description to reflect GitHub provider support, auto-fetched models, and approval policy in step descriptions

## Capabilities

### New Capabilities

- `cli-help-text`: Accurate and complete --help descriptions for settings, doctor, and onboard commands

### Modified Capabilities


## Impact

- `internal/cli/settings/settings.go` — Long description string replacement
- `internal/cli/doctor/doctor.go` — Long description string replacement
- `internal/cli/onboard/onboard.go` — Long description string replacement
- No API, dependency, or behavioral changes — documentation-only update
