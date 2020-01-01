package api

import (
	"gitlab.com/pangold/goim/api/front"
	im "gitlab.com/pangold/goim/api/grpc"
	"gitlab.com/pangold/goim/api/http"
	"gitlab.com/pangold/goim/config"
)

type ApiServer struct {
	frontServer *front.Server
	httpServer  *http.Server
	grpcServer  *im.Server
}

func NewApiServer(conf config.Config) *ApiServer {
	api := &ApiServer{frontServer: front.NewServer(conf)}
	api.httpServer = http.NewServer(api.frontServer, conf.Http)
	api.grpcServer = im.NewGrpcServer(api.frontServer, conf.Grpc)
	return api
}

func (a *ApiServer) Run() {
	go a.httpServer.Run()
	go a.grpcServer.Run()
	a.frontServer.Run()
}