## Context

Lango is a high-performance AI agent built with Go. The current CLI provides only basic commands (`serve`, `version`, `config validate`). Users must manually create `lango.json` configuration files and troubleshoot issues without tooling assistance.

OpenClaw provides rich TUI-based onboarding and diagnostic tools that significantly improve developer experience. This design brings similar capabilities to Lango while maintaining Go's simplicity.

## Goals / Non-Goals

**Goals:**
- Implement `lango doctor` command for configuration and connectivity diagnostics
- Implement `lango onboard` command for guided initial setup
- Use bubbletea for interactive TUI with consistent styling
- Support non-interactive mode with `--json` output for CI/scripting
- Keep the onboarding flow minimal: API key → Model → One channel

**Non-Goals:**
- WebUI or browser-based interface (explicitly deferred)
- Full OpenClaw feature parity (e.g., skills, hooks, canvas)
- Multi-agent routing or advanced session management
- Native app companions (macOS/iOS/Android)

## Decisions

### 1. TUI Library: bubbletea

**Decision**: Use github.com/charmbracelet/bubbletea with lipgloss for styling.

**Rationale**:
- Most popular Go TUI library with active maintenance
- Elm architecture provides clean, testable code
- lipgloss enables consistent cross-platform styling
- huh (built on bubbletea) provides form primitives

**Alternatives considered**:
- `survey`: Simpler but less visually appealing, limited customization
- `tview`: More traditional ncurses-style, steeper learning curve
- `promptui`: Too limited for wizard-style flows

### 2. Doctor Check Categories

**Decision**: Organize checks into modular categories:

```
internal/cli/doctor/
├── doctor.go          # Main command, orchestrates checks
├── checks/
│   ├── config.go      # Config file validation
│   ├── apikey.go      # API key verification
│   ├── channels.go    # Channel token validation
│   ├── database.go    # SQLite accessibility
│   ├── network.go     # Port availability, API connectivity
│   └── checks.go      # Check interface and result types
└── output/
    ├── tui.go         # Interactive TUI output
    └── json.go        # JSON output for --json flag
```

**Rationale**: Each check is isolated, testable, and can be extended independently.

### 3. Onboard Flow: Minimal Steps

**Decision**: Implement 4-step minimal onboarding:
1. Welcome + mode selection (QuickStart vs Advanced)
2. API key input (with validation)
3. Model selection (dropdown)
4. Single channel setup (choose one: Telegram/Discord/Slack)

**Rationale**: Lower barrier to entry. Advanced configuration via `lango configure` (future).

### 4. Configuration Repair with --fix

**Decision**: Doctor `--fix` flag attempts safe auto-repairs:
- Generate missing session database
- Create default config if missing
- Fix common config syntax issues

**Rationale**: Follows OpenClaw pattern. Risky operations require interactive confirmation.

### 5. Package Structure

**Decision**:
```
internal/cli/
├── tui/               # Shared TUI components (styles, widgets)
│   ├── styles.go      # lipgloss styles
│   └── components.go  # Reusable bubbletea components
├── doctor/            # Doctor command
└── onboard/           # Onboard command
```

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| bubbletea adds ~2MB to binary size | Acceptable trade-off for UX improvement |
| API key validation requires network | Allow offline mode, skip validation with warning |
| Channel token validation may rate-limit | Cache validation results, implement backoff |
| Non-TTY environments break TUI | Detect non-TTY, fall back to `--json` mode |

## Open Questions

1. Should `--fix` require explicit confirmation for each repair, or accept-all mode?
2. Should doctor check for available updates to lango binary?
