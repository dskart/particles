import { Container, getContainer } from "@cloudflare/containers";

const HTTP_PORT = 8080;
const SSH_PORT = 2222;

// All clients share one simulation, so every connection goes to this single instance.
const INSTANCE_NAME = "main";

// Go app config (cmd.Config) is read from PARTICLES__* env vars; see "vars" in wrangler.jsonc.
const CONFIG_PREFIX = "PARTICLES__";

function appConfigEnv(env: Env): Record<string, string> {
	return Object.fromEntries(
		Object.entries(env).filter(
			(entry): entry is [string, string] =>
				entry[0].startsWith(CONFIG_PREFIX) && typeof entry[1] === "string" && entry[1] !== "",
		),
	);
}

export class ParticlesContainer extends Container<Env> {
	// Health checks and SSH-over-WebSocket (/ssh) are served on HTTP_PORT.
	defaultPort = HTTP_PORT;
	requiredPorts = [HTTP_PORT, SSH_PORT];
	// Must outlive PARTICLES__APP__MAX_SESS_TIME so active sessions are not cut off.
	sleepAfter = "15m";

	constructor(ctx: ConstructorParameters<typeof Container<Env>>[0], env: Env) {
		super(ctx, env);
		this.envVars = appConfigEnv(env);
	}

	// Raw TCP (plain `ssh`) forwarded from the Worker's connect() handler.
	// Requires Cloudflare's inbound TCP beta (Spectrum -> Worker).
	async connect(socket: Socket): Promise<void> {
		await this.startAndWaitForPorts(SSH_PORT);
		const containerSocket = this.ctx.container!.getTcpPort(SSH_PORT).connect(`10.0.0.1:${SSH_PORT}`);
		await containerSocket.opened;
		await Promise.all([
			socket.readable.pipeTo(containerSocket.writable),
			containerSocket.readable.pipeTo(socket.writable),
		]);
	}

	override onStart() {
		console.log("particles container started");
	}

	override onStop() {
		console.log("particles container stopped");
	}

	override onError(error: unknown) {
		console.error("particles container error:", error);
		throw error;
	}
}

export default {
	// HTTP + WebSocket: `/` (instructions), `/healthz`, and `/ssh` (SSH over WebSocket).
	async fetch(request, env): Promise<Response> {
		return getContainer(env.PARTICLES, INSTANCE_NAME).fetch(request);
	},

	// Inbound TCP (beta): route raw SSH connections to the shared container.
	async connect(socket, env): Promise<void> {
		const stub = getContainer(env.PARTICLES, INSTANCE_NAME);
		const upstream = stub.connect(`particles:${SSH_PORT}`);
		await Promise.all([socket.readable.pipeTo(upstream.writable), upstream.readable.pipeTo(socket.writable)]);
	},
} satisfies ExportedHandler<Env>;
