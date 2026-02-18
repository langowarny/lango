## Why

Recent features (proactive librarian, automation default delivery channels) implemented in 2026-02-17~18 are not reflected in README.md. The Configuration Reference table is missing the `librarian.*` section and `*.defaultDeliverTo` fields, and the multi-agent table has stale librarian tool/role descriptions.

## What Changes

- Add `librarian.*` configuration section (7 fields) to the Configuration Reference table
- Add `cron.defaultDeliverTo`, `background.defaultDeliverTo`, `workflow.defaultDeliverTo` fields to their respective sections
- Update the multi-agent librarian row with proactive knowledge extraction role and new tools (`librarian_pending_inquiries`, `librarian_dismiss_inquiry`)
- Add "proactive knowledge librarian" to the Self-Learning feature bullet

## Capabilities

### New Capabilities

(none — documentation-only change)

### Modified Capabilities

(none — no spec-level requirement changes, only README documentation updates)

## Impact

- `README.md` — sole file modified
- No code, API, or dependency changes
