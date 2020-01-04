package api

import (
	grpc "gitlab.com/pangold/goim/api/grpc"
	"gitlab.com/pangold/goim/api/http"
	"gitlab.com/pangold/goim/api/middleware"
	"gitlab.com/pangold/goim/api/middleware/system"
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
	syncSession  middleware.SyncSession
	dispatcher   middleware.Dispatcher
}

func NewApiServer(conf config.Config) *ApiServer {
	api := &ApiServer{
		front: front.NewServer(conf),
		sessions: session.NewSessions(),
	}
	// handle new coming connection
	api.front.SetConnectedHandler(api.handleConnection)
	// handle new disconnected connection
	api.front.SetDisconnectedHandler(api.handleDisconnection)
	// handle received message(from front server)
	// to chat with others.
	api.front.SetMessageHandler(api.handleMessage)
	// Http api, designs for single node
	api.httpServer = http.NewServer(api.front, api.sessions, conf.Http)
	// Grpc api, designs for cluster
	api.grpcServer = grpc.NewServer(api.front, api.sessions, conf.Grpc)
	// Default middleware for dispatching message/session
	sm := system.NewSystemMiddleware(api.grpcServer)
	// Default
	// We use grpc to dispatch received message to backend services.
	// You can also custom your own middleware to dispatch message to:
	// Your own service, MQ / Redis / DB, ignore them, or others
	api.dispatcher = sm
	// Default
	// We just ignore session synchronization
	// api.syncSession = sm
	return api
}

func (a *ApiServer) Run() {
	go a.grpcServer.Run()
	go a.httpServer.Run()
	a.front.Run()
}

func (a *ApiServer) ResetDispatcher(dis middleware.Dispatcher) {
	a.dispatcher = dis
}

func (a *ApiServer) ResetSyncSession(ses middleware.SyncSession) {
	a.syncSession = ses
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