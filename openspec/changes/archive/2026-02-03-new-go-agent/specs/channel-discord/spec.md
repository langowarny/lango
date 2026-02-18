## ADDED Requirements

### Requirement: Discord bot connection
The system SHALL connect to Discord using the Bot API with a provided bot token.

#### Scenario: Successful bot connection
- **WHEN** the application starts with a valid DISCORD_BOT_TOKEN
- **THEN** the system SHALL establish a gateway connection to Discord

#### Scenario: Gateway reconnection
- **WHEN** the Discord gateway connection drops
- **THEN** the system SHALL automatically reconnect with session resumption

### Requirement: Message reception
The system SHALL receive and process incoming messages from Discord channels.

#### Scenario: Direct message received
- **WHEN** a user sends a DM to the bot
- **THEN** the message SHALL be forwarded to the agent

#### Scenario: Channel message with mention
- **WHEN** a user mentions the bot in a text channel
- **THEN** the message SHALL be processed if the channel is configured

### Requirement: Slash command support
The system SHALL register and handle Discord slash commands.

#### Scenario: Slash command registration
- **WHEN** the bot connects to Discord
- **THEN** configured slash commands SHALL be registered with the Discord API

#### Scenario: Slash command invocation
- **WHEN** a user invokes a slash command
- **THEN** the system SHALL execute the corresponding handler

### Requirement: Message sending
The system SHALL send agent responses back to Discord channels.

#### Scenario: Send text response
- **WHEN** the agent generates a response
- **THEN** the response SHALL be sent to the originating channel

#### Scenario: Discord markdown formatting
- **WHEN** a response contains markdown
- **THEN** the response SHALL use Discord-compatible markdown syntax

### Requirement: Guild configuration
The system SHALL support per-guild configuration for Discord servers.

#### Scenario: Guild-specific settings
- **WHEN** a guild has custom configuration
- **THEN** the bot SHALL use guild-specific settings for that server
