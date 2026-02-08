# Lango 🚀

A high-performance AI agent built with Go, supporting multiple channels (Telegram, Discord, Slack) and advanced tools.

## Features

- 🔥 **Fast** - Single binary, <100ms startup, <100MB memory
- 🔒 **Secure** - No Node.js dependencies, reduced attack surface
- 🔌 **Multi-Channel** - Telegram, Discord, Slack support
- 🛠️ **Rich Tools** - Shell execution, file system, browser automation
- 💾 **Persistent** - SQLite session storage
- 🌐 **Gateway** - WebSocket/HTTP server for control plane

## Quick Start

### Installation

```bash
# Build from source
git clone https://github.com/langowarny/lango.git
cd lango
make build

# Or install directly
go install github.com/langowarny/lango/cmd/lango@latest
```

### Configuration

Create `lango.json`:

```json
{
  "server": {
    "host": "localhost",
    "port": 18789
  },
  "agent": {
    "provider": "gemini",
    "model": "gemini-2.0-flash-exp",
    "apiKey": "${GOOGLE_API_KEY}",
    "maxConversationTurns": 20
  },
  "channels": {
    "telegram": {
      "enabled": true,
      "botToken": "${TELEGRAM_BOT_TOKEN}"
    }
  },
  "logging": {
    "level": "info",
    "format": "console"
  }
}
```

### Run

```bash
# Start the server (ensure GOOGLE_API_KEY is set)
export GOOGLE_API_KEY=your_key_here
lango serve

# Or with custom config
lango serve --config /path/to/lango.json

# Validate configuration
lango config validate
```

### Getting Started

Use the interactive onboard wizard for first-time setup:

```bash
lango onboard
```

This guides you through:
1. API key configuration
2. Model selection
3. Channel setup (Telegram, Discord, or Slack)

### Diagnostics

Run the doctor command to check your setup:

```bash
# Check configuration and environment
lango doctor

# Auto-fix common issues
lango doctor --fix

# JSON output for scripting
lango doctor --json
```

## Architecture

```
lango/
├── cmd/lango/          # CLI entry point
├── internal/
│   ├── agent/          # Agent runtime with multi-provider support
│   ├── provider/       # AI Providers (OpenAI, Anthropic, Gemini)
│   ├── gateway/        # WebSocket/HTTP server
│   ├── channels/       # Telegram, Discord, Slack
│   ├── tools/          # exec, filesystem, browser
│   ├── session/        # SQLite session store
│   ├── config/         # Configuration loading
│   ├── logging/        # Zap logger
└── pkg/                # Public packages
```

## AI Providers

Lango supports multiple AI providers with a unified interface and automatic fallback.

### Supported Providers
- **OpenAI** (`openai`): GPT-4, GPT-3.5, and compatible APIs (Ollama, Groq, etc.)
- **Anthropic** (`anthropic`): Claude 3, Claude 3.5
- **Gemini** (`gemini`): Google Gemini models

### Configuration Example

```json
{
  "agent": {
    "provider": "openai",
    "model": "gpt-4o",
    "fallbackProvider": "anthropic",
    "fallbackModel": "claude-3-5-sonnet-20241022"
  },
  "providers": {
    "openai": {
      "apiKey": "${OPENAI_API_KEY}"
    },
    "anthropic": {
      "apiKey": "${ANTHROPIC_API_KEY}"
    },
    "ollama": {
        "baseUrl": "http://localhost:11434"
    }
  },
  "security": {
    "interceptor": {
      "enabled": true,
      "redactPii": true
    },
    "signer": {
      "provider": "local"
    }
  }
}
```

### Onboarding TUI
Use `lango onboard` to interactively configure providers, models, and security settings. The TUI allows you to manage multiple providers and set up local encryption.

## Configuration Reference

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `server.host` | string | `localhost` | Bind address |
| `server.port` | int | `18789` | Listen port |
| `agent.provider` | string | `gemini` | Primary AI provider ID |
| `agent.model` | string | `gemini-2.0-flash-exp` | Primary model ID |
| `agent.fallbackProvider` | string | - | Fallback provider ID |
| `agent.fallbackModel` | string | - | Fallback model ID |
| `agent.apiKey` | string | - | legacy: use `providers` section |
| `providers.<id>.type` | string | - | Provider type (openai, anthropic, gemini) |
| `providers.<id>.apiKey` | string | - | Provider API key |
| `providers.<id>.baseUrl` | string | - | Custom base URL (e.g. for Ollama) |
| `agent.maxTokens` | int | `4096` | Max tokens |
| `agent.maxConversationTurns` | int | `20` | Max conversation history turns |
| `logging.level` | string | `info` | Log level |
| `logging.format` | string | `console` | `json` or `console` |
| `session.databasePath` | string | `~/.lango/sessions.db` | SQLite path |
| `security.signer.provider` | string | `local` | `local` or `rpc` |
| `security.passphrase` | string | - | **DEPRECATED** Use `LANGO_PASSPHRASE` |

## Security

Lango includes built-in security features for AI agents:

### Security Configuration

Lango supports two security modes:

1. **Local Mode** (Default)
   - Encrypts secrets using AES-256-GCM derived from a passphrase (PBKDF2).
   - **Interactive**: Prompts for passphrase on startup (Recommended).
   - **Headless**: Set `LANGO_PASSPHRASE` environment variable.
   - **Migration**: Rotate your passphrase using:
     ```bash
     lango security migrate-passphrase
     ```
   > **⚠️ Warning**: Losing your passphrase results in permanent loss of all encrypted secrets. Lango does not store your passphrase.

2. **RPC Mode** (Production)
   - Offloads cryptographic operations to a hardware-backed companion app or external signer.
   - Keys never leave the secure hardware.

Configure mode in `lango.json`:

```json
{
  "security": {
    "signer": {
      "provider": "local" // or "rpc"
    }
  }
}
```

### Secrets Management

AI agents can securely store and retrieve secrets (API keys, tokens, etc.):

```go
// Stored secrets are encrypted with AES-256-GCM
secrets.store(name: "api-key", value: "sk-...")
secrets.get(name: "api-key")  // Requires user approval
secrets.list()
secrets.delete(name: "api-key")
```

### Cryptographic Operations

```go
crypto.encrypt(data: "sensitive", keyId: "default")
crypto.decrypt(ciphertext: "...", keyId: "default")
crypto.sign(data: "message")
crypto.hash(data: "content", algorithm: "sha256")
crypto.keys()  // List available keys
```

### Companion App Integration (RPC Mode)

Lango supports optional iOS/macOS companion apps for hardware-backed security:

- **mDNS Discovery** - Auto-discovers companion apps on the local network
- **Secure Enclave** - Keys never leave the hardware security module
- **Approval UI** - Native push notifications for sensitive operations

Configure in `lango.json`:

```json
{
  "security": {
    "signer": {
      "provider": "rpc"
    },
    "companion": {
      "enabled": true,
      "address": "ws://192.168.1.100:18790"
    }
  }
}
```

## Development

```bash
# Run tests
make test

# Run linter
make lint

# Build for all platforms
make build-all

# Run locally
make dev
```

## License

MIT
