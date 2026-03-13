# quoctl (Go)

CLI for Quo/OpenPhone API.

## Prereqs
- Go 1.22+
- `QUO_API_KEY` env var

## Build
```bash
make build
```

## Usage
```bash
./quoctl contacts list --max-results 10
./quoctl contacts get <contact_id>
./quoctl messages list --phone-number-id PNxxxx --participants +12105559876 --max-results 20
./quoctl messages send --from +12105551234 --to +12105559876 --content "Hello from quoctl"
./quoctl phone-numbers list
./quoctl users list --max-results 20
./quoctl api get /v1/conversations?maxResults=20
```

## Auth notes
By default, requests use `Authorization: Bearer <QUO_API_KEY>`.

API key resolution order:
1. `--api-key`
2. `QUO_API_KEY` environment variable
3. `.quoctl.env` in current working directory (`QUO_API_KEY=...`)

On first interactive run, if no key is found, `quoctl` prompts for the key and saves it to `.quoctl.env` (mode `600`).

If your workspace expects a raw header value, use:
```bash
./quoctl users list --auth-scheme ApiKey
```
