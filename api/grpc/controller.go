package im

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/codec/protobuf"
	"gitlab.com/pangold/goim/front"
)

type Controller struct {
	front    *front.Server
	sessions *session.Sessions
}

func NewController(f *front.Server, ss *session.Sessions) *Controller {
	return &Controller{
		front:    f,
		sessions: ss,
	}
}

func (c *Controller) GetConnections(ctx context.Context, req *EmptyRequest) (res *ConnectionList, err error) {
	res = &ConnectionList{
		UserIds: c.sessions.GetUserIds(),
	}
	return res, err
}

func (c *Controller) Send(ctx context.Context, req *protobuf.Message) (*Result, error) {
	rtype := ResultType_FAILURE
	if req.GetTargetId() != "" {
		return &Result{Result: &rtype}, errors.New("uid could not be null")
	}
	token := c.sessions.GetTokenByUserId(req.GetTargetId())
	if token == "" {
		return &Result{Result: &rtype}, fmt.Errorf("uid(%s) is not online", req.GetTargetId())
	}
	if err := c.front.Send(token, req); err != nil {
		return nil, err
	}
	rtype = ResultType_SUCCESS
	return &Result{Result: &rtype}, nil
}

func (c *Controller) Broadcast(ctx context.Context, req *protobuf.Message) (*Result, error) {
	c.front.Broadcast(req)
	rtype := ResultType_SUCCESS
	return &Result{Result: &rtype}, nil
}

func (c *Controller) Online(ctx context.Context, req *OnlineRequest) (*OnlineResult, error) {
	if req.GetTargetId() == "" {
		return &OnlineResult{Result: proto.Bool(false)}, errors.New("uid could not be null")
	}
	token := c.sessions.GetTokenByUserId(req.GetTargetId())
	return &OnlineResult{Result: proto.Bool(token != "")}, nil
}

func (c *Controller) Kick(ctx context.Context, req *KickRequest) (*Result, error) {
	rtype := ResultType_FAILURE
	if req.GetTargetId() != "" {
		return &Result{Result: &rtype}, errors.New("uid could not be null")
	}
	token := c.sessions.GetTokenByUserId(req.GetTargetId())
	if token == "" {
		return &Result{Result: &rtype}, fmt.Errorf("uid(%s) is not online", req.GetTargetId())
	}
	c.front.Remove(token)
	rtype = ResultType_SUCCESS
	return &Result{Result: &rtype}, nil
}
