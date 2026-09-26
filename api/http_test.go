package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	logger := zerolog.Nop()
	newHandler := func(publicHost string) http.Handler {
		return NewHTTPHandler(&API{Logger: &logger, Config: Config{PublicHost: publicHost}})
	}

	t.Run("html uses public host", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept", "text/html")
		rec := httptest.NewRecorder()
		newHandler("particles.example.com").ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Header().Get("Content-Type"), "text/html")
		require.Contains(t, rec.Body.String(), "websocat --binary wss://particles.example.com/ssh")
	})

	t.Run("plain text falls back to request host", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
		rec := httptest.NewRecorder()
		newHandler("").ServeHTTP(rec, req)

		require.Equal(t, "ssh -o ProxyCommand=\"websocat --binary ws://localhost:8080/ssh\" particles\n", rec.Body.String())
	})
}
