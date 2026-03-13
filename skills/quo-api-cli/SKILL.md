---
name: quo-api-cli
description: Use a local Go CLI (`quoctl`) to interact with the Quo/OpenPhone REST API for contacts, messages, phone numbers, users, and raw endpoint calls. Trigger when the user asks to query/send/manage Quo data from terminal workflows, or to automate Quo API actions without writing curl each time.
---

# Quo API CLI Skill

Use the local `quoctl` binary to perform Quo API actions quickly and consistently.

## Preconditions

- Ensure binary exists at `projects/quo-cli-skill/quoctl/quoctl` (build it if missing).
- API key can come from `--api-key`, `QUO_API_KEY`, or `.quoctl.env`.
- On first interactive run without a key, the CLI prompts and writes `.quoctl.env` automatically.

## Build

```bash
cd /Users/openclaw/.openclaw/workspace-ellie/projects/quo-cli-skill/quoctl
make build
```

## Common commands

```bash
./quoctl contacts list --max-results 10
./quoctl contacts get <contact_id>
./quoctl messages list --phone-number-id PNxxxx --participants +12105559876 --max-results 20
./quoctl messages send --from +12105551234 --to +12105559876 --content "Hello"
./quoctl phone-numbers list
./quoctl users list --max-results 20
./quoctl api get /v1/conversations?maxResults=20
./quoctl api post /v1/contacts --data '{"defaultFields":{"firstName":"Joe"}}'
```

## Auth mode toggle

Default auth header is `Authorization: Bearer <key>`.
If the account expects raw API key header value, run commands with:

```bash
--auth-scheme ApiKey
```

## Troubleshooting

- If API returns 401, verify key and auth scheme.
- If API returns 404, confirm endpoint path starts with `/v1/...`.
- If JSON parse errors occur, validate `--data` payload is valid JSON.
