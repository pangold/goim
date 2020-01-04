package codec

import (
	"gitlab.com/pangold/goim/codec/protobuf"
	"gitlab.com/pangold/goim/front/codec"
)

type Codec struct {
	c              *codec.Codec
	decodeHandler func(interface{}, *MessageT)
}

func NewCodec() *Codec {
	c := &Codec {
		c: codec.NewCodec(),
	}
	c.c.SetDecodeHandler(c.handleDecodec)
	return c
}

func (c *Codec) EnableResend(enable bool) {
	c.c.EnableResend(enable)
}

// use to send.
func (c *Codec) SetEncodeHandler(handler func(interface{}, []byte)) {
	c.c.SetEncodeHandler(handler)
}

func (c *Codec) SetDecodeHandler(handler func(interface{}, *MessageT)) {
	c.decodeHandler = handler
}

func (c *Codec) Decode(conn interface{}, data []byte) {
	c.c.Decode(conn, data)
}

func (c *Codec) Encode(conn interface{}, msg *MessageT) error {
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
	return c.c.Encode(conn, m)
}

func (c *Codec) handleDecodec(conn interface{}, msg *protobuf.Message) {
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
	c.decodeHandler(conn, m)
}