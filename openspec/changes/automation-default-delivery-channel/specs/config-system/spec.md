## MODIFIED Requirements

### Requirement: Automation config DefaultDeliverTo fields
CronConfig, BackgroundConfig, and WorkflowConfig SHALL each include a `DefaultDeliverTo []string` field with mapstructure tag "defaultDeliverTo". The config loader SHALL register viper defaults for all three fields.

#### Scenario: Default config values
- **WHEN** the application starts with no explicit defaultDeliverTo configuration
- **THEN** the DefaultDeliverTo fields SHALL default to nil (empty slice)

#### Scenario: Config file specifies defaults
- **WHEN** the config file sets cron.defaultDeliverTo to ["telegram"]
- **THEN** the loaded CronConfig.DefaultDeliverTo SHALL contain ["telegram"]
