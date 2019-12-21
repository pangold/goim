package tcp

import (
	"fmt"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/conn/common"
	"log"
	"net"
	"time"
)

type Client struct {
	Connection
	config config.TcpConfig
}

func newTcpConnection(conf config.TcpConfig) net.Conn {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", conf.Address)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalln("server is not starting")
		return nil
	}
	return conn
}

func NewTcpClient(token string, conf config.TcpConfig) *Client {
	client := &Client {
		Connection: Connection {
			conn:      newTcpConnection(conf),
			send:      make(chan []byte, 1024),
			pool:      nil,
			token:     token,
			stopped:   false,
			remaining: nil,
		},
		config: conf,
	}
	client.Send(common.NewTokenMessage([]byte(token)).Serialize())
	go client.sendLoop()
	go client.receiveLoop()
	return client
}

func (c *Client) SendMessage(message string) {
	c.Send([]byte(message))
}

func (c *Client) handleMessage(message []byte) {
	c.remaining = append(c.remaining, message...)
	m, count := common.NewInternalMessage().Deserialize(c.remaining)
	if m != nil {
		c.handleInternalMessage(m)
		c.remaining = c.remaining[count:]
	}
	fmt.Println(string(c.remaining))
	c.remaining = nil
}

func (c *Client) receiveLoop() {
	defer c.Stop()
	for {
		msg := make([]byte, maxMessageSize)
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			// log.Println("tcp read error: timeout, does heartbeat work?")
			return
		}
		if _, err := c.conn.Read(msg); err != nil {
			// log.Println("tcp read error: ", err.Error())
			return
		}
		c.handleMessage(msg)
	}
}
