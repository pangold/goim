package middleware

import (
	pb "gitlab.com/pangold/goim/api/grpc/proto"
	"gitlab.com/pangold/goim/protocol"
)

type Dispatcher interface {
	Dispatch(*protocol.Message)
}

type SyncSession interface {
	SessionIn(*pb.Session)
	SessionOut(*pb.Session)
}