package service

import (
	pb "gitlab.com/pangold/goim/api/grpc/proto"
	"gitlab.com/pangold/goim/protocol"
)

type ImDispatchService struct {
	messageIn  chan *protocol.Message
	sessionIn  chan *pb.Session
	sessionOut chan *pb.Session
}

func NewImDispatchService() *ImDispatchService {
	return &ImDispatchService{
		messageIn:  make(chan *protocol.Message, 1024),
		sessionIn:  make(chan *pb.Session, 1024),
		sessionOut: make(chan *pb.Session, 1024),
	}
}

func (c *ImDispatchService) PutMessage(msg *protocol.Message) {
	c.messageIn <- msg
}

func (c *ImDispatchService) PutSessionIn(session *pb.Session) {
	c.sessionIn <- session
}

func (c *ImDispatchService) PutSessionOut(session *pb.Session) {
	c.sessionOut <- session
}

func (c *ImDispatchService) Dispatch(req *pb.Empty, srv pb.ImDispatchService_DispatchServer) error {
	for {
		select {
		case msg := <-c.messageIn:
			srv.Send(msg)
		}
	}
	return nil
}

func (c *ImDispatchService) SessionIn(req *pb.Empty, srv pb.ImDispatchService_SessionInServer) error {
	for {
		select {
		case msg := <-c.sessionIn:
			srv.Send(msg)
		}
	}
	return nil
}

func (c *ImDispatchService) SessionOut(req *pb.Empty, srv pb.ImDispatchService_SessionOutServer) error {
	for {
		select {
		case msg := <-c.sessionOut:
			srv.Send(msg)
		}
	}
	return nil
}