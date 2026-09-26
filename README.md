# Particles

[![lint](https://github.com/dskart/particles/actions/workflows/lint.yml/badge.svg)](https://github.com/dskart/particles/actions/workflows/lint.yml)
[![test](https://github.com/dskart/particles/actions/workflows/test.yml/badge.svg)](https://github.com/dskart/particles/actions/workflows/test.yml)
[![build](https://github.com/dskart/particles/actions/workflows/build.yml/badge.svg)](https://github.com/dskart/particles/actions/workflows/build.yml)

A shared particle simulation you play in your terminal over SSH. Everyone who connects sees the same world and
gets their own particle color.

## Play

Tunnel SSH through a WebSocket with [websocat](https://github.com/vi/websocat) (`brew install websocat`):

```sh
ssh -o ProxyCommand="websocat --binary wss://particles.<your-subdomain>.workers.dev/ssh" particles
```

- **Click** in the box to drop particles
- **Esc** or **Ctrl+C** to quit
- Sessions end automatically after 5 minutes

> Cloudflare Workers only accept HTTP/WebSocket today. Plain `ssh` will work once Cloudflare's inbound TCP beta is
> enabled for the Worker.

## Develop

Tools are managed with [mise](https://mise.jdx.dev):

```sh
mise install
(cd cloudflare && npm ci)

go run . serve        # SSH on :2222, HTTP/WebSocket on :8080
ssh -p 2222 localhost

mise run fmt          # format Go + Worker
mise run lint         # format check, golangci-lint, prettier, tsc
mise run vet          # go vet
mise run test         # go test -race
mise run build        # Go binary + container image (needs Docker)
```

Configuration comes from `PARTICLES__*` env vars. mise loads them from a local `.env` if present. See
[AGENTS.md](AGENTS.md#configuration) for the full list.

## Deploy

The app runs as a [Cloudflare Container](https://developers.cloudflare.com/containers/) behind a Worker
(`cloudflare/`). Production config is under `"vars"` in `cloudflare/wrangler.jsonc`.

```sh
cd cloudflare
ssh-keygen -t ed25519 -f particles_host -N ""
npx wrangler secret put PARTICLES__API__SSH_HOST_KEY < particles_host
npx wrangler deploy
```

Deploying requires a Workers Paid plan and Docker running locally. See [AGENTS.md](AGENTS.md) for architecture
details.
