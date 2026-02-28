## Context

The `skills/.placeholder/SKILL.md` file exists to satisfy `go:embed **/SKILL.md` when no real skill files are present. Currently, `EnsureDefaults()` deploys this file to the user directory and `ListActive()` attempts to parse it, producing a WARN log on every startup.

## Goals / Non-Goals

**Goals:**
- Eliminate the spurious WARN log by filtering hidden directories (starting with `.`) in both `ListActive()` and `EnsureDefaults()`
- Align code behavior with the existing spec requirement ("placeholder SHALL NOT be deployed as a usable skill")

**Non-Goals:**
- Changing the `.placeholder` file itself or the embed pattern
- Adding a general-purpose directory filter mechanism

## Decisions

**Filter by hidden directory convention (`.` prefix) rather than hardcoding `.placeholder`**
- Rationale: The `.` prefix convention is a well-understood Unix pattern for hidden/internal entries. Filtering by prefix is forward-compatible — any future build-only artifacts can use the same convention without code changes. Hardcoding `.placeholder` would be brittle and not generalizable.

**Filter at both `ListActive()` and `EnsureDefaults()`**
- Rationale: Defense in depth. Even if one filter is bypassed (e.g., manual file copy), the other prevents the invalid entry from surfacing.

## Risks / Trade-offs

- [Risk: User creates a skill starting with `.`] → Mitigated: dot-prefixed directories are conventionally hidden/internal; no legitimate skill name should start with `.`.
