---
name: quo-api-cli
description: Use a local Go CLI (`quoctl`) to interact with the Quo/OpenPhone REST API for contacts, messages, phone numbers, users, and raw endpoint calls. Trigger when the user asks to query/send/manage Quo data from terminal workflows, or to automate Quo API actions without writing curl each time.
---

# Quo API CLI Skill

Use the local `quoctl` binary to perform Quo API actions quickly and consistently.

---

## Preconditions

- Ensure binary exists at:
  - `projects/quo-cli-skill/quoctl/quoctl`
- Build if missing:

```bash
cd /Users/openclaw/.openclaw/workspace-ellie/projects/quo-cli-skill/quoctl
make build
```

- API key can come from:
  1. `--api-key`
  2. `QUO_API_KEY`
  3. `.quoctl.env`

- On first interactive run without a key, the CLI prompts and writes `.quoctl.env` automatically.

---

## Skill usage pattern

When user asks for Quo/OpenPhone data or actions:

1. Identify intent (contacts/messages/users/phone numbers/raw endpoint)
2. Run matching `quoctl` command
3. Return concise summary + key JSON fields
4. If API fails, include actionable error and next step

---

## OpenClaw prompt examples (what should trigger this skill)

- "List my latest 10 Quo contacts"
- "Get contact CT123abc"
- "Show recent messages for PN123abc with +12105559876"
- "Send a text from +12105551234 to +12105559876 saying ‘Tour confirmed for 3 PM.’"
- "List my OpenPhone users"
- "Call `/v1/conversations` and return the latest 20"

---

## Command examples used by this skill

### Contacts

```bash
./quoctl contacts list --max-results 10
./quoctl contacts get CT123abc
```

### Messages

```bash
./quoctl messages list \
  --phone-number-id PN123abc \
  --participants +12105559876 \
  --max-results 20

./quoctl messages send \
  --from +12105551234 \
  --to +12105559876 \
  --content "Tour confirmed for 3 PM."
```

### Phone numbers and users

```bash
./quoctl phone-numbers list
./quoctl users list --max-results 20
```

### Raw API mode

```bash
./quoctl api get /v1/conversations?maxResults=20
./quoctl api post /v1/contacts --data '{"defaultFields":{"firstName":"Joe"}}'
./quoctl api patch /v1/contacts/CT123abc --data '{"defaultFields":{"firstName":"Joseph"}}'
./quoctl api delete /v1/contacts/CT123abc
```

---

## Auth mode toggle

Default auth header:
- `Authorization: Bearer <key>`

If account expects raw key in Authorization header:

```bash
--auth-scheme ApiKey
```

---

## Troubleshooting

- **401 Unauthorized**
  - Verify API key
  - Try `--auth-scheme ApiKey`
- **404 Not Found**
  - Confirm endpoint begins with `/v1/...`
  - Confirm IDs are valid (`PN...`, `CT...`, etc.)
- **400 Validation error**
  - Check required fields (`content`, `from`, `to`, etc.)
  - Validate JSON in `--data`
