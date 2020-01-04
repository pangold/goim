package api

import (
	"gitlab.com/pangold/goim/protocol"
)

type ImDispatchService struct {
	messageIn  chan *protocol.Message
	sessionIn  chan *Session
	sessionOut chan *Session
}

func NewImDispatchService() *ImDispatchService {
	return &ImDispatchService{
		messageIn:  make(chan *protocol.Message, 1024),
		sessionIn:  make(chan *Session, 1024),
		sessionOut: make(chan *Session, 1024),
	}
}

func (c *ImDispatchService) PutMessage(msg *protocol.Message) {
	c.messageIn <- msg
}

func (c *ImDispatchService) PutSessionIn(session *Session) {
	c.sessionIn <- session
}

func (c *ImDispatchService) PutSessionOut(session *Session) {
	c.sessionOut <- session
}

func (c *ImDispatchService) Dispatch(req *Empty, srv ImDispatchService_DispatchServer) error {
	for {
		select {
		case msg := <-c.messageIn:
			srv.Send(msg)
		}
	}
	return nil
}

func (c *ImDispatchService) SessionIn(req *Empty, srv ImDispatchService_SessionInServer) error {
	for {
		select {
		case msg := <-c.sessionIn:
			srv.Send(msg)
		}
	}
	return nil
}

func (c *ImDispatchService) SessionOut(req *Empty, srv ImDispatchService_SessionOutServer) error {
	for {
		select {
		case msg := <-c.sessionOut:
			srv.Send(msg)
		}
	}
	return nil
}