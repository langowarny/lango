## 1. Settings Menu & Form Builders

- [x] 1.1 Add cron, background, workflow categories to menu.go (after payment, before save)
- [x] 1.2 Add handleMenuSelection cases for cron, background, workflow in editor.go
- [x] 1.3 Implement NewCronForm in forms_impl.go (5 fields: enabled, timezone, max_jobs, session_mode, history_retention)
- [x] 1.4 Implement NewBackgroundForm in forms_impl.go (3 fields: enabled, yield_ms, max_tasks)
- [x] 1.5 Implement NewWorkflowForm in forms_impl.go (4 fields: enabled, max_steps, timeout, state_dir)

## 2. Form-to-Config Mapping

- [x] 2.1 Add cron field cases to UpdateConfigFromForm in state_update.go (5 cases)
- [x] 2.2 Add background field cases to UpdateConfigFromForm in state_update.go (3 cases)
- [x] 2.3 Add workflow field cases to UpdateConfigFromForm in state_update.go (4 cases)

## 3. Exec Guardrail

- [x] 3.1 Implement blockLangoExec helper in tools.go (detect lango cron/bg/workflow commands)
- [x] 3.2 Update buildExecTools signature to accept automationAvailable map
- [x] 3.3 Wire blockLangoExec check into exec handler (before sv.ExecuteTool)
- [x] 3.4 Wire blockLangoExec check into exec_bg handler (before sv.StartBackground)
- [x] 3.5 Update buildTools signature and pass automationAvailable map
- [x] 3.6 Build automationAvailable map in app.go and pass to buildTools

## 4. Prompt Reinforcement

- [x] 4.1 Add exec-prohibition text to buildAutomationPromptSection in wiring.go

## 5. Verification

- [x] 5.1 Run go build ./... — confirm zero errors
- [x] 5.2 Run go test ./... — confirm all tests pass
