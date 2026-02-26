## MODIFIED Requirements

### Requirement: Exec prohibition in automation prompt
The automation prompt section SHALL include an explicit instruction prohibiting the use of exec to run ANY lango CLI command, not only automation subcommands. The prohibition SHALL list all known subcommands (cron, bg, workflow, graph, memory, p2p, security, payment, config, doctor, and others) and explain that every lango CLI invocation requires passphrase authentication during bootstrap and will fail in non-interactive subprocess contexts.

#### Scenario: Prompt includes comprehensive exec prohibition
- **WHEN** any automation feature (cron, background, or workflow) is enabled
- **THEN** the automation prompt section SHALL contain text instructing the agent to NEVER use exec to run ANY "lango" CLI command, covering all subcommands including but not limited to cron, bg, workflow, graph, memory, p2p, security, payment, config, and doctor
- **AND** the prohibition SHALL explain that spawning a new lango process requires passphrase authentication and will fail in non-interactive mode
- **AND** the prohibition SHALL instruct the agent to ask the user to run commands directly in their terminal when no built-in tool equivalent exists

## ADDED Requirements

### Requirement: Comprehensive CLI exec guard
The `blockLangoExec()` function SHALL block ALL `lango` CLI invocations attempted through `exec` or `exec_bg` tools, using a two-phase approach: (1) specific subcommand guards with per-command tool alternative messages, and (2) a catch-all guard for any remaining `lango` prefix.

#### Scenario: Block subcommand with in-process equivalent
- **WHEN** the agent attempts to exec a `lango` subcommand that has in-process tool equivalents (graph, memory, p2p, security, payment, cron, bg, workflow)
- **THEN** the system SHALL return a blocked message listing the specific built-in tools to use instead

#### Scenario: Block subcommand without in-process equivalent
- **WHEN** the agent attempts to exec a `lango` subcommand that has no in-process equivalent (config, doctor, settings, serve, onboard, agent)
- **THEN** the system SHALL return a blocked message explaining that passphrase authentication is required and the user should run the command directly in their terminal

#### Scenario: Allow non-lango commands
- **WHEN** the agent attempts to exec a command that does not start with `lango ` or equal `lango`
- **THEN** the system SHALL allow the command to proceed (return empty string)

#### Scenario: Case-insensitive matching
- **WHEN** the agent attempts to exec a lango command in any case (e.g., `LANGO SECURITY DB-MIGRATE`)
- **THEN** the system SHALL still block and return the appropriate guidance message

### Requirement: Exec tool prompt safety rules
The `TOOL_USAGE.md` prompt SHALL include an explicit top-level rule under the Exec Tool section warning against using exec to run any `lango` CLI command. The rule SHALL list specific subcommands as examples and explain the passphrase failure mechanism. The rule SHALL also instruct the agent to inform the user and ask them to run commands directly when no built-in tool equivalent exists.

#### Scenario: TOOL_USAGE.md contains exec safety rule
- **WHEN** the agent's tool usage prompt is loaded
- **THEN** the first bullet point under "### Exec Tool" SHALL warn against running any lango CLI command via exec
