## Context

The `settings`, `doctor`, and `onboard` CLI commands have had significant feature additions but their `--help` Long descriptions still reflect the old state. This is a documentation-only change to cobra command Long fields — no behavioral or API changes.

## Goals / Non-Goals

**Goals:**
- Update `settings` help to list all 6 group sections with 28 categories and `/` search
- Update `doctor` help to list all 14 checks and mention `--fix`/`--json` flags
- Update `onboard` help to reflect GitHub provider, auto-fetch models, and approval policy

**Non-Goals:**
- No changes to command behavior, flags, or runtime logic
- No new commands or subcommands
- No changes to Short descriptions (they remain accurate)

## Decisions

1. **Group-based listing for settings** — Present categories organized by their 6 UI groups rather than a flat list, matching the actual TUI structure users will see.
2. **Ordered check list for doctor** — List all 14 checks in the same order as `AllChecks()` for consistency with runtime output.
3. **Step-level detail for onboard** — Update each step's description to match the actual provider list, model auto-fetch, and approval policy features.

## Risks / Trade-offs

- [Drift risk] Help text can become stale again as features are added → Mitigated by including help text updates as part of the standard "update downstream artifacts" rule in CLAUDE.md.
