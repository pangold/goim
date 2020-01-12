package api

import (
	"gitlab.com/pangold/goim/api/business"
	"gitlab.com/pangold/goim/api/business/system"
	grpc "gitlab.com/pangold/goim/api/grpc"
	"gitlab.com/pangold/goim/api/http"
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/front"
	"gitlab.com/pangold/goim/protocol"
)

type ApiServer struct {
	front       *front.Server
	grpcServer  *grpc.Server
	httpServer  *http.Server
	sessions    *session.Sessions
	syncSession business.SyncSession
	dispatcher  business.Dispatcher
}

func NewApiServer(conf config.Config) *ApiServer {
	api := &ApiServer{
		front: front.NewServer(conf.Front),
		sessions: session.NewSessions(system.NewToken(conf.Token.SecretKey)),
	}
	// handle new coming connection
	api.front.SetConnectedHandler(api.handleConnection)
	// handle new disconnected connection
	api.front.SetDisconnectedHandler(api.handleDisconnection)
	// handle received message(from front server)
	// to chat with others.
	api.front.SetMessageHandler(api.handleMessage)
	// Http api, designs for single node
	api.httpServer = http.NewServer(api.front, api.sessions, conf.Back.Http)
	// Grpc api, designs for cluster
	api.grpcServer = grpc.NewServer(api.front, api.sessions, conf.Back.Grpc)
	// Default middleware for dispatching message/session
	dispatcher := system.NewDispatchServer(conf.Back.Dispatch)
	// Default
	// We use grpc to dispatch received message to backend services.
	// You can also custom your own middleware to dispatch message to:
	// Your own service, MQ / Redis / DB, ignore them, or others
	api.dispatcher = dispatcher
	// Default
	api.syncSession = dispatcher
	//
	return api
}

func (a *ApiServer) Run() {
	go a.grpcServer.Run()
	go a.httpServer.Run()
	a.front.Run()
}

func (a *ApiServer) ResetDispatcher(dis business.Dispatcher) {
	a.dispatcher = dis
}

func (a *ApiServer) ResetSyncSession(ses business.SyncSession) {
	a.syncSession = ses
}

func (a *ApiServer) ResetToken(token business.Token) {
	a.sessions.ResetTokenExplainer(token)
}

func (a *ApiServer) handleConnection(token string) error {
	return a.sessions.Add(token, func(session *protocol.Session) error {
		if a.syncSession != nil {
			a.syncSession.SessionIn(session)
		}
		return nil
	})
}

func (a *ApiServer) handleDisconnection(token string) {
	session := a.sessions.Remove(token)
	if a.syncSession != nil {
		a.syncSession.SessionOut(session)
	}
}

func (a *ApiServer) handleMessage(msg *protocol.Message, token string) error {
	if a.dispatcher != nil {
		a.dispatcher.Dispatch(msg)
	}
	return nil
}