package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Clients are SSH ProxyCommands (e.g. websocat), not browsers.
	CheckOrigin: func(*http.Request) bool { return true },
}

// NewHTTPHandler serves health checks and an SSH-over-WebSocket bridge on /ssh,
// used when the SSH port cannot be exposed directly (e.g. Cloudflare Containers).
func NewHTTPHandler(a *API) http.Handler {
	mux := http.NewServeMux()

	ok := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, "ok")
	}
	mux.HandleFunc("GET /ping", ok)
	mux.HandleFunc("GET /healthz", ok)

	mux.HandleFunc("GET /ssh", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			a.Logger.Err(err).Msg("websocket upgrade failed")
			return
		}
		a.Logger.Info().Str("remoteAddr", r.RemoteAddr).Msg("SSH over websocket connection")
		a.Server.HandleConn(newWSConn(c))
	})

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		scheme := "wss"
		if r.TLS == nil && r.Header.Get("X-Forwarded-Proto") != "https" && r.Header.Get("Cf-Visitor") == "" {
			scheme = "ws"
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintf(w, "Particles: a shared particle simulation over SSH.\n\n"+
			"Connect with websocat (https://github.com/vi/websocat):\n\n"+
			"  ssh -o ProxyCommand=\"websocat --binary %s://%s/ssh\" particles\n", scheme, r.Host)
	})

	return mux
}
