## REMOVED Requirements

### Requirement: Session database path field
**Reason**: The `db_path` field in the Session Configuration form displayed a path (`~/.lango/data.db`) that was never used in production. The bootstrap process opens `~/.lango/lango.db` and the session store reuses that client, making the field misleading.
**Migration**: No user action needed. The database path is determined automatically by bootstrap.

#### Scenario: Session form no longer shows database path
- **WHEN** user selects "Session" from the settings menu
- **THEN** the form SHALL display only TTL and Max History Turns fields
- **THEN** the form SHALL NOT display a Database Path field

## MODIFIED Requirements

### Configuration Coverage
The settings editor SHALL support editing all configuration sections previously handled by `lango onboard`:
1. **Providers** — Add, edit, delete multi-provider configurations
2. **Agent** — Provider, Model, MaxTokens, Temperature, PromptsDir, Fallback
3. **Server** — Host, Port, HTTP/WebSocket toggles
4. **Channels** — Telegram, Discord, Slack enable/disable + tokens
5. **Tools** — Exec timeout, Browser, Filesystem limits
6. **Session** — TTL, Max history turns
7. **Security** — Interceptor (PII, policy, timeout, tools), Signer (provider, RPC, KeyID)
8. **Auth** — OIDC provider management (add, edit, delete)
9. **Knowledge** — Enabled, max learnings/knowledge/context
10. **Skill** — Enabled, skills directory
11. **Observational Memory** — Enabled, provider, model, thresholds, budget, context limits
12. **Embedding & RAG** — Provider, model, dimensions, local URL, RAG settings
13. **Graph Store** — Enabled, backend, DB path, traversal depth, expansion results
14. **Multi-Agent** — Orchestration toggle
15. **A2A Protocol** — Enabled, base URL, agent name/description
16. **Payment** — Wallet, chain ID, RPC URL, USDC contract, limits, X402
17. **Cron Scheduler** — Enabled, timezone, max concurrent jobs, session mode, history retention
18. **Background Tasks** — Enabled, yield time, max concurrent tasks
19. **Workflow Engine** — Enabled, max concurrent steps, default timeout, state directory
20. **Librarian** — Enabled, observation threshold, inquiry cooldown, max inquiries, auto-save confidence, provider, model

#### Scenario: Menu categories
- **WHEN** user launches `lango settings`
- **THEN** the menu SHALL display categories in order: Providers, Agent, Server, Channels, Tools, Session, Security, Auth, Knowledge, Skill, Observational Memory, Embedding & RAG, Graph Store, Multi-Agent, A2A Protocol, Payment, Cron Scheduler, Background Tasks, Workflow Engine, Librarian, Save & Exit, Cancel
