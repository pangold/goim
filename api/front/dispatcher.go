package front

import (
	"gitlab.com/pangold/goim/codec/protobuf"
)

// Handle if message requested via long term server
// This is fit if only one node(no cluster), no database.
// In this condition, sessions will just be stored in map simply.

type Dispatcher struct {
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (s *Dispatcher) Dispatch(msg *protobuf.Message) error {
	// TODO: to backend service
	// TODO: micro service rpc request backend service to upload message
	return nil
}