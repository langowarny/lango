## ADDED Requirements

### Requirement: Cron configuration
The config system SHALL support a `cron` section with fields: enabled (bool), timezone (string), maxConcurrentJobs (int), defaultSessionMode (string), historyRetention (duration string).

#### Scenario: Default cron config
- **WHEN** no cron config is specified
- **THEN** defaults SHALL be: enabled=false, timezone="UTC", maxConcurrentJobs=5, defaultSessionMode="isolated", historyRetention="720h"

### Requirement: Background configuration
The config system SHALL support a `background` section with fields: enabled (bool), yieldMs (int), maxConcurrentTasks (int).

#### Scenario: Default background config
- **WHEN** no background config is specified
- **THEN** defaults SHALL be: enabled=false, yieldMs=30000, maxConcurrentTasks=3

### Requirement: Workflow configuration
The config system SHALL support a `workflow` section with fields: enabled (bool), maxConcurrentSteps (int), defaultTimeout (duration string), stateDir (string).

#### Scenario: Default workflow config
- **WHEN** no workflow config is specified
- **THEN** defaults SHALL be: enabled=false, maxConcurrentSteps=4, defaultTimeout="10m", stateDir="~/.lango/workflows/"
