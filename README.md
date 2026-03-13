# quo-cli-skill

Quo/OpenPhone API toolkit with two parts:

- `quoctl/` — Go CLI for Quo/OpenPhone API  
  → See [`quoctl/README.md`](./quoctl/README.md)
- `skills/quo-api-cli/` — OpenClaw skill wrapper that uses the CLI  
  → See [`skills/quo-api-cli/SKILL.md`](./skills/quo-api-cli/SKILL.md)

---

## Project structure

```text
quo-cli-skill/
├── quoctl/
│   ├── cmd/quoctl/main.go
│   ├── Makefile
│   └── README.md
└── skills/
    └── quo-api-cli/
        └── SKILL.md
```

---

## Quick start

```bash
cd quoctl
make build
./quoctl help
```

---

## Authentication

`quoctl` resolves API key in this order:
1. `--api-key`
2. `QUO_API_KEY`
3. `.quoctl.env` in current directory

If no key is found on an interactive first run, it prompts for a key and writes `.quoctl.env` with mode `600`.

---

## CLI usage examples (`quoctl`)

### Contacts

```bash
./quoctl contacts list --max-results 10
./quoctl contacts get <contact_id>
```

### Messages

```bash
./quoctl messages list \
  --phone-number-id PNxxxx \
  --participants +12105559876 \
  --max-results 20

./quoctl messages send \
  --from +12105551234 \
  --to +12105559876 \
  --content "Hello from quoctl"
```

### Phone numbers and users

```bash
./quoctl phone-numbers list
./quoctl users list --max-results 20
```

### Raw API calls

```bash
./quoctl api get /v1/conversations?maxResults=20
./quoctl api post /v1/contacts --data '{"defaultFields":{"firstName":"Joe"}}'
```

### One-off API key usage

```bash
./quoctl users list --api-key 'YOUR_API_KEY' --max-results 5
```

---

## OpenClaw skill usage examples

Skill file:
`skills/quo-api-cli/SKILL.md`

### Example prompts that should trigger/use the skill

- "List the latest 10 Quo contacts"
- "Send a message from +12105551234 to +12105559876 saying ‘Maintenance is scheduled for tomorrow.’"
- "Get recent conversations from Quo"
- "Create a contact named Joe via Quo API"

### Example command flow used by the skill

```bash
cd /Users/openclaw/.openclaw/workspace-ellie/projects/quo-cli-skill/quoctl
make build
./quoctl contacts list --max-results 10
./quoctl messages send --from +12105551234 --to +12105559876 --content "Hello"
```

---

## Notes

- Default auth header: `Authorization: Bearer <QUO_API_KEY>`
- If your account expects raw API key header, append:

```bash
--auth-scheme ApiKey
```
