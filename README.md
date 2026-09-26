# Particles ✨

A shared particle simulation you play in your terminal over SSH, served by a Go backend with tcell rendering.

## Getting Started

### Prerequisites

Toolchains (Go, Node, golangci-lint) are pinned in `mise.toml`. Install [mise](https://mise.jdx.dev), then:

```bash
mise install
```

You also need [websocat](https://github.com/vi/websocat) to connect over WebSocket.

```bash
go run . serve
ssh -p 2222 localhost
```

## Deployment

The app runs on [Cloudflare Containers](https://developers.cloudflare.com/containers/). A Worker in `cloudflare/` forwards requests to the
container built from `Dockerfile`. Deploying requires a Workers Paid plan and Docker running locally.

```bash
cd cloudflare && npm ci
npx wrangler login   # once
npx wrangler secret put PARTICLES__API__SSH_HOST_KEY < particles_host   # once
npx wrangler deploy
```

The Worker is served on the custom domain `particles.raphaelvanhoffelen.com` (`routes` in `cloudflare/wrangler.jsonc`),
and `PARTICLES__API__PUBLIC_HOST` sets the hostname shown on the instructions page.

## Play

SSH is tunneled over a WebSocket, so install [websocat](https://github.com/vi/websocat) first:

```bash
brew install websocat
ssh -o ProxyCommand="websocat --binary wss://particles.raphaelvanhoffelen.com/ssh" particles
```

Or add this to `~/.ssh/config` and run `ssh particles`:

```
Host particles
  ProxyCommand websocat --binary wss://particles.raphaelvanhoffelen.com/ssh
```

Click inside the box to drop particles, and press Esc or Ctrl+C to quit. [particles.raphaelvanhoffelen.com](https://particles.raphaelvanhoffelen.com)
shows the same instructions.

See [AGENTS.md](AGENTS.md) for more details on the project layout and conventions.
