## ADDED Requirements

### Requirement: Cron form-to-config mapping
The UpdateConfigFromForm function SHALL map cron form fields to CronConfig struct fields.

#### Scenario: Cron fields update config
- **WHEN** a form with cron fields (cron_enabled, cron_timezone, cron_max_jobs, cron_session_mode, cron_history_retention) is submitted
- **THEN** the corresponding CronConfig fields (Enabled, Timezone, MaxConcurrentJobs, DefaultSessionMode, HistoryRetention) SHALL be updated

### Requirement: Background form-to-config mapping
The UpdateConfigFromForm function SHALL map background form fields to BackgroundConfig struct fields.

#### Scenario: Background fields update config
- **WHEN** a form with background fields (bg_enabled, bg_yield_ms, bg_max_tasks) is submitted
- **THEN** the corresponding BackgroundConfig fields (Enabled, YieldMs, MaxConcurrentTasks) SHALL be updated

### Requirement: Workflow form-to-config mapping
The UpdateConfigFromForm function SHALL map workflow form fields to WorkflowConfig struct fields.

#### Scenario: Workflow fields update config
- **WHEN** a form with workflow fields (wf_enabled, wf_max_steps, wf_timeout, wf_state_dir) is submitted
- **THEN** the corresponding WorkflowConfig fields (Enabled, MaxConcurrentSteps, DefaultTimeout as parsed duration, StateDir) SHALL be updated
