package api

import (
	grpc "gitlab.com/pangold/goim/api/grpc"
	"gitlab.com/pangold/goim/api/http"
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/front"
	"gitlab.com/pangold/goim/protocol"
)

type ApiServer struct {
	front    *front.Server
	grpcServer *grpc.Server
	httpServer *http.Server
	sessions session.Sessions
}

func NewApiServer(conf config.Config) *ApiServer {
	api := &ApiServer{front: front.NewServer(conf)}
	api.front.SetConnectedHandler(api.handleConnection)
	api.front.SetDisconnectedHandler(api.handleDisconnection)
	api.front.SetMessageHandler(api.handleMessage)
	api.httpServer = http.NewServer(api.front, &api.sessions, conf.Http)
	api.grpcServer = grpc.NewServer(api.front, &api.sessions, conf.Grpc)
	return api
}

func (a *ApiServer) Run() {
	go a.grpcServer.Run()
	go a.httpServer.Run()
	a.front.Run()
}

// TODO: dispatch to backend service(cluster) to store in db/redis/etcd
// TODO: filter plugin if user id is invalid
func (a *ApiServer) handleConnection(token string) error {
	return a.sessions.Add(token, func(session *grpc.Session) error {
		// do filter
		a.grpcServer.Dispatcher.PutSessionIn(session)
		return nil
	})
}

// TODO: dispatch to backend service(cluster) to erase from db/redis/etcd
func (a *ApiServer) handleDisconnection(token string) {
	session := a.sessions.Remove(token)
	a.grpcServer.Dispatcher.PutSessionOut(session)
}

// TODO: dispatch to backend service
// TODO: micro service rpc request backend service to upload message
func (a *ApiServer) handleMessage(msg *protocol.Message, token string) error {
	// do dispatch
	a.grpcServer.Dispatcher.PutMessage(msg)
	return nil
}