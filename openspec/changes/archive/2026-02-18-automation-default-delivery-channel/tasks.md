## 1. Config Layer

- [x] 1.1 Add `DefaultDeliverTo []string` field to CronConfig, BackgroundConfig, and WorkflowConfig in `internal/config/types.go`
- [x] 1.2 Register viper defaults for all three `defaultDeliverTo` fields in `internal/config/loader.go`

## 2. Settings TUI

- [x] 2.1 Add "Default Deliver To" InputText field to NewCronForm, NewBackgroundForm, and NewWorkflowForm in `internal/cli/settings/forms_impl.go`
- [x] 2.2 Add case mappings for cron_default_deliver, bg_default_deliver, wf_default_deliver in `internal/cli/tuicore/state_update.go`

## 3. Tool Handlers — Auto-detect and Fallback

- [x] 3.1 Add `detectChannelFromContext` helper in `internal/app/tools.go`
- [x] 3.2 Update `buildCronTools` signature to accept `defaultDeliverTo []string` and add fallback logic in cron_add handler
- [x] 3.3 Update `buildBackgroundTools` signature to accept `defaultDeliverTo []string` and add fallback logic in bg_submit handler
- [x] 3.4 Update `buildWorkflowTools` signature to accept `defaultDeliverTo []string` and add fallback logic in workflow_run handler
- [x] 3.5 Update call sites in `internal/app/app.go` to pass config DefaultDeliverTo slices

## 4. Executor Warning Logs

- [x] 4.1 Add Warn-level log in cron executor when DeliverTo is empty (`internal/cron/executor.go`)
- [x] 4.2 Upgrade Debug to Warn in background notification when OriginChannel is empty (`internal/background/notification.go`)
- [x] 4.3 Add Warn-level log in workflow engine when DeliverTo is empty after successful completion (`internal/workflow/engine.go`)

## 5. Automation Prompts

- [x] 5.1 Update cron, background, and workflow prompt sections in `internal/app/wiring.go` to mention auto-detection

## 6. Verification

- [x] 6.1 Run `go build ./...` and `go test ./...` to verify all changes compile and pass
