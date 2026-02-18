## ADDED Requirements

### Requirement: Cron Scheduler settings menu category
The Settings TUI editor SHALL include a "Cron Scheduler" menu category that opens a form for configuring cron scheduling settings.

#### Scenario: Cron menu appears in settings
- **WHEN** the user opens the Settings editor
- **THEN** a "Cron Scheduler" category with description "Scheduled jobs, timezone, history" SHALL appear in the menu after "Payment"

#### Scenario: Cron form displays all fields
- **WHEN** the user selects the "Cron Scheduler" menu item
- **THEN** the form SHALL display fields for: Enabled (bool), Timezone (text), Max Concurrent Jobs (int), Session Mode (select: isolated/main), History Retention (text)

### Requirement: Background Tasks settings menu category
The Settings TUI editor SHALL include a "Background Tasks" menu category that opens a form for configuring background task settings.

#### Scenario: Background menu appears in settings
- **WHEN** the user opens the Settings editor
- **THEN** a "Background Tasks" category with description "Async tasks, concurrency limits" SHALL appear in the menu after "Cron Scheduler"

#### Scenario: Background form displays all fields
- **WHEN** the user selects the "Background Tasks" menu item
- **THEN** the form SHALL display fields for: Enabled (bool), Yield Time in ms (int), Max Concurrent Tasks (int)

### Requirement: Workflow Engine settings menu category
The Settings TUI editor SHALL include a "Workflow Engine" menu category that opens a form for configuring workflow engine settings.

#### Scenario: Workflow menu appears in settings
- **WHEN** the user opens the Settings editor
- **THEN** a "Workflow Engine" category with description "DAG workflows, timeouts, state" SHALL appear in the menu after "Background Tasks"

#### Scenario: Workflow form displays all fields
- **WHEN** the user selects the "Workflow Engine" menu item
- **THEN** the form SHALL display fields for: Enabled (bool), Max Concurrent Steps (int), Default Timeout (text/duration), State Directory (text)
