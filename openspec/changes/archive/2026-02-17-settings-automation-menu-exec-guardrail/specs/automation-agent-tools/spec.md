## ADDED Requirements

### Requirement: Exec prohibition in automation prompt
The automation prompt section SHALL include an explicit instruction prohibiting the use of exec to run lango automation CLI commands.

#### Scenario: Prompt includes exec prohibition
- **WHEN** any automation feature (cron, background, or workflow) is enabled
- **THEN** the automation prompt section SHALL contain text instructing the agent to NEVER use exec to run "lango cron", "lango bg", or "lango workflow" commands, with explanation that spawning a new lango process requires passphrase authentication
