package system

import (
	api "gitlab.com/pangold/goim/api/grpc"
	pb "gitlab.com/pangold/goim/api/grpc/proto"
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

func (s *System) SessionIn(session *pb.Session) {
	s.grpcServer.SessionIn(session)
}

func (s *System) SessionOut(session *pb.Session) {
	s.grpcServer.SessionOut(session)
}
