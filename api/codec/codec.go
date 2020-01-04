package codec

import (
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/codec/protobuf"
	"log"
)

// wrap splitter and merger

type Codec struct {
	decoder        *Decoder
	encoder        *Encoder
	remaining      []byte
	segments       []*protobuf.Segment
}

func NewCodec() *Codec {
	c := &Codec{
		decoder: NewDecoder(),
		encoder: NewEncoder(),
	}
	return c
}

func (c *Codec) Encode(msg *protobuf.Message) []*protobuf.Segment {
	ss, err := c.encoder.Encode(msg)
	if err != nil {
		log.Printf(err.Error())
		return nil
	}
	return ss
}

func (c *Codec) Decode(data []byte) (res []*protobuf.Message) {
	c.remaining = append(c.remaining, data...)
	for {
		if len(c.remaining) == 0 {
			break
		}
		seg := &protobuf.Segment{}
		if err := proto.Unmarshal(c.remaining, seg); err != nil {
			break
		}
		c.remaining = c.remaining[seg.XXX_Size():]
		if len(c.segments) != int(seg.GetTotal()) {
			c.segments = make([]*protobuf.Segment, seg.GetTotal())
		}
		c.segments[seg.GetIndex()] = seg
		if len(c.segments) == int(seg.GetTotal()) {
			msg, _ := c.decoder.Decode(c.segments)
			res = append(res, msg)
		}
	}
	return res
}

