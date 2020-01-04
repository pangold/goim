package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/front"
	"log"
)

type Router struct {
	config      config.HttpConfig
	router     *gin.Engine
	controller *Controller
}

func NewRouter(front *front.Server, ss *session.Sessions, conf config.HttpConfig) *Router {
	r := &Router{
		config:     conf,
		router:     gin.Default(),
		controller: NewController(front, ss),
	}
	basicRouter(r.router, r.controller)
	return r
}

func (r *Router) Run() {
	log.Printf("Http backend service start running %s", r.config.Address)
	if err := r.router.Run(r.config.Address); err != nil {
		panic(err)
	}
}

func basicRouter(router *gin.Engine, basic *Controller) {
	b := router.Group("/api")
	b.GET ("/", basic.List)
	b.POST("/send", basic.Send)
	b.POST("/broadcast", basic.Broadcast)
	b.POST("/online", basic.Online)
	b.POST("/kick", basic.Kick)
}

//func friendRouter(router *gin.Engine, friend *controller.Friend) {
//	f := router.Group("/api/friend")
//	f.GET ("/", friend.List)
//	f.POST("/make", friend.Make)
//	f.POST("/accept", friend.Accept)
//	f.POST("/reject", friend.Reject)
//	f.POST("/breakup", friend.Breakup)
//	f.POST("/recommend", friend.Recommend)
//}
//
//func groupRouter(router *gin.Engine, group *controller.Group) {
//	g := router.Group("/api/group")
//	g.GET ("/", group.List)
//	g.POST("/create", group.Create)
//}
