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

Connect with:

```bash
ssh -o ProxyCommand="websocat --binary wss://particles.<your-subdomain>.workers.dev/ssh" particles
```

See [AGENTS.md](AGENTS.md) for more details on the project layout and conventions.
