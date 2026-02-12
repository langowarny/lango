## MODIFIED Requirements

### Requirement: Application Bootstrap
The application SHALL initialize components in the following order: config → supervisor → session store → tools → ADK agent → channels → gateway. Security initialization (crypto provider, passphrase) SHALL be skipped entirely when `security.signer.provider` is not configured. The application MUST NOT fail to start due to missing security configuration.

#### Scenario: Startup Sequence
- **WHEN** the application starts with valid provider and channel configuration
- **THEN** it SHALL initialize supervisor, session store, tools (exec + filesystem), ADK agent, channels, and gateway in order
- **THEN** it SHALL log "security disabled, set security.signer.provider to enable" when no signer is configured

#### Scenario: Startup Without Security
- **WHEN** the application starts with no `security.signer.provider` configured
- **THEN** it SHALL skip all security initialization (passphrase, crypto provider, secrets tool, crypto tool)
- **THEN** the agent SHALL still start and respond to messages through channels

### Requirement: Component Wiring
The application SHALL be organized into focused files: `app.go` (lifecycle orchestration, Start/Stop), `wiring.go` (component initialization: supervisor, session store, agent, gateway), `tools.go` (tool definitions and registration for exec and filesystem), `channels.go` (channel initialization and message handlers), `types.go` (type definitions). No single file SHALL exceed 200 lines.

#### Scenario: Tool Registration
- **WHEN** the application initializes tools
- **THEN** it SHALL register only exec tools (exec, exec_bg, exec_status, exec_stop) and filesystem tools (fs_read, fs_list, fs_write, fs_edit, fs_mkdir, fs_delete)
- **THEN** it SHALL NOT register browser, crypto, or secrets tools

#### Scenario: Agent Injection into Channels
- **WHEN** channels are initialized
- **THEN** each channel handler SHALL route messages through `adk.Agent.Run()` as the sole execution path

## REMOVED Requirements

### Requirement: Browser Tool Registration
**Reason**: Browser tool is removed from MVP to eliminate go-rod dependency and reduce startup complexity.
**Migration**: Browser tool source retained in `internal/tools/browser/`. Re-enable in Phase 2 by adding registration back to tools.go.

### Requirement: Crypto Tool Registration
**Reason**: Crypto tool is removed from MVP. Security tools are optional.
**Migration**: Crypto tool source retained in `internal/tools/crypto/`. Re-enable when security signer is configured.

### Requirement: Secrets Tool Registration
**Reason**: Secrets tool is removed from MVP. Security tools are optional.
**Migration**: Secrets tool source retained in `internal/tools/secrets/`. Re-enable when security signer is configured.
