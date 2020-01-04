package api

// For backend services
// Considering the security.

import (
	pb "gitlab.com/pangold/goim/api/grpc/proto"
	"gitlab.com/pangold/goim/api/grpc/service"
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/front"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Server struct {
	config      config.GrpcConfig
	front      *front.Server
	sessions   *session.Sessions
	Dispatcher *service.ImDispatchService
}

func NewServer(front *front.Server, ss *session.Sessions, conf config.GrpcConfig) *Server {
	return &Server{
		config:     conf,
		front:      front,
		sessions:   ss,
		Dispatcher: service.NewImDispatchService(),
	}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		panic("failed to listen: %s" + err.Error())
	}
	log.Printf("grpc server start running, address: %s", s.config.Address)
	ss := grpc.NewServer()
	pb.RegisterImDispatchServiceServer(ss, s.Dispatcher)
	pb.RegisterImApiServiceServer(ss, service.NewImApiService(s.front, s.sessions))
	//
	reflection.Register(ss)
	if err := ss.Serve(listener); err != nil {
		panic("failed to serve" + err.Error())
	}
}