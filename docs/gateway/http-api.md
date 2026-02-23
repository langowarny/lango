---
title: HTTP API
---

# HTTP API

Lango exposes an HTTP API for health monitoring, agent discovery, authentication, and chat interaction.

## Server Configuration

> **Settings:** `lango settings` → Server

```json
{
  "server": {
    "host": "localhost",
    "port": 18789,
    "httpEnabled": true
  }
}
```

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `server.host` | `string` | `localhost` | Bind address for the HTTP server |
| `server.port` | `int` | `18789` | Port number |
| `server.httpEnabled` | `bool` | `true` | Enable the HTTP server |

## Endpoints

### Health Check

```
GET /health
```

Returns the server health status. Use this for monitoring, load balancer probes, and Docker health checks.

```bash
curl http://localhost:18789/health
```

### Agent Card (A2A)

```
GET /.well-known/agent.json
```

Returns the agent's A2A agent card when the A2A protocol is enabled (`a2a.enabled: true`). This endpoint follows the [Agent-to-Agent protocol](../features/a2a-protocol.md) specification for remote agent discovery.

### Authentication

Authentication endpoints are available when OIDC is configured. See [Authentication](../security/authentication.md) for details on the OAuth login flow and token management.

### Chat

The main chat endpoint accepts user messages and returns agent responses. When WebSocket is enabled, responses are streamed in real time via WebSocket events alongside the standard HTTP response.

### P2P Network

When P2P networking is enabled (`p2p.enabled: true`), the gateway exposes read-only endpoints for querying the running node's state. These endpoints are public (no authentication required) and return only node metadata.

#### `GET /api/p2p/status`

Returns the local node's peer ID, listen addresses, and connected peer count.

```bash
curl http://localhost:18789/api/p2p/status
```

```json
{
  "peerId": "12D3KooW...",
  "listenAddrs": ["/ip4/0.0.0.0/tcp/9000"],
  "connectedPeers": 2,
  "mdnsEnabled": true
}
```

#### `GET /api/p2p/peers`

Returns the list of currently connected peers with their IDs and multiaddresses.

```bash
curl http://localhost:18789/api/p2p/peers
```

```json
{
  "peers": [
    {
      "peerId": "12D3KooW...",
      "addrs": ["/ip4/172.20.0.3/tcp/9002"]
    }
  ],
  "count": 1
}
```

#### `GET /api/p2p/identity`

Returns the local DID derived from the wallet and the libp2p peer ID.

```bash
curl http://localhost:18789/api/p2p/identity
```

```json
{
  "did": "did:lango:02abc...",
  "peerId": "12D3KooW..."
}
```

If no identity provider is configured, `did` is `null`.

!!! note
    These REST endpoints query the **running server's persistent P2P node**. The CLI commands (`lango p2p status`, etc.) create ephemeral nodes for one-off operations. For monitoring and automation, prefer the REST API.

## Related

- [WebSocket](websocket.md) -- Real-time streaming events
- [Authentication](../security/authentication.md) -- OIDC and OAuth configuration
- [A2A Protocol](../features/a2a-protocol.md) -- Agent-to-Agent discovery
