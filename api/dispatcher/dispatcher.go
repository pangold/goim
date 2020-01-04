package dispatcher

import (
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/codec/protobuf"
)

// Business Logic Layer
//
// Handle if message requested via long connection server
// In this situation, relationship / group will just be stored in map simply.
// But data will be lost if shutdown.

type Dispatcher struct {

}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

// the reason why returns array, because of group.
// backend service can response asynchronizely(return empty array)
// or synchronizely(return non-empty array)
func (d *Dispatcher) Dispatch(msg *protobuf.Message) []*protobuf.Message {
	var res = make([]*protobuf.Message, 0)
	// TODO: dispatch to backend service
	// TODO: micro service rpc request backend service to upload message
	//
	// another way is process this message here(not recommended)
	return res
}

func (d *Dispatcher) SessionIn(s *session.Session) error {
	// TODO: dispatch to backend service(cluster) to store in db/redis/etcd
	// TODO: filter plugin if user id is invalid
	//
	// or a simple pub/sub to dispatch new coming session
	return nil
}

func (d *Dispatcher) SessionOut(token string) {
	// TODO: dispatch to backend service(cluster) to erase from db/redis/etcd
	// 
	// or a simple pub/sub to dispatch new coming session
}