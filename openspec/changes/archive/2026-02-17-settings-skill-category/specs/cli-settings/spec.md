## MODIFIED Requirements

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

#### Scenario: Menu categories
- **WHEN** user launches `lango settings`
- **THEN** the menu SHALL display categories in order: Providers, Agent, Server, Channels, Tools, Session, Security, Auth, Knowledge, Skill, Observational Memory, Embedding & RAG, Graph Store, Multi-Agent, A2A Protocol, Payment, Save & Exit, Cancel

## ADDED Requirements

### Requirement: Skill configuration form
The settings editor SHALL provide a Skill configuration form with the following fields:
- **Enabled** (`skill_enabled`) — Boolean toggle for enabling the file-based skill system
- **Skills Directory** (`skill_dir`) — Text input for the directory path containing SKILL.md files

#### Scenario: Edit skill settings
- **WHEN** user selects "Skill" from the settings menu
- **THEN** the editor SHALL display a form with Enabled toggle and Skills Directory text field pre-populated from `config.Skill`

#### Scenario: Save skill settings
- **WHEN** user edits skill fields and navigates back (Esc)
- **THEN** the changes SHALL be applied to `config.Skill.Enabled` and `config.Skill.SkillsDir`
