package system

import (
	api "gitlab.com/pangold/goim/api/grpc"
	"gitlab.com/pangold/goim/protocol"
)

type System struct {
	grpcServer *api.Server
}

func NewSystemMiddleware(g *api.Server) *System {
	return &System{
		grpcServer: g,
	}
}

func (s *System) Dispatch(msg *protocol.Message) {
	s.grpcServer.Dispatch(msg)
}

func (s *System) SessionIn(session *protocol.Session) {
	s.grpcServer.SessionIn(session)
}

func (s *System) SessionOut(session *protocol.Session) {
	s.grpcServer.SessionOut(session)
}
