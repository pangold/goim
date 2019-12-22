package conn

import (
	"errors"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/conn/common"
	"gitlab.com/pangold/goim/conn/interfaces"
	"gitlab.com/pangold/goim/conn/tcp"
	"gitlab.com/pangold/goim/conn/websocket"
)

type Server struct {
	servers              []interfaces.Server
	pool                 interfaces.Pool
	connectedHandler    *func(string)
	disconnectedHandler *func(string)
	messageHandler      *func([]byte, string) error
	tokenHandler        *func(string) error
}

//
func NewServer(conf config.Config) *Server {
	c := &Server{ pool: common.NewPool() }
	for _, proto := range conf.Protocols {
		if proto == "tcp" {
			c.servers = append(c.servers, tcp.NewTcpServer(c.pool, conf.Tcp))
		} else if proto == "ws" {
			c.servers = append(c.servers, websocket.NewWsServer(c.pool, conf.Ws))
		} else if proto == "wss" {
			c.servers = append(c.servers, websocket.NewWsServer(c.pool, conf.Wss))
		} else {
			panic("unsupported protocol")
		}
	}
	return c
}

// must be invoked before SetConnectedHandler
func (c *Server) SetMessageHandler(handler func([]byte, string) error) {
	c.messageHandler = &handler
}

// must be invoked before SetConnectedHandler
func (c *Server) SetTokenHandler(handler func(string) error) {
	c.tokenHandler = &handler
}

func (c *Server) SetConnectedHandler(handler func(string)) {
	// callback from pool is Connection type
	// callback out from here is string type(token)
	// that is the difference.
	c.connectedHandler = &handler
	c.pool.SetConnectedHandler(c.handleConnection)
}

func (c *Server) SetDisconnectedHandler(handler func(string)) {
	c.disconnectedHandler = &handler
	c.pool.SetDisconnectedHandler(c.handleDisconnection)
}

func (c *Server) handleConnection(connection interfaces.Conn) error {
	token := connection.GetToken()
	if c.tokenHandler != nil && len(token) > 0 {
		// check if it is valid token
		if err := (*c.tokenHandler)(token); err != nil {
			return err
		}
	}
	if len(token) > 0 && c.connectedHandler != nil {
		(*c.connectedHandler)(token)
	}
	connection.SetMessageHandler(c.messageHandler)
	return nil
}

// remember: don't try to call Send or Receive relatives function here,
// maybe it's been closed. that's a risk.
func (c *Server) handleDisconnection(connection interfaces.Conn) error {
	if c.disconnectedHandler != nil {
		(*c.disconnectedHandler)(connection.GetToken())
	}
	return nil
}

func (c *Server) Run() {
	if len(c.servers) == 0 {
		panic("no protocol is being specified")
	}
	for _, server := range c.servers[1:] {
		go server.Run()
	}
	c.servers[0].Run()
}

func (c *Server) Send(token string, data []byte) error {
	connections := c.pool.GetConnections()
	if conn, ok := (*connections)[token]; ok {
		conn.Send(data)
		return nil
	}
	return errors.New("error: no such connection, token: " + token)
}

func (c *Server) Broadcast(data []byte) (result []string) {
	connections := c.pool.GetConnections()
	for token, conn := range *connections {
		conn.Send(data)
		result = append(result, token)
	}
	return result
}

func (c *Server) Remove(token string) bool {
	connections := c.pool.GetConnections()
	if conn, ok := (*connections)[token]; ok {
		conn.Stop()
		return true
	}
	return false
}

func (c *Server) RemoveAll() (result []string) {
	connections := c.pool.GetConnections()
	for token, conn := range *connections {
		result = append(result, token)
		conn.Stop()
	}
	return result
}

func (c *Server) CheckOnline(token string) bool {
	connections := c.pool.GetConnections()
	if _, ok := (*connections)[token]; ok {
		return true
	}
	return false
}

func (c *Server) GetAll() (result []string) {
	connections := c.pool.GetConnections()
	for token, _ := range *connections {
		result = append(result, token)
	}
	return result
}