# Data Model: Local Runtime Foundation

## Ownership and invariants

The runtime is the canonical owner of all entities below. The CLI reads or requests typed operations; the web dashboard only receives projections. The event reducer is the sole writer of project state and project-local event sequence values.

## Project

| Field | Description | Validation |
|---|---|---|
| `id` | Stable local project identity | Derived from canonical Git root and Git common directory; never from the display name alone. |
| `name` | Display name | Derived from the root directory name; treated as presentation only. |
| `root` | Canonical repository root | Absolute, symlink-resolved path returned by Git; must be a Git worktree. |
| `gitCommonDir` | Git metadata location | Derived by Git and retained for stable identity/worktree awareness. |
| `detectedAt` | Latest successful detection time | UTC timestamp. |
| `capabilities` | Local health of Git, Spec Kit, and providers | Optional integrations may be unavailable without making the project invalid. |

## Runtime Workspace

| Field | Description | Validation |
|---|---|---|
| `id` | Unique runtime identity | Cryptographically random per start. |
| `projectId` | Managed project | Must match the canonical detected project. |
| `origin` | Browser-visible local origin | Exact `http://127.0.0.1:<port>` address. |
| `protocolVersion` | Active public protocol version | Positive integer; initial version is 1. |
| `startedAt` | Start time | UTC timestamp. |
| `state` | Runtime lifecycle state | `starting`, `ready`, `stopping`, `stopped`, or `stale`. |

**Transitions**: `starting → ready → stopping → stopped`. A persisted record that cannot be confirmed by health becomes `stale`; `start` may replace it only after cleanup. A failed start transitions directly to `stopped` and must release resources.

## Runtime Record and Runtime Control

| Entity | Persisted fields | Rules |
|---|---|---|
| `RuntimeRecord` | Runtime ID, project ID/root, origin, process ID as advisory metadata, protocol version, start time | Written atomically; health response matching runtime and project is authoritative over the process ID. |
| `RuntimeControl` | Runtime ID, random control secret, expiry/lifetime | Valid only while the matching runtime ID is in the `ready` or `stopping` state. It expires immediately on shutdown, interruption, stale-record recovery, or restart; it is stored separately with owner-only permissions; omitted from status, health, events, logs, diagnostics, and support output. |

## Event

| Field | Description | Validation |
|---|---|---|
| `schemaVersion` | Wire-schema version | Positive integer; existing semantics are never silently changed. |
| `eventId` | Stable event identity | Unique locally generated identifier. |
| `sequence` | Project-local total order | Assigned only by the reducer; strictly greater than the prior persisted sequence. |
| `timestamp` | Event creation time | UTC timestamp. |
| `projectId` | Event's project | Must equal the reducer's project. |
| `source` | Producing subsystem | Provider-neutral source name and optional adapter version; initial source is `core`. |
| `type` | Event meaning | Initial set: `project.detected`, `project.changed`, `runtime.ready`, `runtime.stopping`. |
| `payload` | Type-specific normalized data | Redacted copy only; must validate against the versioned schema. |

**Persistence lifecycle**: redact input → assign sequence/ID → append event → update snapshot → broadcast. A client only advances its cursor after it has applied the event. Retention cleanup never removes the latest valid snapshot; a reconnect behind the retained range receives `resync_required` and that snapshot.

**Retention policy**: The default window is the most recent 10,000 events per project. Configuration may reduce that window no lower than 1,000 events; it may increase it only through explicit local configuration. Events preceding the retained window are unavailable and always produce `resync_required`, never a partial replay presented as continuous history.

## Event Snapshot

| Field | Description | Validation |
|---|---|---|
| `sequence` | Cursor represented by the snapshot | Must be the latest durably persisted event sequence. |
| `project` | Current project projection | Contains only browser-safe project fields. |
| `runtime` | Current health-safe runtime projection | Excludes control secret and filesystem locations outside the project. |
| `diagnostics` | Current typed checks | Every failed/warned check includes a user-facing next action. |

## Diagnostic Check

| Field | Description | Validation |
|---|---|---|
| `name` | Stable check identifier | Examples: `git`, `runtime-record`, `loopback-listener`, `embedded-dashboard`, `speckit`, `provider`. |
| `status` | Result severity | `pass`, `warn`, or `fail`. |
| `detail` | Safe explanatory detail | Must be redacted and omit secrets/private unrelated paths. |
| `action` | Corrective next action | Required for `warn` and `fail`; optional for `pass`. |

## Path Authorization Request

| Field | Description | Validation |
|---|---|---|
| `path` | Requested local location | Canonicalized before use. |
| `purpose` | Typed operation requesting access | Required; arbitrary command execution is not a purpose. |
| `approvedRoot` | Project, known workspace, or user-selected root | Request path must be contained by this exact canonical root. |

No browser request can expand its own approved root. This milestone exposes no generic filesystem-browsing or shell-execution operation.
