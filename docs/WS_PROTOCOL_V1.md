# WebSocket Protocol v1

This document defines the application-level protocol for WebSocket events in Maicare.

## Goals

- Support multiple event domains beyond notifications.
- Keep client routing stable as payloads evolve.
- Preserve backward compatibility during migration.

## Connection Flow

1. Client requests one-time ticket: `POST /auth/ws-ticket` with `Authorization: Bearer <access_token>`.
2. Server returns ticket, expiration, and WS path.
3. Client connects: `GET /ws?ticket=<ticket>`.
4. Server validates and consumes ticket exactly once, then upgrades the socket.

## Event Envelope (Required)

Every server-to-client WS event MUST be a single JSON object with this envelope:

```json
{
  "v": "1",
  "id": "01HV7V8Q2H4R0X9K8K6W7E2M3N",
  "type": "notification.created",
  "ts": "2026-03-04T12:30:05Z",
  "data": {}
}
```

### Envelope Fields

- `v` (string, required): protocol version. For this spec: `"1"`.
- `id` (string, required): event id (ULID or UUID).
- `type` (string, required): event name (`domain.entity.action`).
- `ts` (string, required): RFC3339 UTC send timestamp.
- `data` (object, required): event payload for `type`.
- `seq` (integer, optional): per-connection sequence number for gap detection.
- `correlation_id` (string, optional): request/job trace linkage.
- `idempotency_key` (string, optional): dedupe key when applicable.
- `replay` (boolean, optional): true when event is replayed/catch-up.
- `error` (object, optional): present for `system.error` events.

## Event Naming

- Format: `domain.entity.action`
- Examples:
  - `notification.created`
  - `presence.updated`
  - `system.error`

Reserved prefix: `system.*` for protocol/meta events.

## Canonical Event Examples

### `notification.created`

```json
{
  "v": "1",
  "id": "01HV7V8Q2H4R0X9K8K6W7E2M3N",
  "type": "notification.created",
  "ts": "2026-03-04T12:30:05Z",
  "seq": 1842,
  "data": {
    "notification_id": "8d4c7f9f-6ab2-4ab9-a9be-fbd7f5c0e1f1",
    "notification_type": "new_appointment",
    "message": "New appointment created",
    "is_read": false,
    "created_at": "2026-03-04T12:30:04Z",
    "data": {
      "new_appointment": {
        "appointment_id": "f4f8d5da-1111-4444-9999-4f74d7b4e118"
      }
    }
  }
}
```

### `presence.updated`

```json
{
  "v": "1",
  "id": "01HV7V9A7N6S1Q2R3T4U5V6W7X",
  "type": "presence.updated",
  "ts": "2026-03-04T12:30:10Z",
  "seq": 1843,
  "data": {
    "user_id": "dcb9ec45-c112-4b4d-bf77-e7ad3f4a2bc6",
    "status": "online",
    "last_seen_at": "2026-03-04T12:30:10Z",
    "source": "web"
  }
}
```

### `system.error`

```json
{
  "v": "1",
  "id": "01HV7V9W3Y1A2B3C4D5E6F7G8H",
  "type": "system.error",
  "ts": "2026-03-04T12:30:12Z",
  "data": {},
  "error": {
    "code": "RATE_LIMITED",
    "message": "Too many subscriptions",
    "retryable": true
  }
}
```

## Framing and Transport Rules

- One WS frame MUST contain exactly one JSON event envelope.
- No concatenated JSON objects in one frame.
- Server sends events in-order per connection.
- Clients should treat `id` as dedupe key.

## Compatibility and Versioning

- Existing notification-only payload is legacy `v0` behavior.
- New clients should consume only v1 envelopes.
- During migration, server may dual-publish (v0 and v1) until all clients migrate.
- Additive changes are allowed in `data`; envelope required fields stay stable.

## Reconnect and Catch-up

- Client reconnects with exponential backoff + jitter.
- On reconnect, client must request a new WS ticket.
- For notification catch-up, client calls existing REST notifications API.
- Replay semantics can be added later using `replay=true` and `since` parameters.

## Validation

- Backend validates envelope (`v`, `id`, `type`, `ts`, `data`) before send.
- Frontend validates envelope and dispatches by `type`.
- Unknown `type` should be ignored safely and logged at debug level.

## Observability (Protocol-Level)

Track at minimum:

- `ws_connections_active`
- `ws_events_sent_total{type}`
- `ws_events_dropped_total{reason}`
- `ws_auth_failures_total{reason}`
- `ws_frame_size_bytes`

Structured logs should include: `conn_id`, `user_id`, `protocol_v`, `event_id`, `type`, `seq`.

## Migration Plan (Repo Touchpoints)

1. **Stabilize framing**
   - Update writer behavior in `hub/client.go` to stop concatenating messages.
2. **Introduce v1 envelope generation**
   - Start in `service/notification/service.go` (emit envelope for notification events).
3. **Protocol routing/versioning**
   - Add v1 route/subprotocol handling in `api/websocket_handler.go`.
4. **Client rollout**
   - Frontend consumes envelope by `type` and ignores unknown events.
5. **Legacy removal**
   - Remove v0 payload mode after migration completion.
