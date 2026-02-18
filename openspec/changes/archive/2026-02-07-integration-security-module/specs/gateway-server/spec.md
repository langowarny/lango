## ADDED Requirements

### Requirement: Companion WebSocket Handler
The system SHALL handle WebSocket connections from companion apps separately from regular client connections.

#### Scenario: Companion connection on /companion endpoint
- **WHEN** a companion app connects to /companion WebSocket endpoint
- **THEN** the system SHALL validate the connection via mTLS
- **AND** mark the connection as a companion type

#### Scenario: Wire RPCProvider sender
- **WHEN** companion connection is established
- **THEN** the system SHALL configure RPCProvider.SetSender to send messages via the companion WebSocket

### Requirement: Companion Message Routing
The system SHALL route crypto RPC responses from companion to RPCProvider.

#### Scenario: Sign response routing
- **WHEN** companion sends a message with event type "sign.response"
- **THEN** the system SHALL call RPCProvider.HandleSignResponse

#### Scenario: Encrypt response routing
- **WHEN** companion sends a message with event type "encrypt.response"
- **THEN** the system SHALL call RPCProvider.HandleEncryptResponse

#### Scenario: Decrypt response routing
- **WHEN** companion sends a message with event type "decrypt.response"
- **THEN** the system SHALL call RPCProvider.HandleDecryptResponse

### Requirement: Approval Request Broadcasting
The system SHALL send approval requests to connected companions.

#### Scenario: Broadcast approval request
- **WHEN** ApprovalMiddleware requires user approval
- **THEN** the system SHALL send approval.request to all connected companions
- **AND** wait for approval.response from any companion
