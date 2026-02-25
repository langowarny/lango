## Context

P1-5 introduced `SubprocessExecutor` that runs tools in child processes with clean environments (only PATH/HOME). This provides memory isolation but no filesystem, network, or syscall restriction. For untrusted remote peer tool invocations, stronger isolation is needed.

Docker provides container-level isolation: separate filesystem namespace, configurable network mode ("none" for complete isolation), read-only rootfs, memory/CPU limits, and automatic cleanup.

## Goals / Non-Goals

**Goals:**
- Abstract container execution behind a `ContainerRuntime` interface for pluggable runtimes
- Implement Docker-based isolation with resource limits, network isolation, and read-only rootfs
- Automatic fallback to subprocess when container runtime unavailable
- Reuse existing JSON stdin/stdout protocol from P1-5
- Optional container pool for reduced cold-start latency
- CLI for status inspection, smoke testing, and orphan cleanup

**Non-Goals:**
- gVisor implementation (stub only, future work)
- Custom seccomp profiles (rely on Docker defaults)
- Image registry or automated image building (manual `make sandbox-image`)
- Windows container support

## Decisions

1. **`ContainerRuntime` interface** — `Run(ctx, ContainerConfig)`, `Cleanup(ctx, id)`, `IsAvailable(ctx)`, `Name()`. Clean abstraction allows Docker, gVisor, native fallback without changing callers.

2. **Docker SDK over CLI** — `github.com/docker/docker/client` provides type-safe API, proper error handling, and stream hijacking for stdin/stdout communication. No shell escaping risks.

3. **Runtime probe chain** — `NewContainerExecutor` probes: Docker (if available) → gVisor (if available) → Native (always available). Config `runtime` field can force specific runtime or "auto" for probe.

4. **Docker stream headers** — Docker multiplexes stdout/stderr with 8-byte headers per frame. The `stripDockerStreamHeaders` function handles this when raw JSON parsing fails.

5. **Ephemeral containers** — Each tool invocation creates a container, executes, collects result, and removes the container. Labels (`lango.sandbox=true`, `lango.tool=<name>`) enable orphan cleanup.

6. **Pool as optional optimization** — `ContainerPool` (channel-based, PoolSize > 0 to enable) pre-warms containers. Default disabled (PoolSize: 0) to avoid Docker daemon dependency at startup.

## Risks / Trade-offs

- **Docker daemon dependency**: Container mode requires Docker installed and running. Mitigated by automatic fallback to subprocess.
- **Cold start latency**: Container creation adds 100ms-2s per invocation. Pool helps but adds complexity.
- **Image availability**: `lango-sandbox:latest` must be built and available. `lango p2p sandbox test` validates this.
- **Resource overhead**: Each container adds ~10-50MB memory overhead. Controlled by pool size and idle timeout.
