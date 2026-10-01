# WebSocket Contract: `craftlane.events.v1`

## Connection

- **Endpoint**: `GET /ws`
- **Required subprotocol**: `craftlane.events.v1`
- **Required production Host**: exact runtime host, `127.0.0.1:<actual-port>`
- **Required production Origin**: exact runtime origin, `http://127.0.0.1:<actual-port>`
- **Authentication**: no token is required for this read-only event stream; mutation secrets are never accepted in the URL, event payloads, or diagnostic messages.

The server rejects a missing or mismatched Origin before upgrading. It emits no permissive CORS headers.

## Client bootstrap message

The first client message must be sent promptly after connection:

```json
{
  "type": "subscribe",
  "protocolVersion": 1,
  "lastSequence": 42
}
```

`lastSequence` is the last event successfully applied by the client, or `0` for a new client.

## Server messages

### `snapshot`

```json
{
  "type": "snapshot",
  "cursor": 47,
  "snapshot": {
    "sequence": 47,
    "project": { "id": "…", "name": "…", "root": "…" },
    "runtime": { "id": "…", "origin": "http://127.0.0.1:7331", "protocolVersion": 1, "state": "ready", "startedAt": "…" },
    "diagnostics": []
  }
}
```

### `event`

```json
{
  "type": "event",
  "event": {
    "schemaVersion": 1,
    "eventId": "…",
    "sequence": 43,
    "timestamp": "…",
    "projectId": "…",
    "source": { "name": "core" },
    "type": "project.changed",
    "payload": {}
  }
}
```

### `resync_required`

```json
{
  "type": "resync_required",
  "cursor": 47,
  "reason": "requested_sequence_outside_retention"
}
```

The server follows this with a `snapshot`. The client discards its continuity assumption, applies the snapshot, records its cursor, and continues with events after that cursor.

The default retained window is the latest 10,000 project events. Local configuration may not lower the window below 1,000 events. Any requested sequence preceding the retained window receives `resync_required` rather than a partial replay.

### `error` and `heartbeat`

`error` contains a stable safe error code and an actionable message. `heartbeat` carries the latest known cursor and exists only to maintain and observe connection health; it does not change canonical state.

## Ordering and replay rules

1. The reducer assigns strictly increasing project-local sequence values and durably appends a redacted event before broadcasting it.
2. On subscription, the server registers the client, captures an upper cursor, replays retained events where `lastSequence < sequence <= cursor` in ascending order, then sends future events where `sequence > cursor`.
3. Clients discard any event at or below their applied cursor and deduplicate by `eventId`.
4. Each client queue is bounded. A client that cannot consume its queue is closed and must reconnect; it never blocks the reducer or other clients.
5. Unknown optional fields and unknown event types must be ignored safely by compatible clients. A breaking semantic change requires a protocol-version increase.
