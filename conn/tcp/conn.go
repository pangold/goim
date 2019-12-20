package tcp

import (
	"github.com/gorilla/websocket"
	"log"
	"net"
	"time"
)

type Conn struct {
	server   *Server
	conn      net.Conn
	received *func([]byte, string) error
	send      chan []byte
	token     string
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func (c *Conn) sendHeartbeat() error {
	return nil
}

func (c *Conn) receiveHeatbeat() error {

	return nil
}

func (c *Conn) sendLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.conn.Close()
		ticker.Stop()
	}()
	//
	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if !ok {
				c.conn.Write([]byte{})
				return
			}
			if _, err := c.conn.Write(message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.sendHeartbeat(); err != nil {
				return
			}
		}
	}
}

func (c *Conn) receiveLoop() {
	defer func() {
		c.server.unregister <- c
	}()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	for {
		msg := make([]byte, 0)
		if _, err := c.conn.Read(msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Fatalf("error: %v", err)
			}
			break
		}
		if c.received != nil {
			if err := (*c.received)(msg, c.token); err != nil {
				log.Fatalf("error: unexpected data")
			}
		}
		// temp
		c.Send(msg)
	}
}

func (c *Conn) Close() {
	if err := c.conn.Close(); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func (c *Conn) Send(data []byte) {
	c.send <- data
}

