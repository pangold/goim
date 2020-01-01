package http

// For backend services
// Considering the security.

import (
	"gitlab.com/pangold/goim/api/front"
	"gitlab.com/pangold/goim/config"
)

type Server struct {
	router *Router
}

func NewServer(front *front.Server, conf config.HttpConfig) *Server {
	return &Server{
		router: NewRouter(front, conf),
	}
}

func (s *Server) Run() {
	s.router.Run()
}
