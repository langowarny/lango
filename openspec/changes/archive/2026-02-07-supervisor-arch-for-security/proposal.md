## Why

Currently, `lango` runs as a single process where the Agent Runtime and Tools share the same environment and memory space as the configuration secrets (API Keys, Bot Tokens). This poses a security risk: if the Agent "hallucinates" or is tricked into running a command like `env` via the `exec` tool, it could expose these secrets.

This change introduces a **Supervisor Architecture** to isolate sensitive secrets from the Agent Runtime.

## What Changes

We will separate the application into two logical components (initially in-process, later potentially out-of-process):

1.  **Supervisor (internal/supervisor)**:
    - Owns the `config.Config` and all secrets (API Keys, etc.).
    - Initializes the real AI Provider Clients.
    - Spawns and manages the Agent Runtime.
    - Acts as a proxy for sensitive operations.

2.  **Runtime (internal/agent)**:
    - **No longer** holds API Keys or sensitive config.
    - Uses a `ProviderProxy` to request AI generation from the Supervisor.
    - Uses "Stub Tools" for sensitive operations (like `exec`) which forward requests to the Supervisor for validation and execution.

## Capabilities

### New Capabilities
- `supervisor-architecture`: Defines the Supervisor role, secret management, and the Runtime boundary.

### Modified Capabilities
- `config-system`: Update requirements to specify that secrets SHALL NOT be available to the Runtime environment directly.

## Impact

- **Breaking Change**: `agent.New()` signature will change to accept a `Provider` interface instead of `Config` with API keys.
- **Breaking Change**: `exec` tool will now require a Supervisor connection (or stub implementation).
- **Refactor**: `internal/app` bootstrapping will be significantly rewritten to initialize Supervisor first.
