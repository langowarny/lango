## MODIFIED Requirements

### Requirement: Settings forms for default delivery channels
The Cron, Background, and Workflow settings forms SHALL each include a "Default Deliver To" text input field that accepts comma-separated channel names. The state update handler SHALL map these fields to the respective config DefaultDeliverTo slices using the splitCSV helper.

#### Scenario: Cron default deliver field
- **WHEN** the user opens the Cron Scheduler settings form
- **THEN** the form SHALL display a "Default Deliver To" field with placeholder "telegram,discord,slack (comma-separated)"

#### Scenario: Background default deliver field
- **WHEN** the user opens the Background Tasks settings form
- **THEN** the form SHALL display a "Default Deliver To" field with placeholder "telegram,discord,slack (comma-separated)"

#### Scenario: Workflow default deliver field
- **WHEN** the user opens the Workflow Engine settings form
- **THEN** the form SHALL display a "Default Deliver To" field with placeholder "telegram,discord,slack (comma-separated)"

#### Scenario: State update mapping
- **WHEN** the user enters "telegram,discord" in the cron default deliver field
- **THEN** the config state SHALL update Cron.DefaultDeliverTo to ["telegram", "discord"]
