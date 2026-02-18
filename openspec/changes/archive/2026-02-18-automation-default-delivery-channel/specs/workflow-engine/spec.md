## MODIFIED Requirements

### Requirement: Workflow delivery channel resolution
The workflow_run tool handler SHALL inject delivery channels into the parsed Workflow when DeliverTo is empty, using the three-tier fallback chain: YAML deliver_to → session auto-detection → workflow.defaultDeliverTo config. The engine SHALL log a Warn-level message when a workflow completes with no delivery channel configured.

#### Scenario: YAML deliver_to specified
- **WHEN** a workflow YAML includes a non-empty deliver_to field
- **THEN** the system SHALL use the YAML-specified channels without fallback

#### Scenario: Auto-detect from Slack session
- **WHEN** workflow_run is called with a YAML lacking deliver_to AND the session key starts with "slack:"
- **THEN** the system SHALL inject ["slack"] into the workflow's DeliverTo

#### Scenario: Config default used
- **WHEN** workflow_run is called with a YAML lacking deliver_to AND session auto-detection returns empty AND workflow.defaultDeliverTo is configured
- **THEN** the system SHALL use the config default channels

#### Scenario: No delivery channel warning
- **WHEN** a workflow completes successfully with empty DeliverTo
- **THEN** the engine SHALL log a Warn-level message including the workflow name and a configuration hint
