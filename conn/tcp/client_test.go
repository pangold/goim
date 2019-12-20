package tcp

import (
	"gitlab.com/pangold/goim/config"
	"testing"
)

func TestServer_Run(t *testing.T) {
	client := NewTcpClient(config.Conf.Tcp)
	client.Run()
}
