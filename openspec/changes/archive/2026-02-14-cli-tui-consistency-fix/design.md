## Context

The Onboard TUI wizard (`internal/cli/onboard/`) uses a menu-driven form system: each config section (Agent, Server, Channels, Tools, Security, Providers) has a dedicated `New*Form` function in `forms_impl.go` that creates form fields, and `state_update.go` maps those fields back to `config.Config`. After recent backend additions (Browser tool recovery, Knowledge system, Agent fallback, Session max history), several config fields exist in `config.Config` but have no corresponding TUI form entries. Additionally, Doctor's SecurityCheck uses an if/else chain that doesn't recognize the new `enclave` signer provider.

## Goals / Non-Goals

**Goals:**
- Every user-configurable field in `config.Config` is editable through the Onboard TUI
- Doctor SecurityCheck recognizes all valid signer providers (`local`, `rpc`, `enclave`)
- Command Long descriptions accurately reflect current functionality
- No behavioral regressions in existing forms or checks

**Non-Goals:**
- Redesigning the form system or menu layout
- Adding validation logic beyond type parsing (e.g., path existence checks)
- Modifying config defaults or config file format
- Adding new Doctor checks beyond updating existing ones

## Decisions

### 1. Extend existing forms rather than restructure

**Decision**: Add fields directly to existing `New*Form` functions, add a new `NewKnowledgeForm` for the Knowledge section.

**Rationale**: The form system is simple and works. Each form maps 1:1 with a config section. Knowledge has no existing form, so it needs a new function, but the pattern is identical. Restructuring would be unnecessary scope creep.

### 2. Place Knowledge menu between Security and Providers

**Decision**: Insert `{"knowledge", "🧠 Knowledge", "Learning, Skills, Context limits"}` after Security in the menu.

**Rationale**: Knowledge is a feature-level config (not infra), and logically groups near Security (both are system-behavior settings). Providers is always last before Save/Cancel.

### 3. Switch statement for signer provider check

**Decision**: Replace if/else-if chain with a switch statement in `SecurityCheck.Run`.

**Rationale**: Switch is cleaner, handles `enclave` (no warnings—most secure), `rpc` (no warnings—production-ready), `local` (warn), and default (fail for unknown). Adding future providers is trivial.

### 4. Fallback provider uses InputSelect with empty option

**Decision**: `fallback_provider` uses `InputSelect` with options `["", "anthropic", "openai", "gemini", "ollama"]`.

**Rationale**: Empty string means "no fallback." This is consistent with how the config zero-value works—empty string disables the feature.

## Risks / Trade-offs

- **[Risk] Form field count grows large in Security form (now 10 fields)** → Acceptable for now; if it grows further, consider sub-sections or collapsible groups in future TUI redesign.
- **[Risk] Duration fields (`browser_session_timeout`, `ttl`) accept freeform text** → Go's `time.ParseDuration` handles validation at the `state_update.go` level; invalid input silently keeps the previous value. This matches existing behavior for `exec_timeout`.
