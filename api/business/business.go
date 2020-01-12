package business

import (
	"gitlab.com/pangold/goim/protocol"
)

type Token interface {
	ExplainToken(token string) *protocol.Session
}

type Dispatcher interface {
	Dispatch(*protocol.Message)
}

type SyncSession interface {
	SessionIn(*protocol.Session)
	SessionOut(*protocol.Session)
}