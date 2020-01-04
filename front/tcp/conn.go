package tcp

import (
	"errors"
	"gitlab.com/pangold/goim/front/interfaces"
	"log"
	"net"
	"time"
)

type Connection struct {
	conn            net.Conn
	messageHandler  func([]byte, interface{}) error
	codec           interface{}
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

func (c *Connection) BindCodec(codec interface{}) {
	c.codec = codec
}

func (c *Connection) GetCodec() interface{} {
	return c.codec
}

func (c *Connection) SetMessageHandler(handler func([]byte, interface{}) error) {
	c.messageHandler = handler
}

func (c *Connection) Stop() {
	// not close directly. but close by SendLoop when exit loop
	if !c.stopped {
		close(c.send)
		c.stopped = true
	}
}

func (c *Connection) Close() {
	if err := c.conn.Close(); err != nil {
		log.Printf("tcp connection close error: %v", err)
	}
}

func (c *Connection) GetToken() string {
	return c.token
}

func (c *Connection) Send(message []byte) {
	if !c.stopped {
		c.send <- message
	}
}

func (c *Connection) sendLoop() {
	ticker := time.NewTicker(pingPeriod * 10)
	defer c.conn.Close()
	defer ticker.Stop()
	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				// log.Printf("write error: timeout")
				return
			}
			if !ok {
				c.conn.Write(NewGoodbyeMessage().Serialize())
				// log.Printf("tcp connection say goodbye")
				return
			}
			if _, err := c.conn.Write(message); err != nil {
				// log.Printf("tcp write error: %v", err)
				return
			}
		case <-ticker.C:
			if _, err := c.conn.Write(NewHeartbeatMessage().Serialize()); err != nil {
				// log.Printf("tcp send hearbeat error: %v", err)
				return
			}
		}
	}
}

func (c *Connection) handleInternalMessage(m *InternalMessage) {
	switch m.Kind {
	case HEARTBEAT:
		// ReceiveLoop has PongWait detection
		log.Println("heart beat.")
	case GOODBYE:
		log.Println("client say goodbye")
	case TOKEN:
		if c.token == "" {
			c.token = string(m.Body)
			c.pool.Register(c)
		} else {
			log.Printf("error: unexpected token request, original: %s, now: %s", c.token, string(m.Body))
		}
	}
}

// callback message(normal message)
func (c *Connection) handleMessage(msg []byte) error {
	c.remaining = append(c.remaining, msg...)
	m, count := NewInternalMessage().Deserialize(c.remaining)
	if m != nil {
		c.handleInternalMessage(m)
		c.remaining = c.remaining[count:]
	}
	// token must be requested at the first time.
	if c.token == "" {
		// no token, no requests will be rightful.
		return errors.New("unauthorized request")
	}
	// message callback handler for invokers
	if err := c.messageHandler(c.remaining, c); err != nil {
		return errors.New("unexpected data")
	}
	// FIXME: with echo reply, receive loop would not exit.
	// c.Send(c.remaining) // temp
	c.remaining = nil
	return nil
}

func (c *Connection) receiveLoop() {
	defer c.pool.Unregister(c)
	for {
		msg := make([]byte, maxMessageSize)
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			// log.Printf("tcp read error: %v", err)
			return
		}
		if _, err := c.conn.Read(msg); err != nil {
			// log.Printf("tcp read error: %v", err)
			return
		}
		if err := c.handleMessage(msg); err != nil {
			log.Printf("tcp read dispatch error: %v", err)
			return
		}
	}
}

