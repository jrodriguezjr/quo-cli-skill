# Webhooks Design Plan (Quo/OpenPhone)

Status: Draft (pre-implementation)

## Goal

Add webhook management to `quoctl` so users can create/list/get/delete webhooks from the CLI, then consume events in their own receiver service.

## Why

Current CLI is pull-based (polling endpoints). Webhooks enable push-based event delivery for lower latency and fewer API polling calls.

---

## API surface from OpenAPI spec

- `GET /v1/webhooks`
- `GET /v1/webhooks/{id}`
- `DELETE /v1/webhooks/{id}`
- `POST /v1/webhooks/messages`
- `POST /v1/webhooks/calls`
- `POST /v1/webhooks/call-summaries`
- `POST /v1/webhooks/call-transcripts`

Create endpoints require:
- `url` (callback URL)
- `events` (event names)

Optional fields include:
- `label`
- `resourceIds`
- `status`
- `userId`

---

## Proposed CLI commands

### Read operations

- `quoctl webhooks list`
- `quoctl webhooks get <id>`

### Delete operation

- `quoctl webhooks delete <id>`

### Create operations

- `quoctl webhooks create messages --url <https://...> --events <csv> [--label ...] [--resource-ids <csv>] [--status active|inactive] [--user-id US...]`
- `quoctl webhooks create calls --url <https://...> --events <csv> [...]`
- `quoctl webhooks create call-summaries --url <https://...> --events <csv> [...]`
- `quoctl webhooks create call-transcripts --url <https://...> --events <csv> [...]`

---

## Validation rules (planned)

1. `--url` required and must start with `https://` (production recommendation).
2. `--events` required and non-empty CSV.
3. `--resource-ids` optional CSV; split to array.
4. `--status` optional; if provided, validate allowed values from spec.
5. Return non-zero exit code on API errors (consistent with existing CLI behavior).

---

## Assumptions to verify before coding

1. Exact allowed `events` values per webhook type.
2. Whether callback signing headers are provided by Quo (and exact verification method).
3. Retry/backoff behavior for non-2xx receiver responses.
4. Timeouts and max payload size expectations.

---

## Security / operations guidance (repo docs)

- Do not expose private/internal endpoints directly.
- Use HTTPS callback URLs.
- Verify webhook signatures when available.
- Implement idempotency in receiver (dedupe by event id).
- Log payload IDs + timestamps for replay/debugging.

---

## Rollout plan

### Phase 1 (this repo)
- Add webhook CRUD/create CLI commands.
- Add docs + examples in `quoctl/README.md` and root `README.md`.
- Add OpenClaw skill examples for webhook management commands.

### Phase 2 (optional)
- Add minimal receiver example service (`examples/webhook-receiver`).
- Add signature verification examples if/when documented.

---

## Example commands (target UX)

```bash
# list
./quoctl webhooks list

# get
./quoctl webhooks get WH123abc

# create messages webhook
./quoctl webhooks create messages \
  --url https://example.com/webhooks/quo/messages \
  --events message.received,message.sent \
  --label "ops-message-events"

# delete
./quoctl webhooks delete WH123abc
```
