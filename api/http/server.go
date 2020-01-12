package http

// For backend services
// Considering the security.

import (
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/front"
)

type Server struct {
	router *Router
}

func NewServer(front *front.Server, ss *session.Sessions, conf config.HostConfig) *Server {
	return &Server{
		router: NewRouter(front, ss, conf),
	}
}

func (s *Server) Run() {
	s.router.Run()
}
