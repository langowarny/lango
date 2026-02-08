## ADDED Requirements

### Requirement: Application Bootstrap
The system SHALL initialize all core components through a centralized application entry point (`internal/app`).

#### Scenario: Startup Sequence
- **WHEN** the application starts
- **THEN** it SHALL load configuration
- **THEN** it SHALL initialize the SQLite session store via Ent
- **THEN** it SHALL initialize the Agent Runtime with the session store and configured tools
- **THEN** it SHALL initialize Channels and Gateway, injecting the Agent
- **THEN** it SHALL start all background services (Gateway, Channels)

#### Scenario: Graceful Shutdown
- **WHEN** the application receives a termination signal (SIGINT/SIGTERM)
- **THEN** it SHALL stop the Gateway server
- **THEN** it SHALL stop all active Channels
- **THEN** it SHALL close the Database connection
- **THEN** it SHALL allow a grace period for active requests to complete

### Requirement: Component Wiring
The system SHALL inject dependencies between components to enable communication.

#### Scenario: Agent Injection into Gateway
- **WHEN** the Gateway is initialized
- **THEN** it SHALL receive a reference to the active Agent Runtime
- **AND** it SHALL use this reference to delegate "chat.message" RPC calls to the Agent

#### Scenario: Agent Injection into Channels
- **WHEN** a Channel (Telegram, Discord, Slack) is initialized
- **THEN** it SHALL receive a reference to the active Agent Runtime
- **AND** it SHALL execute the Agent for incoming messages that match the channel's criteria

#### Scenario: Tool Registration
- **WHEN** the Agent is initialized
- **THEN** the system SHALL create and register Tool instances (Browser, FS, Exec) based on `lango.json` configuration
