# Research: Local Runtime Foundation

## Decisions

### 1. Use a standard-library Go runtime with one small WebSocket dependency

**Decision**: Use Go's `net/http`, `context`, `os/exec`, `os/signal`, and `embed` for the native runtime. Use `github.com/coder/websocket` only for the WebSocket boundary.

**Rationale**: The standard library supplies the process lifecycle, HTTP, static assets, and controlled subprocess support required for this milestone. Coder WebSocket is maintained and context-aware. This avoids adding a CLI or server framework before either is needed.

**Alternatives considered**:

- Cobra or another CLI framework: rejected; the small command set does not justify another abstraction.
- A full HTTP framework: rejected; it adds little over `net/http` for a small local API.
- Polling: rejected; it cannot meet the ordered reconnect requirement.

### 2. Serve one explicit production origin on IPv4 loopback

**Decision**: Bind the production listener explicitly to `127.0.0.1:<port>` and serve the dashboard, API, and WebSocket endpoint from that one origin. The default is port 7331; when the implicit default is occupied, probe a bounded nearby range and report the selected address. An explicitly requested occupied port fails clearly. IPv6 dual-listener support is deferred.

**Rationale**: A numeric IPv4 loopback listener prevents accidental LAN exposure and avoids hostname resolution ambiguity. A single production origin eliminates a permissive CORS requirement.

**Alternatives considered**:

- `:port` or `0.0.0.0`: rejected because it can expose the runtime to the network.
- `localhost`: rejected as an address because it can resolve differently by host configuration.
- Separate UI/API production ports: rejected because it increases origin and token complexity.

### 3. Use a detached `serve` runtime with health-confirmed lifecycle control

**Decision**: The short-lived `start` command resolves the canonical Git root, launches an internal `serve` process for that root, waits for a health response that matches the generated runtime ID, and then attempts browser opening. `status`, `stop`, and `restart` validate the record against that health identity rather than trusting a process identifier. The long-lived runtime owns one signal-cancelled root context and all its goroutines.

**Rationale**: This preserves normal terminal usage, gives every concurrent component a clear owner, and makes stale-record recovery reliable. `http.Server.Shutdown` supports graceful stop rather than abrupt termination.

**Alternatives considered**:

- Keep the invoking CLI process resident: rejected because it occupies the developer's terminal and complicates status/stop commands.
- Trust a stored process ID alone: rejected because process IDs can be stale or reused.
- Hard-kill on normal shutdown: rejected because it risks lost events and corrupt snapshots.

### 4. Persist a redacted snapshot plus bounded JSONL event history

**Decision**: Define a storage interface and implement per-project atomic JSON snapshots/runtime records plus append-only JSONL event segments retaining 10,000 events by default and no fewer than 1,000. One event reducer assigns sequence values and performs redaction, durable append, snapshot update, then broadcast in that order.

**Rationale**: A single writer provides project-local monotonic order without shared-state races. JSON/JSONL is inspectable, offline, CGO-free, and sufficient for local-session replay. Atomic replace prevents partially written snapshots.

**Alternatives considered**:

- In-memory history: rejected because reconnecting clients lose events.
- SQLite now: deferred because it adds distribution complexity before normal session replay requires it.
- Per-client state mutation: rejected because it compromises ordering and race safety.

### 5. Define an explicit replay protocol and slow-client policy

**Decision**: Negotiate WebSocket subprotocol `craftlane.events.v1`. The client sends `{ "protocolVersion": 1, "lastSequence": <uint64> }`; the server registers the client, captures a snapshot cursor, replays retained events through that cursor in ascending order, then delivers later events. If the requested sequence predates retention, it sends `resync_required` and a fresh snapshot. Each client has a bounded outbound queue and is disconnected when it cannot keep up.

**Rationale**: Registration before bounded replay eliminates the snapshot/subscribe race. A clear resync result prevents clients from assuming continuity that no longer exists. Disconnecting slow clients protects canonical state progress.

**Alternatives considered**:

- Broadcast-only delivery: rejected because it loses disconnected updates.
- Snapshot followed by subscription without a cursor: rejected because it can create gaps or duplicates.
- Unbounded client queues: rejected because they can exhaust memory or stall the reducer.

### 6. Make security policy explicit at every browser and storage boundary

**Decision**: Accept only the exact runtime Host; accept a WebSocket Origin only when it is the exact runtime origin; return no permissive CORS headers. State-changing routes require a per-runtime control secret in `X-Craftlane-Token` and exact browser Origin when present. The control secret is valid only while its matching runtime is ready or stopping; it expires on shutdown, interruption, stale recovery, or restart. It is held in an owner-only runtime-control file, never returned by health/status/events/diagnostics, and removed on clean shutdown. All browser paths are validated against approved roots. A recursive copy-on-write redactor runs before event construction, persistence, logging, diagnostics, or presentation.

**Rationale**: Loopback alone does not protect a local service from malicious websites. Exact Host and Origin checks, a non-CORS-safe header, and a per-runtime secret prevent untrusted web origins from initiating a state change. Redaction before persistence avoids an irreversible leak.

**Alternatives considered**:

- Loopback-only or CORS-only protection: rejected because either leaves browser-origin attack paths.
- Static token, cookie, URL token, or local-storage token: rejected because they increase replay or disclosure risk.
- Log-only/UI-only redaction: rejected because raw values may already have leaked to a durable sink.

### 7. Embed a React/Vite build, and defer nonessential browser complexity

**Decision**: Keep contributor web source in `web/`, use React/TypeScript/Vite, and build static assets into a child directory of the Go server package before compiling or testing any Go package that embeds them. The binary embeds and serves those assets. Development may run Vite on an explicit local development origin that proxies `/api` and `/ws` to the runtime. The first dashboard uses semantic HTML, native Fetch and WebSocket, and no spatial view or client state library.

**Rationale**: It meets the planned web stack while keeping the executable self-contained. A package can only embed files below its own directory, so build output must be placed under the server package. Deferring Three.js and a state library prevents speculative complexity.

**Alternatives considered**:

- Ship static files beside the binary: rejected because it weakens single-binary installation.
- Run Vite in production: rejected because end users would need web-development tooling.
- Add Three.js or Zustand immediately: rejected because neither is needed for this dashboard.

### 8. Test the released integration boundary, not just its parts

**Decision**: Use Go unit and integration tests for project detection, records, event sequencing/replay, security, redaction, and shutdown. Use Vitest/React Testing Library for dashboard logic and Playwright against a built executable in a temporary Git repository.

**Rationale**: The highest release risk is the actual binary serving embedded assets and a local runtime, which dev-server-only tests do not cover. Temporary Git repositories test Git behavior without depending on the contributor checkout.

**Alternatives considered**:

- Browser-only tests: rejected because protocol, storage, and security failures can be masked.
- Vite-dev-server-only tests: rejected because they do not verify embedding or production routing.

## Sources

- [Go `net/http` server documentation](https://pkg.go.dev/net/http#Server.Serve)
- [Go `embed` documentation](https://pkg.go.dev/embed)
- [Go signal lifecycle documentation](https://pkg.go.dev/os/signal)
- [Go cancellable command documentation](https://pkg.go.dev/os/exec#CommandContext)
- [Vite production build guide](https://vite.dev/guide/build)
- [Playwright web server guide](https://playwright.dev/docs/test-webserver)
- [RFC 6455 Origin verification guidance](https://www.rfc-editor.org/rfc/rfc6455#section-10.2)
- [Fetch CORS preflight behavior](https://fetch.spec.whatwg.org/#cors-preflight-fetch)
