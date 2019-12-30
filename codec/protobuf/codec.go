package protobuf

import (
	"github.com/golang/protobuf/proto"
	"log"
)

// wrap splitter and merger

type Codec struct {
	decoder        *Decoder
	encoder        *Encoder
	encodeHandler  func(interface{}, []byte)
	decodeHandler  func(interface{}, *Message)
	remaining      []byte
}

func NewCodec() *Codec {
	c := &Codec{
		decoder: NewDecoder(),
		encoder: NewEncoder(),
	}
	c.encoder.SetSegmentHandler(c.handleSegment)
	c.decoder.SetMessageHandler(c.handleMessage)
	return c
}

func (c *Codec) SetEncodeHandler(h func(interface{}, []byte)) {
	c.encodeHandler = h
}

func (c *Codec) SetDecodeHandler(h func(interface{}, *Message)) {
	c.decodeHandler = h
}

func (c *Codec) EnableResend(enable bool) {
	if enable {
		c.encoder.SetResendHandler(c.handleSegment)
	} else {
		c.encoder.SetResendHandler(nil)
	}
}

func (c *Codec) Encode(conn interface{}, msg *Message) error {
	if err := c.encoder.Send(conn, msg); err != nil {
		return err
	}
	return nil
}

func (c *Codec) Decode(conn interface{}, data []byte) {
	seg := &Segment{}
	//
	c.remaining = append(c.remaining, data...)
	if err := proto.Unmarshal(c.remaining, seg); err != nil {
		//log.Println(err.Error())
		return
	}
	// clear remaining
	c.remaining = c.remaining[seg.XXX_Size():]
	// received an ack segment
	if seg.GetAck() == 1 {
		c.encoder.SetAckSegment(seg)
		return
	}
	// received a non-ack segment(a real segment)
	if err := c.decoder.Push(conn, seg); err != nil {
		log.Println(err.Error())
		return
	}
	// reply segment ack
	if c.encoder.ResendEnabled() {
		c.ack(conn, seg.GetId())
	}
}

func (c *Codec) ack(conn interface{}, id int64) {
	ack := &Segment{
		Id:    proto.Int64(id),
		Index: proto.Int32(0),
		Total: proto.Int32(1),
		Ack:   proto.Int32(1),
		Body:  nil,
	}
	c.handleSegment(conn, ack)
}

func (c *Codec) handleSegment(conn interface{}, seg *Segment) {
	buf, err := proto.Marshal(seg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c.encodeHandler(conn, buf)
}

func (c *Codec) handleResend(conn interface{}, seg *Segment) {
	c.handleSegment(conn, seg)
}

func (c *Codec) handleMessage(conn interface{}, msg *Message) {
	c.decodeHandler(conn, msg)
}
