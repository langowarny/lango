## Purpose

Define the CLI commands for managing workflow execution (run, list, status, cancel, history).

## Requirements

### Requirement: Workflow run command
The CLI SHALL provide `lango workflow run <file.yaml>` that parses and executes a workflow YAML file.

#### Scenario: Run a workflow
- **WHEN** user runs `lango workflow run code-review.flow.yaml`
- **THEN** the CLI SHALL parse the YAML, validate the DAG, and execute the workflow

#### Scenario: Run with schedule registration
- **WHEN** user runs `lango workflow run report.flow.yaml --schedule "0 9 * * MON"`
- **THEN** the CLI SHALL register the workflow with the cron scheduler (not yet implemented, logged as info)

#### Scenario: Invalid YAML file
- **WHEN** user runs `lango workflow run invalid.yaml` with malformed content
- **THEN** the CLI SHALL display a parse error

### Requirement: Workflow list command
The CLI SHALL provide `lango workflow list` that displays workflow runs with columns: Run ID, Workflow, Status, Steps, Started, Completed.

#### Scenario: List workflow runs
- **WHEN** user runs `lango workflow list`
- **THEN** the CLI SHALL display all workflow runs in tabular format

### Requirement: Workflow status command
The CLI SHALL provide `lango workflow status <run-id>` that displays detailed run information including all step statuses.

#### Scenario: View workflow status
- **WHEN** user runs `lango workflow status <uuid>`
- **THEN** the CLI SHALL display the run overview and a table of step statuses

### Requirement: Workflow cancel command
The CLI SHALL provide `lango workflow cancel <run-id>` that cancels a running workflow.

#### Scenario: Cancel a running workflow
- **WHEN** user runs `lango workflow cancel <uuid>`
- **THEN** the CLI SHALL cancel the workflow and display confirmation

### Requirement: Workflow history command
The CLI SHALL provide `lango workflow history` that displays completed workflow runs.

#### Scenario: View workflow history
- **WHEN** user runs `lango workflow history`
- **THEN** the CLI SHALL display recent workflow runs ordered by start time
