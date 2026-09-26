package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestWSConnRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		conn := newWSConn(c)
		defer func() { _ = conn.Close() }()
		_, _ = io.Copy(conn, conn) // echo
	}))
	defer srv.Close()

	c, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http"), nil)
	require.NoError(t, err)
	conn := newWSConn(c)
	defer func() { _ = conn.Close() }()

	// Text frames are ignored, binary frames are streamed.
	require.NoError(t, c.WriteMessage(websocket.TextMessage, []byte("ignored")))
	_, err = conn.Write([]byte("hello "))
	require.NoError(t, err)
	_, err = conn.Write([]byte("world"))
	require.NoError(t, err)

	buf := make([]byte, len("hello world"))
	_, err = io.ReadFull(conn, buf)
	require.NoError(t, err)
	require.Equal(t, "hello world", string(buf))
}
