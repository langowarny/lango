## Context

Lango is a Go-based AI agent platform with self-learning, Graph RAG, and multi-agent orchestration. Currently, all agent interactions are user-initiated: the agent only acts when a user sends a message. This limits Lango to reactive mode — it cannot perform scheduled tasks, handle long-running operations gracefully, or execute complex multi-step workflows autonomously.

The existing architecture provides foundational patterns that these new systems build on:
- **Buffer lifecycle** (Start/Trigger/Stop): Used by EmbeddingBuffer, GraphBuffer, memory.Buffer
- **Session isolation**: EntStore supports multiple concurrent sessions with key-based lookup
- **Channel adapters**: Telegram, Discord, Slack channels with typed Send methods
- **AgentRunner pattern**: ADK agent with RunAndCollect for single-turn execution

## Goals / Non-Goals

**Goals:**
- Enable scheduled agent execution via cron expressions, intervals, and one-time triggers
- Support background task execution with concurrency limiting and completion notifications
- Provide declarative YAML workflow definition with DAG-based parallel step execution
- Wire observational memory compaction for automatic message cleanup after observation
- Deliver results to configured channels (Telegram/Slack/Discord) programmatically
- Persist job/workflow state in Ent ORM for reliability across restarts

**Non-Goals:**
- Distributed scheduling across multiple Lango instances (single-instance only)
- Real-time progress streaming for background tasks (polling-based status only)
- Visual workflow editor or UI-based workflow creation
- Complex workflow primitives (loops, conditionals, sub-workflows) — only DAG with template variables
- Gateway auto-yield timer implementation (Manager API is ready, yield wiring deferred)

## Decisions

### 1. robfig/cron/v3 for Scheduling
**Choice**: Use `robfig/cron/v3` library for cron scheduling.
**Rationale**: Battle-tested Go cron library with timezone support, standard cron expression parsing, and `@every` duration syntax. Alternatives considered:
- Custom timer-based scheduler: More control but reinvents well-solved problems
- `go-co-op/gocron`: Feature-rich but heavier dependency for our needs
- OS-level crontab: Not portable, no programmatic control

### 2. AgentRunner Interface for Decoupling
**Choice**: Define `AgentRunner` interface in each consumer package (cron, background, workflow), implemented by `agentRunnerAdapter` in the app package.
**Rationale**: Avoids import cycles between `internal/app` and the new packages. Each package defines exactly the interface it needs (Interface Segregation). Alternative was a shared interface package, but that adds unnecessary coupling.

### 3. In-Memory Background Task Manager
**Choice**: Background tasks are managed in-memory (not persisted).
**Rationale**: Background tasks are ephemeral — they exist only during the current server lifecycle. Persisting them adds complexity without clear benefit since incomplete tasks would need re-evaluation on restart anyway. Cron jobs and workflows use Ent persistence because they represent durable scheduling intent.

### 4. channelSender Adapter for Delivery
**Choice**: Create a `channelSender` in `internal/app/` that dispatches to concrete channel adapters.
**Rationale**: Cron, background, and workflow packages need to send messages to channels but cannot import channel packages directly (would create circular deps). The adapter pattern isolates delivery concerns. Each consumer package defines its own sender interface (`ChannelSender` / `ChannelNotifier`).

### 5. DAG Execution with Topological Sort Layers
**Choice**: Topological sort returns layers of independent steps; each layer executes in parallel with semaphore-limited concurrency.
**Rationale**: Simple and correct. Layer-based execution naturally handles dependencies — all steps in a layer have their dependencies satisfied. Semaphore prevents resource exhaustion. Alternative was event-driven execution (notify on each step completion), but layers are simpler to reason about and debug.

### 6. Go Templates for Variable Substitution
**Choice**: `{{step-id.result}}` syntax using Go `text/template`.
**Rationale**: Go templates are stdlib, well-understood, and sufficient for our needs (simple variable replacement). Custom template syntax would add parsing complexity without benefit.

### 7. Ent ORM for State Persistence
**Choice**: Four new Ent schemas: CronJob, CronJobHistory, WorkflowRun, WorkflowStepRun.
**Rationale**: Consistent with existing Lango patterns (session, message, skill schemas all use Ent). Provides type-safe queries, automatic migrations, and SQLite compatibility.

## Risks / Trade-offs

- **[Single-instance scheduling]** → Cron jobs execute on the running instance only. If the instance restarts, missed schedules are not retroactively executed. Mitigation: Jobs persist in DB and reload on startup; NextRunAt tracking helps detect gaps.

- **[In-memory background tasks lost on restart]** → Running background tasks are lost if the server restarts. Mitigation: Acceptable trade-off since tasks are ephemeral. Users can re-trigger via cron or workflow.

- **[Workflow state consistency]** → If a step fails mid-execution, the workflow state may be partially updated. Mitigation: Each step's state is saved atomically; Resume() picks up from last completed step.

- **[Channel delivery failures]** → If a target channel is unavailable, delivery silently fails. Mitigation: Errors are logged; cron history records failure status. Future improvement: retry with backoff.

- **[Telegram chat ID requirement]** → Telegram delivery requires at least one allowlisted chat ID in config. Mitigation: Clear error message returned; Discord/Slack use default channels.
