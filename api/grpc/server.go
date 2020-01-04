package im

// For backend services
// Considering the security.

import (
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/front"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Server struct {
	config    config.GrpcConfig
	front    *front.Server
	sessions *session.Sessions

}

func NewGrpcServer(front *front.Server, ss *session.Sessions, conf config.GrpcConfig) *Server {
	return &Server{
		config:   conf,
		front:    front,
		sessions: ss,
	}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		panic("failed to listen: %s" + err.Error())
	}
	log.Printf("grpc server start running, address: %s", s.config.Address)
	ss := grpc.NewServer()
	RegisterApiServer(ss, NewController(s.front, s.sessions))
	//
	reflection.Register(ss)
	if err := ss.Serve(listener); err != nil {
		panic("failed to serve" + err.Error())
	}
}