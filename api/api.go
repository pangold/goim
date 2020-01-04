package api

import (
	im "gitlab.com/pangold/goim/api/grpc"
	"gitlab.com/pangold/goim/api/http"
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/codec/protobuf"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/front"
)

type Server interface {
	Run()
}

type ApiServer struct {
	front    *front.Server
	servers  []Server
	sessions session.Sessions
}

func NewApiServer(conf config.Config) *ApiServer {
	api := &ApiServer{front: front.NewServer(conf)}
	api.front.SetConnectedHandler(api.handleConnection)
	api.front.SetDisconnectedHandler(api.handleDisconnection)
	api.front.SetMessageHandler(api.handleMessage)
	api.servers = append(api.servers, http.NewServer(api.front, &api.sessions, conf.Http))
	api.servers = append(api.servers, im.NewGrpcServer(api.front, &api.sessions, conf.Grpc))
	return api
}

func (a *ApiServer) Run() {
	for _, server := range a.servers {
		go server.Run()
	}
	a.front.Run()
}

// TODO: dispatch to backend service(cluster) to store in db/redis/etcd
// TODO: filter plugin if user id is invalid
func (a *ApiServer) handleConnection(token string) error {
	return a.sessions.Add(token, func(session *session.Session) error {
		// do filter
		return nil
	})
}

// TODO: dispatch to backend service(cluster) to erase from db/redis/etcd
func (a *ApiServer) handleDisconnection(token string) {
	a.sessions.Remove(token)
}

// TODO: dispatch to backend service
// TODO: micro service rpc request backend service to upload message
func (a *ApiServer) handleMessage(msg *protobuf.Message, token string) error {
	// do dispatch
	return nil
}