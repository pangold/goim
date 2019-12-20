package conn

import (
	"errors"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/conn/interfaces"
	"gitlab.com/pangold/goim/conn/tcp"
	"gitlab.com/pangold/goim/conn/ws"
)

type ChatServer struct {
	server               interfaces.Server
	connectedHandler    *func(string)
	disconnectedHandler *func(string)
	messageHandler      *func([]byte, string) error
	tokenHandler        *func(string) error
}

func NewChatServer(conf config.Config) *ChatServer {
	chat := &ChatServer{nil, nil, nil, nil, nil}
	if conf.Protocol == "tcp" {
		chat.server = tcp.NewTcpServer(conf.Tcp)
	} else if conf.Protocol == "ws" {
		chat.server = ws.NewWsServer(conf.Ws)
	} else {
		panic("unsupported protocol")
	}
	return chat
}

// must be invoked before SetConnectedHandler
func (c *ChatServer) SetMessageHandler(handler func([]byte, string) error) {
	c.messageHandler = &handler
}

// must be invoked before SetConnectedHandler
func (c *ChatServer) SetTokenHandler(handler func(string) error) {
	c.tokenHandler = &handler
}

func (c *ChatServer) SetConnectedHandler(handler func(string)) {
	// callback from pool is Connection type
	// callback out from here is string type(token)
	// that is the difference.
	c.connectedHandler = &handler
	c.server.GetPool().SetConnectedHandler(c.handleConnection)
}

func (c *ChatServer) SetDisconnectedHandler(handler func(string)) {
	c.disconnectedHandler = &handler
	c.server.GetPool().SetDisconnectedHandler(c.handleDisconnection)
}

func (c *ChatServer) handleConnection(connection interfaces.Connection) error {
	// when connected, ws has token, but tcp doesn't have
	// think about how to deal with it.
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
	connection.SetTokenHandler(c.tokenHandler) // tcp needs
	connection.SetMessageHandler(c.messageHandler)
	return nil
}

// remember: don't try to call Send or Receive relatives function here,
// maybe it's been closed. that's a risk.
func (c *ChatServer) handleDisconnection(connection interfaces.Connection) error {
	if c.disconnectedHandler != nil {
		(*c.disconnectedHandler)(connection.GetToken())
	}
	return nil
}

func (c *ChatServer) Run() {
	c.server.Run()
}

func (c *ChatServer) Send(token string, data []byte) error {
	connections := c.server.GetPool().GetConnections()
	if conn, ok := (*connections)[token]; ok {
		conn.Send(data)
		return nil
	}
	return errors.New("error: no such connection, token: " + token)
}

func (c *ChatServer) Broadcast(data []byte) (result []string) {
	connections := c.server.GetPool().GetConnections()
	for token, conn := range *connections {
		conn.Send(data)
		result = append(result, token)
	}
	return result
}

func (c *ChatServer) Remove(token string) bool {
	connections := c.server.GetPool().GetConnections()
	if conn, ok := (*connections)[token]; ok {
		conn.Stop()
		return true
	}
	return false
}

func (c *ChatServer) RemoveAll() (result []string) {
	connections := c.server.GetPool().GetConnections()
	for token, conn := range *connections {
		result = append(result, token)
		conn.Stop()
	}
	return result
}

func (c *ChatServer) CheckOnline(token string) bool {
	connections := c.server.GetPool().GetConnections()
	if _, ok := (*connections)[token]; ok {
		return true
	}
	return false
}

func (c *ChatServer) GetAll() (result []string) {
	connections := c.server.GetPool().GetConnections()
	for token, _ := range *connections {
		result = append(result, token)
	}
	return result
}