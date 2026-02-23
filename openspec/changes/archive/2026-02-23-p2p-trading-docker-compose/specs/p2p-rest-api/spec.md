## ADDED Requirements

### Requirement: P2P status endpoint
The gateway SHALL expose `GET /api/p2p/status` that returns the local node's peer ID, listen addresses, and connected peer count as JSON.

#### Scenario: Query P2P status when node is running
- **WHEN** a client sends `GET /api/p2p/status` to the gateway
- **THEN** the response SHALL be HTTP 200 with JSON containing `peerId` (string), `listenAddrs` (string array), and `connectedPeers` (integer)

### Requirement: P2P peers endpoint
The gateway SHALL expose `GET /api/p2p/peers` that returns a list of currently connected peers with their IDs and multiaddresses.

#### Scenario: Query connected peers
- **WHEN** a client sends `GET /api/p2p/peers` to the gateway
- **THEN** the response SHALL be HTTP 200 with JSON containing `peers` (array of objects with `peerId` and `addrs` fields) and `count` (integer)

### Requirement: P2P identity endpoint
The gateway SHALL expose `GET /api/p2p/identity` that returns the local DID string derived from the wallet.

#### Scenario: Query identity with wallet configured
- **WHEN** a client sends `GET /api/p2p/identity` and the identity provider is available
- **THEN** the response SHALL be HTTP 200 with JSON containing `did` (string starting with `did:lango:`) and `peerId` (string)

#### Scenario: Query identity without identity provider
- **WHEN** a client sends `GET /api/p2p/identity` and the identity provider is nil
- **THEN** the response SHALL be HTTP 200 with JSON containing `did` as null and `peerId` (string)

### Requirement: P2P routes registration
The P2P REST endpoints SHALL be registered on the gateway router only when P2P components are initialized (i.e., `p2pComponents` is non-nil).

#### Scenario: P2P disabled
- **WHEN** P2P is disabled in configuration
- **THEN** no `/api/p2p/*` routes SHALL be registered on the gateway
