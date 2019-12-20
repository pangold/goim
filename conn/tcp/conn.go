package tcp

import (
	"github.com/gorilla/websocket"
	"gitlab.com/pangold/goim/conn/interfaces"
	"log"
	"net"
	"time"
)

type Connection struct {
	conn            net.Conn
	messageHandler *func([]byte, string) error
	pool            interfaces.Pool
	send            chan []byte
	token           string
	stopped         bool
	remaining       []byte
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


func (c *Connection) SetMessageHandler(handler *func([]byte, string) error) {
	c.messageHandler = handler
}

func (c *Connection) SetTokenHandler(handler *func(string) error) {
}

func (c *Connection) Stop() {
	// not close directly. but close by SendLoop when exit loop
	if !c.stopped {
		close(c.send)
		c.stopped = true
	}
}

func (c *Connection) GetToken() string {
	return c.token
}

func (c *Connection) Send(message []byte) {
	c.send <- message
}

// tcp send ping message
func (c *Connection) sendHeartbeat() error {
	//c.send <- SerializeHeartbeatMessage()
	return nil
}

func (c *Connection) handleInternalMessage(m *InternalMessage) {
	switch m.kind {
	case HEARTBEAT:
		// ReceiveLoop has PongWait detection
	case GOODBYE:
		log.Println("say goodbye")
	case TOKEN:
		if len(c.token) == 0 {
			c.token = string(m.body)
			c.pool.Register(c)
		} else {
			log.Fatalf("error: unexpected token request, original: %s, now: %s", c.token, string(m.body))
		}
	}
}

// callback message(normal message)
func (c *Connection) dispatchMessage(msg []byte) {
	// TODO: what if bytes remaining?
	// FIXME: Is here something wrong?
	msg = append(c.remaining, msg...)
	m, count := DeserializeInternalMessage(msg)
	if m != nil {
		c.handleInternalMessage(m)
		c.remaining = msg[count:]
	} else if len(c.token) == 0 {
		// no token, no requests will be rightful.
		c.Stop()
		log.Fatalf("error: unauthorized request")
	} else if c.messageHandler != nil {
		if err := (*c.messageHandler)(msg, c.token); err != nil {
			log.Fatalf("error: unexpected data")
		}
		c.remaining = msg[len(msg):]
	} else {
		log.Println(string(msg))
		c.remaining = msg[len(msg):]
	}
}

// 3 extra types of message: Heartbeat, Goodbye, Token
func (c *Connection) sendLoop() {
	ticker := time.NewTicker(pingPeriod * 10)
	defer c.conn.Close()
	defer ticker.Stop()
	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if !ok {
				c.conn.Write(SerializeGoodbyeMessage())
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

func (c *Connection) receiveLoop() {
	defer c.pool.Unregister(c)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	for {
		msg := make([]byte, maxMessageSize)
		if _, err := c.conn.Read(msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Fatalf("error: %v", err)
			}
			break
		}
		c.dispatchMessage(msg)
		c.Send(msg) // temp
	}
}

