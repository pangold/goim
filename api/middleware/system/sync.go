package system

import (
	api "gitlab.com/pangold/goim/api/grpc"
	"gitlab.com/pangold/goim/protocol"
)

type Sync struct {
	grpcServer *api.Server
}

func NewSync(g *api.Server) *Sync {
	return &Sync{
		grpcServer: g,
	}
}

func (s *Sync) Dispatch(msg *protocol.Message) {
	s.grpcServer.Dispatch(msg)
}

func (s *Sync) SessionIn(session *protocol.Session) {
	s.grpcServer.SessionIn(session)
}

func (s *Sync) SessionOut(session *protocol.Session) {
	s.grpcServer.SessionOut(session)
}
