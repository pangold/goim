package protobuf

import (
	"github.com/golang/protobuf/proto"
	"log"
)

// wrap splitter and merger

type Codec struct {
	decoder        *Decoder
	encoder        *Encoder
	splitHandler   func([]byte)
	messageHandler func(*Message)
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

func (c *Codec) SetEncodeHandler(h func([]byte)) {
	c.splitHandler = h
}

func (c *Codec) SetDecodeHandler(h func(*Message)) {
	c.messageHandler = h
}

func (c *Codec) EnableResend(enable bool) {
	if enable {
		c.encoder.SetResendHandler(c.handleSegment)
	} else {
		c.encoder.SetResendHandler(nil)
	}
}

func (c *Codec) Encode(msg *Message) error {
	if err := c.encoder.Send(msg); err != nil {
		return err
	}
	return nil
}

func (c *Codec) Decode(data []byte) int {
	seg := &Segment{}
	//
	c.remaining = append(c.remaining, data...)
	if err := proto.Unmarshal(c.remaining, seg); err != nil {
		//log.Println(err.Error())
		return seg.XXX_Size()
	}
	// clear remaining
	c.remaining = c.remaining[seg.XXX_Size():]
	// received segment ack
	if seg.GetAck() == 1 {
		if err := c.encoder.SetAckSegment(seg); err != nil {
			log.Println(err.Error())
		}
		return seg.XXX_Size()
	}
	// received a non-ack segment(a real segment)
	if err := c.decoder.Push(seg); err != nil {
		log.Println(err.Error())
		return 0
	}
	// reply segment ack
	ack := &Segment{
		Id:    proto.Int64(seg.GetId()),
		Index: proto.Int32(0),
		Total: proto.Int32(1),
		Ack:   proto.Int32(1),
		Body:  nil,
	}
	c.handleSegment(ack)
	return seg.XXX_Size()
}

func (c *Codec) handleSegment(seg *Segment) {
	buf, err := proto.Marshal(seg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c.splitHandler(buf)
}

func (c *Codec) handleResend(seg *Segment) {
	c.handleSegment(seg)
}

func (c *Codec) handleMessage(msg *Message) {
	c.messageHandler(msg)
}
