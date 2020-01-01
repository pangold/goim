package im

import (
	"context"
	"fmt"
	"gitlab.com/pangold/goim/codec/protobuf"
	"gitlab.com/pangold/goim/config"
	"google.golang.org/grpc"
	"log"
	"time"
)

type Client struct {
	conn   *grpc.ClientConn
	context context.Context
	client  ApiClient
}

func NewClient(conf config.GrpcConfig) *Client {
	conn, err := grpc.Dial(conf.Address, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("failed to connect: %v", err))
	}
	client := NewApiClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	return &Client{
		conn:    conn,
		context: ctx,
		client:  client,
	}
}

func (c *Client) Send(uid, tid, gid string, action, t, ack int32, body []byte) bool {
	id := time.Now().UnixNano()
	msg := &protobuf.Message{
		Id:                   &id,
		UserId:               &uid,
		TargetId:             &tid,
		GroupId:              &gid,
		Action:               &action,
		Ack:                  &t,
		Type:                 &ack,
		Body:                 body,
	}
	res, err := c.client.Send(c.context, msg)
	if err != nil {
		log.Printf("failed to send: %v", err)
		return false
	}
	if res.GetResult() == ResultType_FAILURE {
		log.Printf("send result failure")
		return false
	}
	return true
}

func (c *Client) Broadcast(action, t, ack int32, body []byte) bool {
	id := time.Now().UnixNano()
	msg := &protobuf.Message{
		Id:                   &id,
		Action:               &action,
		Ack:                  &t,
		Type:                 &ack,
		Body:                 body,
	}
	res, err := c.client.Broadcast(c.context, msg)
	if err != nil {
		log.Printf("failed to broadcast: %v", err)
		return false
	}
	if res.GetResult() == ResultType_FAILURE {
		log.Printf("broadcast result failure")
		return false
	}
	return true
}

func (c *Client) GetConnections() []string {
	res, err := c.client.GetConnections(c.context, &EmptyRequest{})
	if err != nil {
		log.Printf("failed to get connections, error: %v", err)
		return nil
	}
	return res.GetUserIds()
}

func (c *Client) Online(uid string) bool {
	res, err := c.client.Online(c.context, &OnlineRequest{TargetId: &uid})
	if err != nil {
		log.Printf("failed to get online users(%s), error: %v", uid, err)
		return false
	}
	return res.GetResult()
}

func (c *Client) Kick(uid string) bool {
	res, err := c.client.Kick(c.context, &KickRequest{TargetId: &uid})
	if err != nil {
		log.Printf("failed to kick user(%s), error: %v", uid, err)
		return false
	}
	return res.GetResult() == ResultType_SUCCESS
}
