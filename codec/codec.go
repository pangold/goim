package codec

import (
	"gitlab.com/pangold/goim/codec/protobuf"
)

type Codec struct {
	c              *protobuf.Codec
	ackHandler     func(*MessageT)
	messageHandler func(*MessageT)
}

func NewCodec() *Codec {
	c := &Codec {
		c: protobuf.NewCodec(),
	}
	c.c.SetDecodeHandler(c.handleMessage)
	return c
}

// use to send.
func (c *Codec) SetEncodeHandler(h func([]byte)) {
	c.c.SetEncodeHandler(h)
}

func (c *Codec) SetDecodeHandler(h func(*MessageT)) {
	c.messageHandler = h
}

func (c *Codec) EnableResend(enable bool) {
	c.c.EnableResend(enable)
}

func (c *Codec) Decode(data []byte) int {
	return c.c.Decode(data)
}

func (c *Codec) Encode(msg *MessageT) error {
	m := &protobuf.Message {
		Id:       &msg.Id,
		UserId:   &msg.UserId,
		TargetId: &msg.TargetId,
		GroupId:  &msg.GroupId,
		Action:   &msg.Action,
		Ack:      &msg.Ack,
		Type:     &msg.Type,
		Body:     msg.Body,
	}
	return c.c.Encode(m)
}

func (c *Codec) handleMessage(msg *protobuf.Message) {
	m := &MessageT {
		Id:       msg.GetId(),
		UserId:   msg.GetUserId(),
		TargetId: msg.GetTargetId(),
		GroupId:  msg.GetGroupId(),
		Action:   msg.GetAction(),
		Ack:      msg.GetAck(),
		Type:     msg.GetType(),
		Body:     msg.Body,
	}
	c.messageHandler(m)
}