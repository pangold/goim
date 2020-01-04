package api

// For backend services
// Considering the security.

import (
	pb "gitlab.com/pangold/goim/api/grpc/proto"
	"gitlab.com/pangold/goim/api/grpc/service"
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/front"
	"gitlab.com/pangold/goim/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Server struct {
	config      config.GrpcConfig
	front      *front.Server
	sessions   *session.Sessions
	dispatcher *service.ImDispatchService
}

func NewServer(front *front.Server, ss *session.Sessions, conf config.GrpcConfig) *Server {
	return &Server{
		config:     conf,
		front:      front,
		sessions:   ss,
		dispatcher: service.NewImDispatchService(),
	}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		panic("failed to listen: %s" + err.Error())
	}
	log.Printf("grpc server start running, address: %s", s.config.Address)
	ss := grpc.NewServer()
	pb.RegisterImDispatchServiceServer(ss, s.dispatcher)
	pb.RegisterImApiServiceServer(ss, service.NewImApiService(s.front, s.sessions))
	//
	reflection.Register(ss)
	if err := ss.Serve(listener); err != nil {
		panic("failed to serve" + err.Error())
	}
}

func (s *Server) Dispatch(msg *protocol.Message) {
	s.dispatcher.PutMessage(msg)
}

func (s *Server) SessionIn(session *pb.Session) {
	s.dispatcher.PutSessionIn(session)
}

func (s *Server) SessionOut(session *pb.Session) {
	s.dispatcher.PutSessionOut(session)
}