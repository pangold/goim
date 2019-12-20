package ws

import (
	"github.com/gorilla/websocket"
	"gitlab.com/pangold/goim/conn/interfaces"
	"log"
	"time"
)

type Connection struct {
	conn            *websocket.Conn
	messageHandler  *func([]byte, string) error
	pool             interfaces.Pool
	send             chan []byte
	token            string
	stopped          bool
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
	// ws doesn't need token handler, but tcp does
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

func (c *Connection) SyncSend(message []byte) error {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	if _, err := w.Write(message); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

func (c *Connection) SyncSendHeartbeat() error {
	if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}
	if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
		return err
	}
	return nil
}

func (c *Connection) ReceiveHeartbeat(string) error {
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Fatalf("receive heartbeat error: %v", err)
		return err
	}
	return nil
}

func (c *Connection) DispatchMessage(msg []byte) {
	if c.messageHandler != nil {
		if err := (*c.messageHandler)(msg, c.token); err != nil {
			log.Fatalf("error: unexpected data")
		}
	}
}

func (c *Connection) SendLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer c.conn.Close()
	defer ticker.Stop()
	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Fatalf("write deadline error: %v", err)
				return
			}
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.SyncSend(message); err != nil {
				log.Fatalf("sync send error: %v", err)
				return
			}
		case <-ticker.C:
			if err := c.SyncSendHeartbeat(); err != nil {
				log.Fatalf("sync send heartbeat error: %v", err)
				return
			}
		}
	}
}

// 处理被动关闭连接，如客户端关闭、或者其他错误
func (c *Connection) ReceiveLoop() {
	defer c.pool.Unregister(c)
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(c.ReceiveHeartbeat)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Fatalf("unexpected close error: %v", err)
			}
			return
		}
		c.DispatchMessage(msg)
		c.Send(msg) // temp
	}
}