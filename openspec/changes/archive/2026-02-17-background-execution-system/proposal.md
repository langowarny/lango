## Why

Lango excels at AI intelligence (self-learning, Graph RAG, multi-agent) but lacks autonomous execution — it cannot act without a user prompt. Users need scheduled tasks ("summarize news every morning at 9am"), long-running background operations (auto-yield after 30s), and declarative multi-step workflows (DAG-based parallel execution). These three automation capabilities transform Lango from a reactive assistant into a proactive agent.

## What Changes

- Add cron scheduling system: robfig/cron/v3-based scheduler with persistent jobs (Ent ORM), isolated session execution, and multi-channel delivery (Telegram/Slack/Discord)
- Add background task execution: in-memory task manager with concurrency control, auto-yield for long-running agent turns, and completion notifications
- Add workflow engine: YAML-defined DAG workflows with topological sort, parallel step execution, Go template variable substitution (`{{step-id.result}}`), and Ent-backed state persistence for resume
- Add observational memory compaction: wire `buffer.Trigger()` after agent turns, add `CompactMessages()` to session store for automatic message cleanup after observation
- Add CLI commands for all three systems: `lango cron`, `lango bg`, `lango workflow`
- Add channel sender adapter for programmatic message delivery from cron/background/workflow

## Capabilities

### New Capabilities
- `cron-scheduling`: Persistent cron job scheduling with at/every/cron schedule types, isolated session execution, and multi-channel result delivery
- `background-execution`: In-memory background task manager with concurrency limiting, async agent execution, and completion notifications
- `workflow-engine`: DAG-based YAML workflow engine with parallel step execution, template variable substitution, and Ent-persisted state for resume
- `cli-cron-management`: CLI commands for cron job CRUD (add/list/delete/pause/resume/history)
- `cli-bg-management`: CLI commands for background task management (list/status/cancel/result)
- `cli-workflow-management`: CLI commands for workflow execution (run/list/status/cancel/history)

### Modified Capabilities
- `observational-memory`: Add message compaction after observation — CompactMessages() deletes observed messages and inserts summary
- `gateway-server`: Add turn completion callbacks for buffer triggers (OnTurnComplete)
- `config-system`: Add CronConfig, BackgroundConfig, WorkflowConfig sections

## Impact

- **New packages**: `internal/cron/`, `internal/background/`, `internal/workflow/`
- **New CLI packages**: `internal/cli/cron/`, `internal/cli/bg/`, `internal/cli/workflow/`
- **New Ent schemas**: CronJob, CronJobHistory, WorkflowRun, WorkflowStepRun
- **New dependency**: `github.com/robfig/cron/v3`
- **Modified files**: `internal/config/types.go`, `internal/config/loader.go`, `internal/app/app.go`, `internal/app/wiring.go`, `internal/app/types.go`, `internal/gateway/server.go`, `internal/memory/buffer.go`, `internal/session/ent_store.go`, `cmd/lango/main.go`
- **New adapter**: `internal/app/sender.go` (channelSender for programmatic delivery)
