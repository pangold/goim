package api

import (
	"fmt"
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/front"
	"testing"
	"time"
)

var (
	server *Server
	client *Client
	frontServer *front.Server
	sessions *session.Sessions
)

func init() {
	frontServer = front.NewServer(config.Conf)
	sessions = session.NewSessions()
	server = NewServer(frontServer, sessions, config.Conf.Grpc)
	client = NewClient(config.Conf.Grpc)
	go frontServer.Run()
	go server.Run()
}

func TestServer_Run(t *testing.T) {
	time.Sleep(time.Second * 5)

	conns := client.GetConnections()
	if len(conns) != 0 {
		t.Error("unexpected connection count")
	}

	if !client.Send("10001", "10002", "", 1, 1, 0, []byte("test text")) {
		fmt.Println("failed to send")
	}

	if !client.Broadcast(1, 1, 0, []byte("test")) {
		fmt.Println("failed to broadcast")
	}

	res := client.Online("10001")
	fmt.Println(res)

	client.Kick("10001")
}
