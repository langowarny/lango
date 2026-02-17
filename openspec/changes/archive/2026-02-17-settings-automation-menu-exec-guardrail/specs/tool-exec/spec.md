## ADDED Requirements

### Requirement: Block lango automation commands via exec
The exec and exec_bg tool handlers SHALL detect and block commands that attempt to invoke lango CLI automation subcommands (cron, bg, background, workflow).

#### Scenario: Block lango cron command
- **WHEN** an exec or exec_bg tool receives a command starting with "lango cron"
- **THEN** the tool SHALL return a structured response with blocked=true and a message guiding to use built-in cron tools instead, without executing the command

#### Scenario: Block lango bg command
- **WHEN** an exec or exec_bg tool receives a command starting with "lango bg" or "lango background"
- **THEN** the tool SHALL return a structured response with blocked=true and a message guiding to use built-in background tools instead

#### Scenario: Block lango workflow command
- **WHEN** an exec or exec_bg tool receives a command starting with "lango workflow"
- **THEN** the tool SHALL return a structured response with blocked=true and a message guiding to use built-in workflow tools instead

#### Scenario: Context-aware guidance when feature is disabled
- **WHEN** an exec tool blocks a lango automation command and the corresponding feature is not enabled
- **THEN** the guidance message SHALL instruct the user to enable the feature in Settings

#### Scenario: Allow non-lango commands
- **WHEN** an exec or exec_bg tool receives a command that does not start with "lango cron", "lango bg", "lango background", or "lango workflow"
- **THEN** the tool SHALL execute the command normally without blocking
