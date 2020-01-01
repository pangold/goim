package front

import (
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/codec/protobuf"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/conn"
	"log"
)

type Server struct {
	conn       *conn.Server
	// FIXME: what if I want another codec
	codec      *protobuf.Codec
	sessions   *Sessions
	dispatcher *Dispatcher
}

func NewServer(conf config.Config) *Server {
	s := &Server{
		conn:       conn.NewServer(conf),
		// FIXME: what if I want another codec
		codec:      protobuf.NewCodec(),
		sessions:   NewSessions(),
		dispatcher: NewDispatcher(),
	}
	s.codec.SetEncodeHandler(s.handleEncode)
	s.codec.SetDecodeHandler(s.handleDecode)
	s.conn.SetMessageHandler(s.handleReceived)
	s.conn.SetConnectedHandler(s.handleConnected)
	s.conn.SetDisconnectedHandler(s.handleDisconnected)
	return s
}

func (s *Server) handleFilter(ss *Session) error {
	// TODO: filter plugin
	// TODO: micro service rpc request backend service to check if user id is valid
	return nil
}

func (s *Server) handleConnected(token string) error {
	if err := s.sessions.Add(token, s.handleFilter); err != nil {
		return err
	}
	return nil
}

func (s *Server) handleDisconnected(token string) {
	s.sessions.Remove(token)
}

func (s *Server) handleReceived(data []byte, token string) error {
	s.codec.Decode(token, data)
	return nil
}

func (s *Server) handleDecode(token interface{}, msg *protobuf.Message) {
	if err := s.dispatcher.Dispatch(msg); err != nil {
		// Nothing needs to be responded
		// s.codec.Encode(token, res)
		log.Printf(err.Error())
	}
}

func (s *Server) handleEncode(token interface{}, data []byte) {
	if err := s.conn.Send(token.(string), data); err != nil {
		log.Printf(err.Error())
	}
}

// for internal dispatcher if we need it
func (s *Server) SendEx(token string, msg *protobuf.Message) {
	if err := s.codec.Encode(token, msg); err != nil {
		log.Printf(err.Error())
	}
}

func (s *Server) BroadcastEx(msg *protobuf.Message) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("unabled to broadcast, error: %v", err)
		return
	}
	s.conn.Broadcast(data)
}

// for api server
func (s *Server) Send(token string, data []byte) {
	if err := s.conn.Send(token, data); err != nil {
		log.Printf(err.Error())
	}
}

func (s *Server) Broadcast(data []byte) {
	s.conn.Broadcast(data)
}

func (s *Server) Remove(token string) {
	s.conn.Remove(token)
	s.sessions.Remove(token)
}

func (s *Server) RemoveAll() {
	s.conn.RemoveAll()
	s.sessions.Clear()
}

func (s *Server) GetOnlineUserIds() []string {
	return s.sessions.GetUserIds()
}

func (s *Server) GetOnlineTokenByUserId(uid string) string {
	return s.sessions.GetTokenByUserId(uid)
}

func (s *Server) Run() {
	s.conn.Run()
}