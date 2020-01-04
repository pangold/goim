package codec

import (
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/front/interfaces"
	"gitlab.com/pangold/goim/protocol"
	"log"
)

// wrap splitter and merger

type Codec struct {
	decoder        *Decoder
	encoder        *Encoder
	decodeHandler  func(interfaces.Conn, *proto.Message)
	remaining      []byte
}

func NewCodec() *Codec {
	c := &Codec{
		decoder: NewDecoder(),
		encoder: NewEncoder(),
	}
	c.decoder.MsgHandler = c.handleMessage
	return c
}

func (c *Codec) SetDecodeHandler(h func(interfaces.Conn, *proto.Message)) {
	c.decodeHandler = h
}

func (c *Codec) EnableResend(enable bool) {
	c.encoder.ResendEnabled = enable
}

func (c *Codec) Send(conn interfaces.Conn, msg *proto.Message) error {
	if err := c.encoder.Send(conn, msg); err != nil {
		return err
	}
	return nil
}

func (c *Codec) Decode(conn interfaces.Conn, data []byte) {
	c.remaining = append(c.remaining, data...)
	for {
		if len(c.remaining) == 0 {
			break
		}
		seg := &proto.Segment{}
		if err := proto.Unmarshal(c.remaining, seg); err != nil {
			//log.Println(err.Error())
			break
		}
		// clear remaining
		c.remaining = c.remaining[seg.XXX_Size():]
		// received an ack segment
		if seg.GetAck() == 1 {
			c.encoder.SetAckSegment(seg)
			continue
		}
		//fmt.Printf("segment received: %d\n", len(c.remaining))
		// received a non-ack segment(a real segment)
		if err := c.decoder.Decode(conn, seg); err != nil {
			log.Println(err.Error())
			continue
		}
		// reply segment ack
		if c.encoder.ResendEnabled {
			c.ack(conn, seg.GetId())
		}
	}
}

func (c *Codec) ack(conn interfaces.Conn, id int64) {
	ack := &proto.Segment{
		Id:    proto.Int64(id),
		Index: proto.Int32(0),
		Total: proto.Int32(1),
		Ack:   proto.Int32(1),
		Body:  nil,
	}
	c.encoder.send(conn, ack)
}

func (c *Codec) handleMessage(conn interfaces.Conn, msg *proto.Message) {
	c.decodeHandler(conn, msg)
}
