## Why

The `.placeholder` directory under `skills/` exists solely to satisfy the `go:embed **/SKILL.md` pattern at build time. However, `EnsureDefaults()` deploys it to `~/.lango/skills/.placeholder/` and `ListActive()` then attempts to parse it, producing a WARN log on every app startup: `skip invalid skill {"dir": ".placeholder", "error": "missing frontmatter delimiter (---)"}`. The spec (`skill-system/spec.md`) already states "placeholder SHALL NOT be deployed as a usable skill", but the code does not enforce this.

## What Changes

- **`ListActive()`**: Skip directories whose name starts with `.` (hidden directories), preventing `.placeholder` from being parsed.
- **`EnsureDefaults()`**: Skip embedded skill paths whose directory name starts with `.`, preventing `.placeholder/SKILL.md` from being deployed to the user's skills directory.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `skill-system`: Add requirement that hidden directories (names starting with `.`) are excluded from listing and deployment.

## Impact

- **Code**: `internal/skill/file_store.go` — `ListActive()` and `EnsureDefaults()` methods.
- **Behavior**: The spurious WARN log on every startup is eliminated. No user-visible skills are affected since hidden directories are not valid skill names.
