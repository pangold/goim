package codec

import (
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/protocol"
	"log"
)

// wrap splitter and merger

type Codec struct {
	decoder        *Decoder
	encoder        *Encoder
	remaining      []byte
	segments       []*protocol.Segment
}

func NewCodec() *Codec {
	c := &Codec{
		decoder: NewDecoder(),
		encoder: NewEncoder(),
	}
	return c
}

func (c *Codec) Encode(msg *protocol.Message) []*protocol.Segment {
	ss, err := c.encoder.Encode(msg)
	if err != nil {
		log.Printf(err.Error())
		return nil
	}
	return ss
}

func (c *Codec) Decode(data []byte) (res []*protocol.Message) {
	c.remaining = append(c.remaining, data...)
	for {
		if len(c.remaining) == 0 {
			break
		}
		seg := &protocol.Segment{}
		if err := proto.Unmarshal(c.remaining, seg); err != nil {
			break
		}
		c.remaining = c.remaining[seg.XXX_Size():]
		if msg := c.DecodeSegment(seg); msg != nil {
			res = append(res, msg)
		}
	}
	return res
}

func (c *Codec) DecodeSegment(seg *protocol.Segment) (res *protocol.Message) {
	if len(c.segments) != int(seg.GetTotal()) {
		c.segments = make([]*protocol.Segment, seg.GetTotal())
	}
	c.segments[seg.GetIndex()] = seg
	if len(c.segments) == int(seg.GetTotal()) {
		res, _ = c.decoder.Decode(c.segments)
	}
	return res
}

