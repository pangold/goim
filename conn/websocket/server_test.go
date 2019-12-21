package websocket

import (
	"gitlab.com/pangold/goim/config"
	"testing"
)

func TestServer_Run(t *testing.T) {
	s := NewWsServer(config.Conf.Ws)
	s.Run()
}
