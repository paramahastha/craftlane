# Contract: Provider Observation Ingestion and Projection

## Ingestion

`POST /api/v1/observe/{provider}` where `{provider}` is exactly `claude`, `codex`, or `opencode`.

The route is local bridge infrastructure, not a browser API. It accepts only `POST` to the exact runtime Host and either no `Origin` header (the private local bridge) or an `Origin` exactly matching the active runtime origin; every other Origin is rejected. It also requires a bounded JSON body, a valid separate `X-Craftlane-Observation-Token`, a supported envelope version, matching project identity, and an approved canonical project path when present. It returns safe errors only, emits no permissive CORS headers, and does not log raw bodies.

```json
{
  "version": 1,
  "provider": "codex",
  "projectId": "local-project-id",
  "kind": "tool.completed",
  "sessionKey": "opaque-session-id",
  "agentKey": "opaque-agent-id",
  "toolKey": "opaque-tool-id",
  "toolName": "apply_patch",
  "outcome": "completed",
  "reportedAt": "2026-10-02T12:00:00Z"
}
```

Unknown/forbidden raw fields are rejected or discarded before persistence. In particular, no prompts, messages, tool args/results, commands, transcript data, environment data, or credential values are valid envelope fields.

## Event and snapshot additions

The existing v1 event envelope is unchanged. New normalized events use source `{ "name": "claude" | "codex" | "opencode", "adapterVersion": "…" }`, type `agent.updated`, and a complete browser-safe agent projection. A local event sequence orders accepted Craftlane events only; it does not assert provider causal order.

Snapshots gain optional `agents` and `providerStatuses` collections. Provider status lists available/unavailable capabilities and safe reasons. Compatible existing clients ignore the additive data; updated clients render capability evidence rather than assuming parity.
