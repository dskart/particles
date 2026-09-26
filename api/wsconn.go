package api

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// wsConn adapts a gorilla websocket connection to a net.Conn so it can be handed
// to the SSH server. SSH bytes are carried in binary messages.
type wsConn struct {
	*websocket.Conn
	r       io.Reader
	writeMu sync.Mutex
}

var _ net.Conn = (*wsConn)(nil)

func newWSConn(c *websocket.Conn) *wsConn {
	return &wsConn{Conn: c}
}

func (c *wsConn) Read(b []byte) (int, error) {
	for {
		if c.r == nil {
			msgType, r, err := c.NextReader()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					return 0, io.EOF
				}
				return 0, err
			}
			if msgType != websocket.BinaryMessage {
				continue
			}
			c.r = r
		}

		n, err := c.r.Read(b)
		if err == io.EOF {
			c.r = nil
			if n > 0 {
				return n, nil
			}
			continue
		}
		return n, err
	}
}

func (c *wsConn) Write(b []byte) (int, error) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if err := c.WriteMessage(websocket.BinaryMessage, b); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *wsConn) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}
