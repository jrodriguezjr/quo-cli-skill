# quoctl (Go)

`quoctl` is a Go CLI for interacting with the Quo/OpenPhone Public API.

It wraps common endpoints into easy terminal commands and also provides a raw `api` mode for any endpoint in the spec.

---

## What this CLI can do

- List and fetch contacts
- List messages and send SMS messages
- List phone numbers
- List users
- Call any API endpoint directly (`api get|post|patch|delete`)

---

## Prerequisites

- Go 1.22+
- Network access to `https://api.openphone.com`
- A valid Quo API key

---

## Build

```bash
make build
```

Binary output:
- `./quoctl`

---

## Authentication

`quoctl` resolves API key in this order:

1. `--api-key`
2. `QUO_API_KEY` environment variable
3. `.quoctl.env` in current directory (`QUO_API_KEY=...`)

### First-run experience

If no key is found on an **interactive** run, `quoctl` will:

1. Prompt for the key
2. Save it to `.quoctl.env`
3. Continue execution

The `.quoctl.env` file is written with mode `600`.

### One-off API key example

```bash
./quoctl users list --api-key 'YOUR_API_KEY' --max-results 5
```

### Raw auth header mode

Default header is:
- `Authorization: Bearer <key>`

If your account expects raw API key in the Authorization header, use:

```bash
./quoctl users list --auth-scheme ApiKey
```

---

## Global flags

These flags work on all subcommands:

- `--api-key <key>`
- `--base-url <url>` (default: `https://api.openphone.com`)
- `--auth-scheme <Bearer|ApiKey|None>`
- `--timeout <duration>` (default: `30s`)

Example:

```bash
./quoctl users list --max-results 10 --timeout 60s
```

---

## Command reference + examples

## 1) Contacts

### `contacts list`
List contacts with pagination.

```bash
./quoctl contacts list --max-results 10
./quoctl contacts list --max-results 10 --page-token '<nextPageToken>'
```

### `contacts get`
Get a single contact by ID.

```bash
./quoctl contacts get CT123abc
```

Use this when you already have a contact ID from a prior list/search result.

---

## 2) Messages

### `messages list`
List messages for a specific OpenPhone number and participants set.

Required:
- `--phone-number-id` (e.g., `PN...`)
- `--participants` (comma-separated E.164 list)

```bash
./quoctl messages list \
  --phone-number-id PN123abc \
  --participants +12105559876 \
  --max-results 20

./quoctl messages list \
  --phone-number-id PN123abc \
  --participants +12105559876,+12105550000 \
  --max-results 20 \
  --page-token '<nextPageToken>'
```

### `messages send`
Send an outbound message.

Required:
- `--from` (E.164 or `PN...`)
- `--to` (single recipient)
- `--content`

```bash
./quoctl messages send \
  --from +12105551234 \
  --to +12105559876 \
  --content "Hi — your maintenance request is scheduled for tomorrow at 10am."
```

Optional:
- `--user-id` to force sender context
- `--set-inbox-status done` to mark conversation done

```bash
./quoctl messages send \
  --from PN123abc \
  --to +12105559876 \
  --content "Closing this out now."
  --set-inbox-status done
```

---

## 3) Phone Numbers

### `phone-numbers list`
List phone numbers in the workspace.

```bash
./quoctl phone-numbers list
```

Useful for discovering the `PN...` IDs needed by other commands.

---

## 4) Users

### `users list`
List users with pagination.

```bash
./quoctl users list --max-results 20
./quoctl users list --max-results 20 --page-token '<nextPageToken>'
```

---

## 5) Raw API Mode

Use raw mode when a helper command does not exist yet.

### GET
```bash
./quoctl api get /v1/conversations?maxResults=20
```

### POST
```bash
./quoctl api post /v1/contacts --data '{"defaultFields":{"firstName":"Joe"}}'
```

### PATCH
```bash
./quoctl api patch /v1/contacts/CT123abc --data '{"defaultFields":{"firstName":"Joseph"}}'
```

### DELETE
```bash
./quoctl api delete /v1/contacts/CT123abc
```

---

## Output behavior

- JSON response is printed to stdout (pretty-printed when possible)
- Non-2xx responses still print body, then return non-zero exit with `HTTP <code>` error

This makes it script-friendly while preserving error details.

---

## Troubleshooting

### 401 Unauthorized
- Verify API key value
- Try switching auth scheme (`Bearer` vs `ApiKey`)

### 404 Not Found
- Confirm endpoint path starts with `/v1/...`
- Confirm IDs are valid (`PN...`, `CT...`, etc.)

### Validation errors (400)
- Check required fields (`content`, `from`, `to`, etc.)
- Validate JSON passed to `--data`

### Timeout/network errors
- Increase timeout: `--timeout 60s`
- Confirm outbound access to `api.openphone.com`
