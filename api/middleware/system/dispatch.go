package system

import (
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// dispatch grpc service
type dispatchService struct {
	messageIn  chan *protocol.Message
	sessionIn  chan *protocol.Session
	sessionOut chan *protocol.Session
}

func newImDispatchService() *dispatchService {
	return &dispatchService{
		messageIn:  make(chan *protocol.Message, 1024),
		sessionIn:  make(chan *protocol.Session, 1024),
		sessionOut: make(chan *protocol.Session, 1024),
	}
}

func (c *dispatchService) Dispatch(req *protocol.Empty, srv protocol.ImDispatchService_DispatchServer) error {
	for {
		select {
		case msg := <-c.messageIn:
			srv.Send(msg)
		}
	}
	return nil
}

func (c *dispatchService) SessionIn(req *protocol.Empty, srv protocol.ImDispatchService_SessionInServer) error {
	for {
		select {
		case msg := <-c.sessionIn:
			srv.Send(msg)
		}
	}
	return nil
}

func (c *dispatchService) SessionOut(req *protocol.Empty, srv protocol.ImDispatchService_SessionOutServer) error {
	for {
		select {
		case msg := <-c.sessionOut:
			srv.Send(msg)
		}
	}
	return nil
}

// dispatch server
type DispatchServer struct {
	config      config.HostConfig
	dispatcher *dispatchService
}

func NewDispatchServer(conf config.HostConfig) *DispatchServer {
	dispatcher := &DispatchServer{
		config:     conf,
		dispatcher: newImDispatchService(),
	}
	go dispatcher.Run()
	return dispatcher
}

func (s *DispatchServer) Run() {
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		panic("failed to listen: %s" + err.Error())
	}
	log.Printf("grpc dispatch server start running on %s", s.config.Address)
	ss := grpc.NewServer()
	protocol.RegisterImDispatchServiceServer(ss, s.dispatcher)
	//
	reflection.Register(ss)
	if err := ss.Serve(listener); err != nil {
		panic("failed to serve" + err.Error())
	}
}

func (s *DispatchServer) Dispatch(msg *protocol.Message) {
	s.dispatcher.messageIn <- msg
}

func (s *DispatchServer) SessionIn(session *protocol.Session) {
	s.dispatcher.sessionIn <- session
}

func (s *DispatchServer) SessionOut(session *protocol.Session) {
	s.dispatcher.sessionOut <-session
}