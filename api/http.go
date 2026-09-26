package api

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

//go:embed templates/index.html
var templatesFS embed.FS

var indexTmpl = template.Must(template.ParseFS(templatesFS, "templates/index.html"))

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
		wsURL := sshWebSocketURL(r, a.Config.PublicHost)
		if !strings.Contains(r.Header.Get("Accept"), "text/html") {
			// curl and other non-browser clients get the command as plain text.
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = fmt.Fprintf(w, "ssh -o ProxyCommand=\"websocat --binary %s\" particles\n", wsURL)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := indexTmpl.Execute(w, struct{ WSURL string }{wsURL}); err != nil {
			a.Logger.Err(err).Msg("failed to render index")
		}
	})

	return mux
}

// sshWebSocketURL returns the URL clients pass to websocat. A configured public host is assumed to be served over TLS.
func sshWebSocketURL(r *http.Request, publicHost string) string {
	if publicHost != "" {
		return "wss://" + publicHost + "/ssh"
	}
	scheme := "ws"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" || r.Header.Get("Cf-Visitor") != "" {
		scheme = "wss"
	}
	return scheme + "://" + r.Host + "/ssh"
}
