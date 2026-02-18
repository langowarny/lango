## MODIFIED Requirements

### Requirement: Automation prompt sections mention auto-detection
The automation prompt section SHALL inform the agent that delivery channel parameters are optional and will be auto-detected from the current session context when omitted.

#### Scenario: Cron prompt
- **WHEN** the automation prompt section is built with cron enabled
- **THEN** the cron section SHALL include a note that deliver_to is optional and auto-detected

#### Scenario: Background prompt
- **WHEN** the automation prompt section is built with background enabled
- **THEN** the background section SHALL include a note that channel is optional and auto-detected

#### Scenario: Workflow prompt
- **WHEN** the automation prompt section is built with workflow enabled
- **THEN** the workflow section SHALL include a note that deliver_to in YAML is optional and auto-detected
