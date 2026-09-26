# AGENTS.md

Guidance for coding agents working in this repository.

## Overview

`particles` is a Go SSH server that streams one **shared** particle physics simulation to every connected
client. Clients render it in their terminal (tcell, 2x2 Unicode block characters) and click to add particles.
It is deployed as a [Cloudflare Container](https://developers.cloudflare.com/containers/) fronted by a Worker.

## Tooling

Tool versions are pinned in `mise.toml` (Go 1.24, Node 22). Run `mise install` once, then either activate mise in
your shell or prefix commands with `mise exec --`. mise also loads the gitignored `.env` (local `PARTICLES__*`
overrides) if it exists. There is no direnv.

## Commands

The standard checks are mise tasks (`mise tasks ls`). CI (`.github/workflows/`) runs the same tasks.

```sh
(cd cloudflare && npm ci)   # once, needed by the worker tasks
mise run fmt                # golangci-lint fmt (gofmt, goimports) + prettier
mise run lint               # format check + golangci-lint (.golangci.yml) + prettier --check + tsc
mise run vet                # go vet ./...
mise run test               # go test -race ./...
mise run build              # go build + `wrangler deploy --dry-run` (builds the container image, needs Docker)
```

```sh
# Go app
go run . serve                           # SSH on :2222, HTTP on :8080
go run . serve --port 2222 --http-port 0 # SSH only
ssh -p 2222 localhost                    # connect directly
ssh -o ProxyCommand="websocat --binary ws://localhost:8080/ssh" particles  # connect over WebSocket

# Cloudflare Worker + Container (from cloudflare/)
npm install
npx wrangler dev          # builds the root Dockerfile and runs the container in local Docker, on :8787
npx wrangler types        # regenerate worker-configuration.d.ts after editing wrangler.jsonc
npm run typecheck
npx wrangler deploy       # needs Docker running and a Workers Paid plan
```

## Architecture

- `main.go` → `cmd/` (cobra)
  - `root.go`: loads config (YAML via `--config`, then `PARTICLES__*` env vars), logger, signal handling.
  - `serve.go`: runs the simulation, the SSH server (`--port`, default 2222) and the HTTP server
    (`--http-port`, default 8080, `0` disables) in one errgroup.
- `app/`: simulation and rendering.
  - `simulator.go`: physics loop; a single `Simulator` is shared by all sessions.
  - `particle.go`, `renderer.go`, `box.go`, `colors.go`: particles, 2x2 block rendering, viewport, per-client colors.
  - `app.go`: `HandleSSHSession`, the per-client render/input loop with a session timer (`MaxSessTime`).
  - `sshtty.go`, `screen.go`, `virtual_screen.go`: adapt an SSH session to a `tcell.Screen`.
- `api/`: transports.
  - `api.go`: gliderlabs SSH server, session limit (`MaxNumSessions`), host key parsing.
  - `http.go`: `/ping` and `/healthz` (health checks), `/ssh` (SSH over WebSocket), `/` (connection instructions).
  - `wsconn.go`: adapts a gorilla/websocket connection to `net.Conn` (binary frames) so `/ssh` can hand it
    to `ssh.Server.HandleConn`.
- `pkg/config`: YAML + env unmarshalling. `pkg/logger`: zerolog. `pkg/shutdown`: shutdown hooks.
- `cloudflare/`: Worker (`src/index.ts`) and `wrangler.jsonc`. The container image is the root `Dockerfile`.

## Configuration

`cmd.Config` is filled from env vars named `PARTICLES__<SECTION>__<FIELD>`. The segments come from the `env` struct
tags, joined with `__`.

| Env var                               | Default | Field                       |
| ------------------------------------- | ------- | --------------------------- |
| `PARTICLES__APP__MAX_SESS_TIME`       | `5m`    | `app.Config.MaxSessTime`    |
| `PARTICLES__APP__SIM_CONFIG__WIDTH`   | `80`    | `app.SimConfig.Width` (≥80) |
| `PARTICLES__APP__SIM_CONFIG__HEIGHT`  | `25`    | `app.SimConfig.Height` (≥25) |
| `PARTICLES__APP__SIM_CONFIG__GRAVITY` | `1`     | `app.SimConfig.Gravity`     |
| `PARTICLES__API__MAX_NUM_SESSIONS`    | `10`    | `api.Config.MaxNumSessions` |
| `PARTICLES__API__SSH_HOST_KEY`        | (none)  | `api.Config.SSHHostKey`: PEM private key; if unset, a key is generated at startup |

When you add a config field, also add it to `"vars"` in `cloudflare/wrangler.jsonc` (or to `"secrets"` if it is
sensitive) and run `npx wrangler types`. The Worker forwards every `PARTICLES__*` binding to the container, so no
TypeScript change is needed.

## Cloudflare deployment

```
ssh ─ProxyCommand websocat─► wss://<worker>/ssh ► Worker fetch() ► ParticlesContainer (Durable Object)
    ► container :8080 /ssh ► ssh.Server.HandleConn
```

- **Single instance:** the simulation is in-memory shared state. The Worker always routes to
  `getContainer(env.PARTICLES, "main")`, and `max_instances` is 1. Don't use `getRandom` or scale out.
- **HTTP/WebSocket only:** Workers can't accept raw TCP yet. Inbound TCP through Spectrum is a Cloudflare private
  beta. `connect()` handlers are already wired in `src/index.ts` (Worker → DO → `getTcpPort(2222)`) for when it
  becomes available.
- **Lifecycle:** the container sleeps after `sleepAfter` (15m) without requests, so the first request after that
  cold-starts it and the simulation resets. Keep `sleepAfter` longer than `PARTICLES__APP__MAX_SESS_TIME`.
- **Host key:** set `PARTICLES__API__SSH_HOST_KEY` with `npx wrangler secret put PARTICLES__API__SSH_HOST_KEY < key`.
  Without it, the fingerprint changes on every cold start. Locally, put it in `cloudflare/.dev.vars`
  (see `.dev.vars.example`).
- The image must build for `linux/amd64`. Wrangler builds the root `Dockerfile` with the repo root as context
  (see `.dockerignore`).

## Conventions

- Log with zerolog (`*zerolog.Logger` passed down, sub-loggers via `.With()`); don't use `fmt.Println` for logs.
- Wrap errors with `fmt.Errorf("...: %w", err)`.
- Run long-lived goroutines in the `serve` errgroup, and register cleanup with `shutdown.OnShutdown`.
- Tests use `testify/require`.
