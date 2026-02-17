## 1. Foundation & Config

- [x] 1.1 Add CronConfig, BackgroundConfig, WorkflowConfig to internal/config/types.go
- [x] 1.2 Add default values for new configs in internal/config/loader.go (DefaultConfig + viper defaults)
- [x] 1.3 Add robfig/cron/v3 dependency to go.mod

## 2. Ent Schemas

- [x] 2.1 Create CronJob Ent schema (internal/ent/schema/cron_job.go) with all required fields and indexes
- [x] 2.2 Create CronJobHistory Ent schema (internal/ent/schema/cron_job_history.go) with status enum and indexes
- [x] 2.3 Create WorkflowRun Ent schema (internal/ent/schema/workflow_run.go) with status enum and indexes
- [x] 2.4 Create WorkflowStepRun Ent schema (internal/ent/schema/workflow_step_run.go) with status enum and indexes
- [x] 2.5 Run go generate for Ent code generation

## 3. Cron Scheduling Package

- [x] 3.1 Create internal/cron/job.go — Job, JobResult, HistoryEntry domain types
- [x] 3.2 Create internal/cron/store.go — Store interface + EntStore CRUD implementation
- [x] 3.3 Create internal/cron/scheduler.go — Scheduler with Start/Stop/AddJob/RemoveJob/PauseJob/ResumeJob, semaphore concurrency
- [x] 3.4 Create internal/cron/executor.go — Executor with AgentRunner interface, isolated session key generation, history recording
- [x] 3.5 Create internal/cron/delivery.go — Delivery with ChannelSender interface for multi-channel result dispatch

## 4. Background Execution Package

- [x] 4.1 Create internal/background/task.go — Task struct with Status state machine and mutex-protected transitions
- [x] 4.2 Create internal/background/manager.go — Manager with Submit/Cancel/Status/List/Result and semaphore concurrency
- [x] 4.3 Create internal/background/notification.go — Notification with ChannelNotifier interface
- [x] 4.4 Create internal/background/monitor.go — Monitor with ActiveCount() and Summary()

## 5. Workflow Engine Package

- [x] 5.1 Create internal/workflow/step.go — Workflow, Step, RunResult, RunStatus types with YAML tags
- [x] 5.2 Create internal/workflow/parser.go — Parse/ParseFile/Validate with DFS cycle detection
- [x] 5.3 Create internal/workflow/dag.go — DAG with TopologicalSort returning parallel layers
- [x] 5.4 Create internal/workflow/template.go — RenderPrompt for {{step-id.result}} substitution
- [x] 5.5 Create internal/workflow/state.go — StateStore with Ent-backed WorkflowRun/WorkflowStepRun persistence
- [x] 5.6 Create internal/workflow/engine.go — Engine with Run/Resume/Cancel/Status/ListRuns and DAG execution

## 6. ObservationalMemory Compaction

- [x] 6.1 Add CompactMessages(key, upToIndex, summary) to internal/session/ent_store.go
- [x] 6.2 Add MessageCompactor type, SetCompactor(), and compaction logic to internal/memory/buffer.go
- [x] 6.3 Add TurnCallback type and OnTurnComplete() to internal/gateway/server.go
- [x] 6.4 Fire turn callbacks in handleChatMessage after agent turn

## 7. App Wiring

- [x] 7.1 Add CronScheduler, BackgroundManager, WorkflowEngine fields to internal/app/types.go
- [x] 7.2 Create internal/app/sender.go — channelSender adapter for Telegram/Discord/Slack
- [x] 7.3 Add agentRunnerAdapter, initCron(), initBackground(), initWorkflow() to internal/app/wiring.go
- [x] 7.4 Wire new components in internal/app/app.go New() (steps 11-15)
- [x] 7.5 Wire CronScheduler Start/Stop in app.go Start() and Stop()
- [x] 7.6 Wire memory compaction (SetCompactor + OnTurnComplete callbacks) in app.go

## 8. CLI Commands

- [x] 8.1 Create internal/cli/cron/cron.go — add/list/delete/pause/resume/history with resolveJobID helper
- [x] 8.2 Create internal/cli/bg/bg.go — list/status/cancel/result commands
- [x] 8.3 Create internal/cli/workflow/workflow.go — run/list/status/cancel/history commands
- [x] 8.4 Register cron and workflow CLI commands in cmd/lango/main.go

## 9. Build Verification

- [x] 9.1 Run go build ./... and verify clean compilation
- [x] 9.2 Run go test ./... and verify all tests pass
- [x] 9.3 Run go mod tidy and verify clean
