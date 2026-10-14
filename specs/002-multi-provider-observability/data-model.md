# Data Model: Multi-Provider Agent Observability

## Ownership and invariants

The runtime reducer is the sole owner of normalized agent state, processed-event fingerprints, local sequence, event retention, and snapshots. Provider adapters own only raw-input projection/normalization. Browser clients receive safe projections and never derive authoritative state.

## Provider Status and Capability

| Field | Description | Validation |
|---|---|---|
| `provider` | Adapter identifier | One of `claude`, `codex`, or `opencode`; generic consumers treat it as opaque. |
| `setupState` | Integration readiness | `not_configured`, `needs_trust`, `ready`, `degraded`, or `unavailable`. |
| `capabilities` | Per-category coverage | Map of sessions, subagents, tools, permissions, failures, terminal states, and idle states to `available`/`unavailable` with safe reason. |
| `detail` / `action` | Safe readiness explanation | Redacted; unhealthy states include next action. |
| `adapterVersion` | Adapter/provider compatibility marker | Optional safe version; no environment or account data. |

## Agent and Tool Activity

| Field | Description | Validation |
|---|---|---|
| `id` | Stable local agent identity | Derived from provider + project + opaque provider agent/session key. |
| `provider`, `projectId`, `sessionId` | Ownership/correlation | Provider and project must match the active adapter/runtime. |
| `parentAgentId` | Supplied hierarchy edge | Set only from authoritative provider lineage; never inferred from activity. |
| `state`, `stateSource` | Normalized lifecycle condition | `working`, `waiting`, `waiting_permission`, `idle`, `completed`, `failed`, or `unknown`; source is `reported` or `inferred`. |
| `currentTool` | Safe current/recent operation | Contains opaque correlation ID, allowlisted name, started/completed/failed state, and coarse outcome only. |
| `lastObservedAt`, `reportedAt` | Timing | Receipt time and optional supplied timestamp; never claim receipt order as provider order. |
| `priorHistoryUnavailable` | Late-observation marker | True when first signal follows session start. |
| `staleAfterRestart` | Recovery marker | True until an active pre-restart agent receives new evidence. |

## Observation Envelope and Fingerprint

```text
version, provider, projectId, kind, sessionKey?, agentKey?, parentAgentKey?,
toolKey?, toolName?, outcome?, reportedAt?
```

No envelope may contain raw provider input/output/prompt/message/path/command/environment data. The callback validates a bounded envelope before adapter use. A deterministic fingerprint uses provider, project, lifecycle kind, and stable correlation identifiers. Fingerprints are bounded by local observation retention; a duplicate produces no state/event mutation.

## State rules

| Evidence | Result |
|---|---|
| Supported session/activity signal | Upsert working/unknown agent; mark prior history unavailable when no earlier start was observed. |
| Provider tool start/completion/failure | Update tool activity without retaining raw data. |
| Explicit permission signal | `waiting_permission`, reported source. |
| Supported explicit idle/wait signal | `idle`/`waiting`, reported source. |
| Inactivity threshold for live session | `idle`, inferred source only. |
| Provider-specific durable terminal signal | `completed`/`failed`, reported source. |
| Late nonterminal or duplicate record after terminal state | Ignore; do not revive agent. |
| Runtime restart with formerly active agent | `unknown`, `staleAfterRestart = true` until new evidence. |

Provider semantics determine whether a signal is durable terminal evidence: Claude session end/subagent stop can be terminal; Codex turn Stop is not session end and SessionEnd may be delayed; OpenCode busy/idle/retry signals are not successful terminal completion. Capability status makes these differences visible.

## Snapshot and event projection

`Snapshot.agents` and `Snapshot.providerStatuses` are optional additive fields. Any normalized state change emits `agent.updated` through the existing v1 event envelope with source `{name: provider, adapterVersion: ...}` and a complete safe agent projection. Existing clients ignore unknown additions; updated clients upsert by agent ID and use capability status rather than provider-specific branching.

## Persistence lifecycle

`provider input → bounded bridge projection → callback authorization/validation → adapter normalization → reducer dedup/state guard → redact → append event → save snapshot → broadcast`

Only the normalized/redacted envelope, agent projection, fingerprints, existing bounded event history, and diagnostics are persisted locally.
