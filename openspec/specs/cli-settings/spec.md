# CLI Settings Spec

## Goal
The `lango settings` command provides a comprehensive, interactive menu-based configuration editor for all aspects of the encrypted configuration profile.

## Requirements

### Configuration Coverage
The settings editor SHALL support editing all configuration sections previously handled by `lango onboard`:
1. **Providers** — Add, edit, delete multi-provider configurations
2. **Agent** — Provider, Model, MaxTokens, Temperature, PromptsDir, Fallback
3. **Server** — Host, Port, HTTP/WebSocket toggles
4. **Channels** — Telegram, Discord, Slack enable/disable + tokens
5. **Tools** — Exec timeout, Browser, Filesystem limits
6. **Session** — Database path, TTL, Max history turns
7. **Security** — Interceptor (PII, policy, timeout, tools), Signer (provider, RPC, KeyID)
8. **Auth** — OIDC provider management (add, edit, delete)
9. **Knowledge** — Enabled, max learnings/knowledge/context
10. **Skill** — Enabled, skills directory
11. **Observational Memory** — Enabled, provider, model, thresholds, budget
12. **Embedding & RAG** — Provider, model, dimensions, local URL, RAG settings
13. **Graph Store** — Enabled, backend, DB path, traversal depth, expansion results
14. **Multi-Agent** — Orchestration toggle
15. **A2A Protocol** — Enabled, base URL, agent name/description
16. **Payment** — Wallet, chain ID, RPC URL, USDC contract, limits, X402
17. **Cron Scheduler** — Enabled, timezone, max concurrent jobs, session mode, history retention
18. **Background Tasks** — Enabled, yield time, max concurrent tasks
19. **Workflow Engine** — Enabled, max concurrent steps, default timeout, state directory

### User Interface
- Menu-based navigation with 21 categories (19 sections + Save & Exit + Cancel)
- Free navigation between categories
- Uses shared `tuicore.FormModel` for all forms
- Provider and OIDC provider list views for managing collections

#### Scenario: Menu categories
- **WHEN** user launches `lango settings`
- **THEN** the menu SHALL display categories in order: Providers, Agent, Server, Channels, Tools, Session, Security, Auth, Knowledge, Skill, Observational Memory, Embedding & RAG, Graph Store, Multi-Agent, A2A Protocol, Payment, Cron Scheduler, Background Tasks, Workflow Engine, Save & Exit, Cancel

### Skill configuration form
The settings editor SHALL provide a Skill configuration form with the following fields:
- **Enabled** (`skill_enabled`) — Boolean toggle for enabling the file-based skill system
- **Skills Directory** (`skill_dir`) — Text input for the directory path containing SKILL.md files

#### Scenario: Edit skill settings
- **WHEN** user selects "Skill" from the settings menu
- **THEN** the editor SHALL display a form with Enabled toggle and Skills Directory text field pre-populated from `config.Skill`

#### Scenario: Save skill settings
- **WHEN** user edits skill fields and navigates back (Esc)
- **THEN** the changes SHALL be applied to `config.Skill.Enabled` and `config.Skill.SkillsDir`

### Cron Scheduler configuration form
The settings editor SHALL provide a Cron Scheduler configuration form with the following fields:
- **Enabled** (`cron_enabled`) — Boolean toggle
- **Timezone** (`cron_timezone`) — Text input for timezone (e.g., "UTC", "Asia/Seoul")
- **Max Concurrent Jobs** (`cron_max_jobs`) — Integer input
- **Session Mode** (`cron_session_mode`) — Select: isolated, main
- **History Retention** (`cron_history_retention`) — Text input for retention duration

#### Scenario: Edit cron settings
- **WHEN** user selects "Cron Scheduler" from the settings menu
- **THEN** the editor SHALL display a form with all cron fields pre-populated from `config.Cron`

### Background Tasks configuration form
The settings editor SHALL provide a Background Tasks configuration form with the following fields:
- **Enabled** (`bg_enabled`) — Boolean toggle
- **Yield Time (ms)** (`bg_yield_ms`) — Integer input
- **Max Concurrent Tasks** (`bg_max_tasks`) — Integer input

#### Scenario: Edit background settings
- **WHEN** user selects "Background Tasks" from the settings menu
- **THEN** the editor SHALL display a form with all background fields pre-populated from `config.Background`

### Workflow Engine configuration form
The settings editor SHALL provide a Workflow Engine configuration form with the following fields:
- **Enabled** (`wf_enabled`) — Boolean toggle
- **Max Concurrent Steps** (`wf_max_steps`) — Integer input
- **Default Timeout** (`wf_timeout`) — Text input for duration (e.g., "10m")
- **State Directory** (`wf_state_dir`) — Text input for directory path

#### Scenario: Edit workflow settings
- **WHEN** user selects "Workflow Engine" from the settings menu
- **THEN** the editor SHALL display a form with all workflow fields pre-populated from `config.Workflow`

### Command
```
lango settings [--profile <name>]
```
- Default profile: "default"
- Loads existing profile or creates new with defaults
- Saves via `configstore.Store.Save()` to encrypted profile

#### Scenario: Launch settings
- **WHEN** user runs `lango settings`
- **THEN** the editor SHALL display a welcome screen followed by the configuration menu

#### Scenario: Save from settings
- **WHEN** user selects "Save & Exit" from the menu
- **THEN** the configuration SHALL be saved as an encrypted profile
