# PARTICLES

A shared particle simulation you play over SSH. Everyone who connects sees the same simulation and gets their own particle color.

## Connect

The app is deployed as a Cloudflare Container behind a Worker. Workers only accept HTTP/WebSocket for now, so SSH
is tunneled over a WebSocket with [websocat](https://github.com/vi/websocat) (`brew install websocat`):

```sh
ssh -o ProxyCommand="websocat --binary wss://particles.<your-subdomain>.workers.dev/ssh" particles
```

Plain `ssh` will work once Cloudflare's inbound TCP beta is enabled for the Worker (see `cloudflare/src/index.ts`).

## Develop

```sh
mise install
go run . serve            # SSH on :2222, HTTP/WebSocket on :8080
ssh -p 2222 localhost
```

Configuration is read from `PARTICLES__*` env vars (or a YAML file via `--config`):

```yaml
API:
  MaxNumSessions: 5
```

## Deploy

```sh
cd cloudflare
npm install
ssh-keygen -t ed25519 -f particles_host -N ""
npx wrangler secret put PARTICLES__API__SSH_HOST_KEY < particles_host
npx wrangler deploy
```

Production config lives in `cloudflare/wrangler.jsonc` under `"vars"`. See [AGENTS.md](AGENTS.md) for architecture and details.
