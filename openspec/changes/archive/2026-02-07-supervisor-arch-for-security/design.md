## Context

Currently, the `lango` agent runs as a single monolithic process where configuration secrets (API keys, bot tokens) are loaded into the process environment variables. The Agent Runtime and its Tools (specifically `exec`) inherit this environment. This creates a security vulnerability where a compromised or hallucinating agent could expose secrets by inspecting its own environment or executing commands like `env`.

## Goals / Non-Goals

**Goals:**
- Isolate sensitive secrets (API Keys, Bot Tokens) from the Agent Runtime environment.
- Prevent the `exec` tool from accessing secrets via environment variables.
- Establish a "Supervisor" component that manages secrets and lifecycle.
- Refactor `internal/agent` to be unaware of API keys.

**Non-Goals:**
- Full process isolation (e.g., separate OS processes or containers) in this initial phase. We are implementing "In-Process Logical Separation" as a first step.
- Hardware-backed security enclaves (TEE).

## Decisions

### 1. In-Process Logical Separation
**Decision**: We will implement the Supervisor and Runtime as separate logical components within the same Go process, enforced by strict interface boundaries.
**Rationale**: This allows for incremental refactoring without the immediate operational complexity of managing multiple processes or RPC layers. It prepares the codebase for future out-of-process separation.

### 2. Provider Proxy Pattern
**Decision**: The Agent Runtime will no longer accept API Keys in its configuration. Instead, it will accept a `provider.Provider` interface. The implementation passed to it will be a `ProviderProxy` that forwards generation requests to the Supervisor, which holds the actual keys.
**Rationale**: This ensures the Runtime *never* possesses the API keys in memory in a way that is easily accessible to tool execution contexts (though memory dump attacks are still possible in-process, this mitigates accidental leakage via tools).

### 3. Stub Tools for Privileged Operations
**Decision**: Sensitive tools like `exec` will be replaced in the Runtime with "Stub" implementations. These stubs will forward the execution request to the Supervisor. The Supervisor will enforce policies (e.g., environment variable whitelisting) before executing the actual command.
**Rationale**: Centralizes security policy enforcement in the Trusted Zone (Supervisor) rather than the Sandboxed Zone (Runtime).

### 4. Bootstrapping Order
**Decision**: `internal/app` will be refactored to initialize the `Supervisor` first (loading config and secrets), and then asking the Supervisor to spawn/initialize the `Runtime`.
**Rationale**: Reflects the dependency relationship—Runtime depends on Supervisor.

## Risks / Trade-offs

- **Complexity**: The bootstrapping logic in `internal/app` will become more complex.
- **Refactoring Scope**: This requires changing the `agent.New` signature, which breaks `internal/app` and potentially tests.
- **In-Process Limitation**: Since it's still one process, a native code exploit or memory inspection could still potentially recover keys. This protects primarily against *accidental* leakage via standard tool usage (e.g., `exec "env"`).

## Migration Plan

1.  Create `internal/supervisor` package.
2.  Implement `ProviderProxy`.
3.  Refactor `internal/agent` to accept `provider.Provider`.
4.  Update `internal/app` to wire components using the new pattern.
5.  Verify `exec` tool no longer sees secret env vars.
