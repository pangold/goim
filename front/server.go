package front

import (
	"errors"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/front/codec"
	"gitlab.com/pangold/goim/front/interfaces"
	"gitlab.com/pangold/goim/front/pool"
	"gitlab.com/pangold/goim/front/tcp"
	"gitlab.com/pangold/goim/front/websocket"
	"gitlab.com/pangold/goim/protocol"
)

type Server struct {
	servers              []interfaces.Server
	pool                 interfaces.Pool
	connectedHandler    *func(string) error
	disconnectedHandler *func(string)
	messageHandler      *func(*protocol.Message, string) error
}

func NewServer(conf config.FrontConfig) *Server {
	s := &Server{ nil, pool.NewPool(), nil, nil, nil }
	for _, proto := range conf.Protocols {
		if proto == "tcp" {
			s.servers = append(s.servers, tcp.NewTcpServer(s.pool, conf.Tcp))
		} else if proto == "ws" {
			s.servers = append(s.servers, websocket.NewWsServer(s.pool, conf.Ws))
		} else if proto == "wss" {
			s.servers = append(s.servers, websocket.NewWsServer(s.pool, conf.Wss))
		} else {
			panic("unsupported protocol")
		}
	}
	s.pool.SetConnectedHandler(s.handleConnection)
	s.pool.SetDisconnectedHandler(s.handleDisconnection)
	return s
}

func (s *Server) SetMessageHandler(handler func(*protocol.Message, string) error) {
	s.messageHandler = &handler
}

func (s *Server) SetConnectedHandler(handler func(string) error) {
	s.connectedHandler = &handler
}

func (s *Server) SetDisconnectedHandler(handler func(string)) {
	s.disconnectedHandler = &handler
}

func (s *Server) handleConnection(conn interfaces.Conn) error {
	token := conn.GetToken()
	if token == "" {
		return errors.New("token could not be empty")
	}
	if s.connectedHandler != nil {
		// check token if it is valid
		if err := (*s.connectedHandler)(token); err != nil {
			return err
		}
	}
	c := codec.NewCodec()
	c.SetDecodeHandler(s.handleDecode)
	conn.BindCodec(c)
	conn.SetMessageHandler(s.handleMessage)
	return nil
}

// remember: don't try to call Send or Receive relatives function here,
// maybe it's been closed. that's a risk.
func (s *Server) handleDisconnection(conn interfaces.Conn) error {
	if s.disconnectedHandler != nil {
		(*s.disconnectedHandler)(conn.GetToken())
	}
	return nil
}

func (s *Server) handleMessage(data []byte, conn interface{}) error {
	c := conn.(interfaces.Conn)
	c.GetCodec().(*codec.Codec).Decode(c, data)
	return nil
}

func (s *Server) handleDecode(conn interfaces.Conn, msg *protocol.Message) {
	// callback message
	if s.messageHandler != nil {
		(*s.messageHandler)(msg, conn.GetToken())
	}
}

func (s *Server) Send(token string, msg *protocol.Message) error {
	connections := s.pool.GetConnections()
	if conn, ok := (*connections)[token]; ok {
		conn.GetCodec().(*codec.Codec).Send(conn, msg)
		return nil
	}
	return errors.New("error: no such connection, token: " + token)
}

func (s *Server) Broadcast(msg *protocol.Message) (result []string) {
	connections := s.pool.GetConnections()
	for token, conn := range *connections {
		conn.GetCodec().(*codec.Codec).Send(conn, msg)
		result = append(result, token)
	}
	return result
}

func (s *Server) Remove(token string) bool {
	connections := s.pool.GetConnections()
	if conn, ok := (*connections)[token]; ok {
		conn.Stop()
		return true
	}
	return false
}

func (s *Server) RemoveAll() (result []string) {
	connections := s.pool.GetConnections()
	for token, conn := range *connections {
		result = append(result, token)
		conn.Stop()
	}
	return result
}

func (s *Server) CheckOnline(token string) bool {
	connections := s.pool.GetConnections()
	if _, ok := (*connections)[token]; ok {
		return true
	}
	return false
}

func (s *Server) GetAll() (result []string) {
	connections := s.pool.GetConnections()
	for token, _ := range *connections {
		result = append(result, token)
	}
	return result
}

func (s *Server) Run() {
	if len(s.servers) == 0 {
		panic("no protocol is being specified")
	}
	for _, server := range s.servers[1:] {
		go server.Run()
	}
	s.servers[0].Run()
}