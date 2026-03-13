# quo-cli-skill

Quo/OpenPhone API toolkit:

- `quoctl/` — Go CLI (`contacts`, `messages`, `phone-numbers`, `users`, `api`)
- `skills/quo-api-cli/` — OpenClaw skill wrapper for using the CLI

## Quick start

```bash
cd quoctl
make build
./quoctl help
```

## Auth

`quoctl` resolves API key in this order:
1. `--api-key`
2. `QUO_API_KEY`
3. `.quoctl.env`

On first interactive run without a key, it prompts and saves `.quoctl.env` with mode `600`.
