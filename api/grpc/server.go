package im

import (
	"gitlab.com/pangold/goim/api/front"
	"gitlab.com/pangold/goim/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Server struct {
	config config.GrpcConfig
	front  *front.Server
}

func NewGrpcServer(front *front.Server, conf config.GrpcConfig) *Server {
	return &Server{
		config: conf,
		front:  front,
	}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		panic("failed to listen: %s" + err.Error())
	}
	log.Printf("grpc server start running, address: %s", s.config.Address)
	ss := grpc.NewServer()
	RegisterApiServer(ss, NewController(s.front))
	//
	reflection.Register(ss)
	if err := ss.Serve(listener); err != nil {
		panic("failed to serve" + err.Error())
	}
}