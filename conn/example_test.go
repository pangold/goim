package conn

import (
	"gitlab.com/pangold/goim/config"
	"testing"
)

var (
	server *Server
)

func init() {
	server = NewServer(config.Conf)
}

func TestServer_Run(t *testing.T) {
	server.Run()
}
