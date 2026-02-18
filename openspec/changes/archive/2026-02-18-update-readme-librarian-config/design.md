## Context

README.md Configuration Reference and Multi-Agent tables are out of date. Recent features (proactive librarian system, automation default delivery channels) have been implemented but not documented. The source of truth is `internal/config/types.go` (structs) and `internal/config/loader.go` (defaults).

## Goals / Non-Goals

**Goals:**
- Bring README.md Configuration Reference into parity with `LibrarianConfig`, `CronConfig`, `BackgroundConfig`, and `WorkflowConfig` structs
- Update multi-agent librarian description to reflect proactive knowledge extraction tools
- Add proactive librarian mention to the features list

**Non-Goals:**
- No code changes
- No new spec requirements — purely documentation alignment
- No restructuring of README sections

## Decisions

1. **Table placement**: Librarian section placed after Workflow Engine (follows config struct ordering in `types.go`)
2. **Default values**: Sourced directly from `loader.go` default initialization, not assumed
3. **Tool names**: Taken from actual tool registration (`librarian_pending_inquiries`, `librarian_dismiss_inquiry`)

## Risks / Trade-offs

- [Low] Documentation may drift again as new features are added → Mitigated by OpenSpec workflow archiving changes
