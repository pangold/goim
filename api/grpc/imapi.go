package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/front"
	"gitlab.com/pangold/goim/protocol"
	"log"
)

type ImApiService struct {
	front    *front.Server
	sessions *session.Sessions
}

func NewImApiService(f *front.Server, ss *session.Sessions) *ImApiService {
	return &ImApiService{
		front:    f,
		sessions: ss,
	}
}

func (c *ImApiService) GetConnections(ctx context.Context, req *Empty) (res *Users, err error) {
	res = &Users{
		UserIds: c.sessions.GetUserIds(),
	}
	return res, err
}

func (c *ImApiService) send(msg *protocol.Message) error {
	if msg.GetTargetId() != "" {
		return errors.New("uid could not be null")
	}
	token := c.sessions.GetTokenByUserId(msg.GetTargetId())
	if token == "" {
		return fmt.Errorf("uid(%s) is not online", msg.GetTargetId())
	}
	if err := c.front.Send(token, msg); err != nil {
		return err
	}
	return nil
}

func (c *ImApiService) Send(srv ImApiService_SendServer) error {
	for {
		msg, err := srv.Recv()
		if err != nil {
			log.Printf("grpc send steam break error: %v", err)
			break
		}
		if err := c.send(msg); err != nil {
			log.Printf("send error: %v", err)
		}
	}
	return nil
}

func (c *ImApiService) Broadcast(srv ImApiService_BroadcastServer) error {
	for {
		msg, err := srv.Recv()
		if err != nil {
			log.Printf("grpc broadcast stream break error: %v", err)
			break
		}
		c.front.Broadcast(msg)
	}
	return nil
}

func (c *ImApiService) Online(ctx context.Context, req *User) (*Result, error) {
	if req.GetUserId() == "" {
		return nil, errors.New("uid could not be null")
	}
	token := c.sessions.GetTokenByUserId(req.GetUserId())
	return &Result{Success: proto.Bool(token != "")}, nil
}

func (c *ImApiService) Kick(ctx context.Context, req *User) (*Result, error) {
	if req.GetUserId() != "" {
		return nil, errors.New("uid could not be null")
	}
	token := c.sessions.GetTokenByUserId(req.GetUserId())
	if token == "" {
		return nil, fmt.Errorf("uid(%s) is not online", req.GetUserId())
	}
	c.front.Remove(token)
	return &Result{Success: proto.Bool(true)}, nil
}