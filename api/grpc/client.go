package rpc

import (
	"context"
	"fmt"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/protocol"
	"google.golang.org/grpc"
	"log"
	"time"
)

type Client struct {
	conn         *grpc.ClientConn
	context      context.Context
	ImApi        protocol.ImApiServiceClient
	ImDispatcher protocol.ImDispatchServiceClient
}

func NewClient(conf config.GrpcConfig) *Client {
	conn, err := grpc.Dial(conf.Address, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("failed to connect: %v", err))
	}
	ctx := context.Background()
	// ctx, _ := context.WithTimeout(context.Background(), time.Second)
	return &Client{
		conn:         conn,
		context:      ctx,
		ImApi:        protocol.NewImApiServiceClient(conn),
		ImDispatcher: protocol.NewImDispatchServiceClient(conn),
	}
}

func (c *Client) Send(uid, tid, gid string, action, t, ack int32, body []byte) bool {
	id := time.Now().UnixNano()
	msg := &protocol.Message{
		Id:                   id,
		UserId:               uid,
		TargetId:             tid,
		GroupId:              gid,
		Action:               action,
		Ack:                  t,
		Type:                 ack,
		Body:                 body,
	}
	cli, _ := c.ImApi.Send(context.Background())
	if err := cli.Send(msg); err != nil {
		log.Printf("send message error: %v", err)
		return false
	}
	return true
}

func (c *Client) Broadcast(action, t, ack int32, body []byte) bool {
	id := time.Now().UnixNano()
	msg := &protocol.Message{
		Id:                   id,
		Action:               action,
		Ack:                  t,
		Type:                 ack,
		Body:                 body,
	}
	cli, _ := c.ImApi.Broadcast(context.Background())
	if err := cli.Send(msg); err != nil {
		log.Printf("broadcast message error: %v", err)
		return false
	}
	return true
}

func (c *Client) GetConnections() []string {
	res, err := c.ImApi.GetConnections(c.context, &protocol.Empty{})
	if err != nil {
		log.Printf("failed to get connections, error: %v", err)
		return nil
	}
	return res.GetUserIds()
}

func (c *Client) Online(uid string) bool {
	res, err := c.ImApi.Online(c.context, &protocol.User{UserId: uid})
	if err != nil {
		log.Printf("failed to get online users(%s), error: %v", uid, err)
		return false
	}
	return res.GetSuccess()
}

func (c *Client) Kick(uid string) bool {
	res, err := c.ImApi.Kick(c.context, &protocol.User{UserId: uid})
	if err != nil {
		log.Printf("failed to kick user(%s), error: %v", uid, err)
		return false
	}
	return res.GetSuccess()
}

// handle dispatch messages / sessions
func (c *Client) GetDispatchedMessages() {
	cli, _ := c.ImDispatcher.Dispatch(context.Background(), &protocol.Empty{})
	for {
		msg, err := cli.Recv()
		if err != nil {
			log.Printf("get dispatched message error: %v", err)
			break
		}
		// TODO: your code
		log.Println(msg)
	}
}

func (c *Client) GetSessionIn() {
	cli, _ := c.ImDispatcher.SessionIn(context.Background(), &protocol.Empty{})
	for {
		session, err := cli.Recv()
		if err != nil {
			log.Printf("get session in error: %v", err)
			break
		}
		// TODO: your code
		log.Println(session)
	}
}

func (c *Client) GetSessionOut() {
	cli, _ := c.ImDispatcher.SessionOut(context.Background(), &protocol.Empty{})
	for {
		session, err := cli.Recv()
		if err != nil {
			log.Printf("get session in error: %v", err)
			break
		}
		// TODO: your code
		log.Println(session)
	}
}
