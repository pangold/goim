package middleware

import (
	"gitlab.com/pangold/goim/protocol"
)

type Dispatcher interface {
	Dispatch(*protocol.Message)
}

type SyncSession interface {
	SessionIn(*protocol.Session)
	SessionOut(*protocol.Session)
}