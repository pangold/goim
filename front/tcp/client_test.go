package tcp

import (
	"gitlab.com/pangold/goim/config"
	"strconv"
	"testing"
)

func TestTcpClient_Run(t *testing.T) {
	client := NewTcpClient("f3ba0e819f55cb7", config.Conf.Front.Tcp)
	for i := 0; i < 100000; i++ {
		client.SendMessage(strconv.Itoa(i) +
			"012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789")
	}
	client.Stop()
}
