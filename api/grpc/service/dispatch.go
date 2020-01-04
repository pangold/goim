package service

import (
	"gitlab.com/pangold/goim/protocol"
)

type ImDispatchService struct {
	messageIn  chan *protocol.Message
	sessionIn  chan *protocol.Session
	sessionOut chan *protocol.Session
}

func NewImDispatchService() *ImDispatchService {
	return &ImDispatchService{
		messageIn:  make(chan *protocol.Message, 1024),
		sessionIn:  make(chan *protocol.Session, 1024),
		sessionOut: make(chan *protocol.Session, 1024),
	}
}

func (c *ImDispatchService) PutMessage(msg *protocol.Message) {
	c.messageIn <- msg
}

func (c *ImDispatchService) PutSessionIn(session *protocol.Session) {
	c.sessionIn <- session
}

func (c *ImDispatchService) PutSessionOut(session *protocol.Session) {
	c.sessionOut <- session
}

func (c *ImDispatchService) Dispatch(req *protocol.Empty, srv protocol.ImDispatchService_DispatchServer) error {
	for {
		select {
		case msg := <-c.messageIn:
			srv.Send(msg)
		}
	}
	return nil
}

func (c *ImDispatchService) SessionIn(req *protocol.Empty, srv protocol.ImDispatchService_SessionInServer) error {
	for {
		select {
		case msg := <-c.sessionIn:
			srv.Send(msg)
		}
	}
	return nil
}

func (c *ImDispatchService) SessionOut(req *protocol.Empty, srv protocol.ImDispatchService_SessionOutServer) error {
	for {
		select {
		case msg := <-c.sessionOut:
			srv.Send(msg)
		}
	}
	return nil
}