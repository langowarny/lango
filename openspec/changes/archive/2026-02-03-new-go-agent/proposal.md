## Why

OpenClaw is a powerful personal AI assistant, but being based on TypeScript/Node.js, it suffers from issues such as 50+ npm dependencies, high memory usage (~100-300MB), slow startup times (1-3 seconds), and npm supply chain security risks. Rewriting it in Go enables single binary distribution, fast startup (~50ms), low memory usage (~10-50MB), and the construction of a safer and faster agent utilizing modern agent frameworks like Google ADK-Go.

## What Changes

- **New Project Creation**: Create a Go agent project named `lango`
- **ADK-Go Based Agent**: LLM agent runtime using the Google ADK-Go framework
- **Gateway Server**: WebSocket-based control plane (OpenClaw Gateway counterpart)
- **Channel Integration**: Support for Telegram, Discord, and Slack channels
- **Tool System**: Bash execution, file manipulation, and browser control tools
- **Session Management**: SQLite-based persistent session store
- **Configuration System**: JSON/YAML-based configuration loading
- **CLI**: cobra-based command-line interface
- **Platforms**: macOS and Linux support (cross-compilation)

## Capabilities

### New Capabilities

- `agent-runtime`: ADK-Go based LLM agent execution engine. Includes tool calling, streaming responses, and session management.
- `gateway-server`: WebSocket/HTTP based control plane. Handles client connections, RPC methods, and state broadcasting.
- `channel-telegram`: Telegram bot integration. Handles message sending/receiving, group support, and media processing.
- `channel-discord`: Discord bot integration. Supports messages and slash commands.
- `channel-slack`: Slack app integration. Bolt-style event handling.
- `tool-exec`: Shell command execution tool. Supports PTY, timeouts, and background execution.
- `tool-filesystem`: File read/write/edit tool. Handles binary and text files.
- `tool-browser`: Browser automation tool. rod-based CDP control, screenshots, and DOM manipulation.
- `session-store`: SQLite-based session persistence. Stores history, state, and metadata.
- `config-system`: Configuration loading and validation. Supports JSON/YAML and environment variable substitution.

### Modified Capabilities

None (New Project)

## Impact

- **New Codebase**: Create a new Go project at `/Users/juwonkim/GolandProjects/lango`
- **Dependencies**: ADK-Go, zap, chi, gorilla/websocket, rod, sqlite3, telegram-bot-api, discordgo, slack-go
- **Build**: Go 1.22+, CGO (sqlite), Cross-compilation (darwin/linux)
- **Deployment**: Single binary, Docker image, systemd/launchd service
