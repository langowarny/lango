## Why

The project has grown from a basic CLI tool to a large framework with 36+ internal packages, Docker multi-profile support, Ent ORM, and multi-agent architecture, but the Makefile remains in its initial state. Developer experience and local CI pipeline capabilities need to match the project's maturity.

## What Changes

- Update build variables to explicitly set `CGO_ENABLED=1` (required by sqlite3/sqlite-vec)
- Add 16 new Make targets: `install`, `run`, `coverage`, `fmt`, `fmt-check`, `vet`, `ci`, `docker-build-browser`, `docker-push`, `docker-up`, `docker-up-browser`, `docker-up-sidecar`, `docker-down`, `docker-logs`, `health`
- Improve existing targets: `docker-build` adds `latest` tag, `lint` auto-installs golangci-lint, `clean` removes `.coverage/`, `help` auto-generated from comments, `.PHONY` lists all targets
- Add `bin/` and `.coverage/` to `.gitignore`

## Capabilities

### New Capabilities

_None — this change is build tooling only, no new runtime capabilities._

### Modified Capabilities

- `docker-deployment`: Makefile now provides compose orchestration targets (`docker-up`, `docker-down`, `docker-logs`) and browser-variant build target

## Impact

- `Makefile` — full rewrite
- `.gitignore` — two entries added (`bin/`, `.coverage/`)
- No runtime code changes, no API changes, no dependency changes
