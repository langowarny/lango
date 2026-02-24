## Why

P1-5's subprocess isolation provides basic memory separation for P2P tool execution, but lacks filesystem, network, and syscall-level isolation. A malicious tool invoked by a remote peer could still access the host filesystem or make network calls. Container-based execution provides comprehensive isolation boundaries needed for untrusted remote tool invocations.

## What Changes

- Add `ContainerSandboxConfig` to `ToolIsolationConfig` with runtime, image, network mode, resource limits, and pool settings
- Create `ContainerRuntime` interface abstracting container execution environments
- Implement `DockerRuntime` using Docker Go SDK for full container-based isolation
- Implement `NativeRuntime` wrapping existing `SubprocessExecutor` as a `ContainerRuntime` fallback
- Implement `GVisorRuntime` stub (always unavailable, placeholder for future)
- Create `ContainerExecutor` that probes runtimes in priority order: Docker → gVisor → Native
- Create optional `ContainerPool` for pre-warmed container management
- Add protocol version field to `ExecutionRequest` for forward compatibility
- Update app wiring to use `ContainerExecutor` when container mode is enabled, with subprocess fallback
- Add CLI commands: `lango p2p sandbox status/test/cleanup`
- Create sandbox Docker image (`build/sandbox/Dockerfile`)
- Add `sandbox-image` Makefile target

## Capabilities

### New Capabilities
- `container-sandbox`: Container-based tool execution isolation with Docker, gVisor (stub), and native subprocess fallback runtimes
- `container-pool`: Optional pre-warmed container pool for reduced cold-start latency

### Modified Capabilities
- `tool-isolation`: Extended with container runtime selection and configuration

## Impact

- **Config**: `internal/config/types.go` (ToolIsolationConfig + ContainerSandboxConfig), `internal/config/loader.go` (defaults)
- **New files**: `internal/sandbox/container_runtime.go`, `docker_runtime.go`, `native_runtime.go`, `gvisor_runtime.go`, `container_executor.go`, `container_pool.go`
- **Modified**: `internal/sandbox/executor.go` (Version field), `internal/app/app.go` (wiring)
- **CLI**: `internal/cli/p2p/sandbox.go` (status/test/cleanup), `internal/cli/p2p/p2p.go` (registration)
- **Build**: `build/sandbox/Dockerfile`, `Makefile` (sandbox-image target)
- **Dependencies**: `github.com/docker/docker` (Docker Go SDK)
