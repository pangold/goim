package im

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/api/front"
	"gitlab.com/pangold/goim/codec/protobuf"
)

type Controller struct {
	front *front.Server
}

func NewController(front *front.Server) *Controller {
	return &Controller{
		front: front,
	}
}

func (c *Controller) GetConnections(ctx context.Context, req *EmptyRequest) (res *ConnectionList, err error) {
	res = &ConnectionList{
		UserIds: c.front.GetOnlineUserIds(),
	}
	return res, err
}

func (c *Controller) Send(ctx context.Context, req *protobuf.Message) (*Result, error) {
	rtype := ResultType_FAILURE
	if req.GetTargetId() != "" {
		return &Result{Result: &rtype}, errors.New("uid could not be null")
	}
	token := c.front.GetOnlineTokenByUserId(req.GetTargetId())
	if token == "" {
		return &Result{Result: &rtype}, fmt.Errorf("uid(%s) is not online", req.GetTargetId())
	}
	c.front.SendEx(token, req)
	rtype = ResultType_SUCCESS
	return &Result{Result: &rtype}, nil
}

func (c *Controller) Broadcast(ctx context.Context, req *protobuf.Message) (*Result, error) {
	c.front.BroadcastEx(req)
	rtype := ResultType_SUCCESS
	return &Result{Result: &rtype}, nil
}

func (c *Controller) Online(ctx context.Context, req *OnlineRequest) (*OnlineResult, error) {
	if req.GetTargetId() == "" {
		return &OnlineResult{Result: proto.Bool(false)}, errors.New("uid could not be null")
	}
	token := c.front.GetOnlineTokenByUserId(req.GetTargetId())
	return &OnlineResult{Result: proto.Bool(token != "")}, nil
}

func (c *Controller) Kick(ctx context.Context, req *KickRequest) (*Result, error) {
	rtype := ResultType_FAILURE
	if req.GetTargetId() != "" {
		return &Result{Result: &rtype}, errors.New("uid could not be null")
	}
	token := c.front.GetOnlineTokenByUserId(req.GetTargetId())
	if token == "" {
		return &Result{Result: &rtype}, fmt.Errorf("uid(%s) is not online", req.GetTargetId())
	}
	c.front.Remove(token)
	rtype = ResultType_SUCCESS
	return &Result{Result: &rtype}, nil
}
