package tcp

import (
	"fmt"
	"gitlab.com/pangold/goim/config"
	"log"
	"net"
	"time"
)

type Client struct {
	config config.TcpConfig
	conn net.Conn
	send chan []byte
	remaining []byte
}

func NewTcpClient(conf config.TcpConfig) *Client {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", conf.Address)
	c, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println("server is not starting")
	}
	client := &Client {
		config: conf,
		conn: c,
		send: make(chan []byte, 1024),
	}
	return client
}

func (c *Client) Run() {
	defer c.conn.Close()
	c.SendToken("f3ba0e819f55cb7")
	// go c.SendHeartbeatLoop()
	go c.SendLoop()
	go c.ReceiveLoop()
	for {
		c.SendMessage("input")
		time.Sleep(time.Second)
	}
}

func (c *Client) SendHeartbeatLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.SendHeartbeat()
		}
	}
}

func (c *Client) ReceiveLoop() {
	for {
		// FIXME: maxMessageSize not safe
		msg := make([]byte, maxMessageSize)
		if _, err := c.conn.Read(msg); err != nil {
			log.Printf("received error: %v", err)
			break
		}
		c.DispatchMessage(msg)
	}
}

func (c *Client) SendLoop() {
	ticker := time.NewTicker(pingPeriod)
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
		}
	}
}

func (c *Client) DispatchMessage(msg []byte) {
	// FIXME: Is here something wrong?
	msg = append(c.remaining, msg...)
	m, count := DeserializeInternalMessage(msg)
	if m != nil {
		c.HandleInternalMessage(m)
		c.remaining = msg[count:]
	} else {
		fmt.Println(string(msg))
		c.remaining = msg[len(msg):]
	}
}

func (c *Client) HandleInternalMessage(m *InternalMessage) {
	switch m.kind {
	case HEARTBEAT:
		// ReceiveLoop has PongWait detection
	case GOODBYE:
		//
	case TOKEN:
		//
	}
}

func (c *Client) SendHeartbeat() {
	// c.Send(SerializeHeartbeatMessage())
}

func (c *Client) SendToken(token string) {
	c.Send(SerializeTokenMessage([]byte(token)))
}

func (c *Client) Send(message []byte) {
	c.send <- message
}

func (c *Client) SendMessage(message string) {
	c.Send([]byte(message))
}
