# Craftlane protocol

The foundation protocol is version 1. HTTP and WebSocket payloads are local-only,
additive, and browser-safe. Breaking semantic changes require a new version.

Observation adds optional `agents` and `providerStatuses` collections to
snapshots and `agent.updated` events with a provider source. Existing clients
may ignore those fields. Provider callbacks use `POST /api/v1/observe/{provider}`
with a separate observation token; requests require the exact Host and either
no Origin or the exact runtime Origin.
