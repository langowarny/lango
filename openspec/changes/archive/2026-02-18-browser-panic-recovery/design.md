## Context

The `go-rod/rod` library panics instead of returning errors when the CDP/WebSocket connection to Chrome drops unexpectedly. In Docker sidecar deployments (`chromedp/headless-shell`), the Chrome container can restart or become temporarily unreachable, causing panics that crash the entire Lango process. Currently there is no panic recovery at any layer: not in the browser tool methods, not in the tool handler wrappers, and not in the WebSocket goroutines.

## Goals / Non-Goals

**Goals:**
- Prevent rod/CDP panics from crashing the Lango process
- Convert panics into structured errors that callers can handle
- Auto-reconnect to Chrome after a connection loss (single retry)
- Protect WebSocket goroutines from unhandled panics
- Ensure Chrome is healthy before Lango starts in sidecar deployment

**Non-Goals:**
- Circuit breaker or exponential backoff (over-engineering for current scale)
- Monitoring/alerting integration for panic events
- Browser pool or connection multiplexing

## Decisions

### Decision 1: Layered Defense (3 layers of panic recovery)

**Choice**: Recover panics at three independent layers.
**Rationale**: Defense in depth — if the innermost layer misses a panic (e.g., panic in a goroutine spawned by rod), the outer layers catch it. Each layer operates independently.

- **Layer 1 (Core)**: `safeRodCall`/`safeRodCallValue` in `browser.go` — wraps every direct rod API call
- **Layer 2 (Application)**: `wrapBrowserHandler` in `tools.go` — wraps the tool handler closure
- **Layer 3 (Infrastructure)**: `readPump`/`writePump`/`handleRPC` recovery in `server.go`

**Alternative considered**: Single recovery at handler level only. Rejected because rod may panic in internal goroutines that the handler-level defer cannot catch.

### Decision 2: Generic panic recovery helper using Go 1.25 generics

**Choice**: `safeRodCallValue[T any]` generic function for value-returning calls.
**Rationale**: Avoids type assertion boilerplate. Go 1.25.4 fully supports generics. The non-generic `safeRodCall` variant handles error-only returns.

### Decision 3: Single retry on ErrBrowserPanic

**Choice**: Close the browser connection and retry exactly once on `ErrBrowserPanic`.
**Rationale**: A single retry covers the common case (Chrome restarted, connection stale). Infinite retries risk resource exhaustion. The retry happens at two points: `SessionManager.EnsureSession()` and `wrapBrowserHandler`.

### Decision 4: Sentinel error type

**Choice**: `var ErrBrowserPanic = errors.New("browser panic recovered")` with `%w` wrapping.
**Rationale**: Allows `errors.Is()` matching at any layer. Callers can distinguish browser panics from normal errors and take appropriate action (e.g., reconnect).

### Decision 5: Log + return for panics (exception to handle-once rule)

**Choice**: Both log and return the error when recovering from a panic.
**Rationale**: Panics are exceptional events that warrant immediate visibility in logs. The error is also returned so callers can react. This is a deliberate exception to the normal "handle once" principle.

### Decision 6: RPC handler isolation

**Choice**: Extract handler invocation into `handleRPC()` with its own `defer recover()`.
**Rationale**: Prevents a single handler panic from tearing down the entire `readPump` goroutine and disconnecting the client. The client receives an error response instead.

## Risks / Trade-offs

- **[Risk]** Rod panics in internally-spawned goroutines cannot be caught by any defer in our code → **Mitigation**: Layer 3 (readPump/writePump recovery) catches the resulting goroutine crash at the connection level. The specific client disconnects gracefully rather than crashing the process.
- **[Risk]** Single retry may not be enough if Chrome takes longer to restart → **Mitigation**: Docker healthcheck with `start_period: 10s` ensures Chrome is ready before Lango starts. Runtime restarts are handled by the retry; if it still fails, the user gets a clear error.
- **[Trade-off]** Logging panic + returning error violates handle-once → **Accepted**: Panics are rare, exceptional events. The log entry provides immediate observability while the error enables programmatic handling.
