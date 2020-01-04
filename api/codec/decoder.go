package codec

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/protocol"
)

type Decoder struct {

}

func NewDecoder() *Decoder {
	return &Decoder{}
}

// because of reset mechanism, seg may be exist
func (d *Decoder) Decode(segs []*protocol.Segment) (*protocol.Message, error) {
	// optimize for single segment
	if segs[0].GetTotal() == 1 {
		return d.single(segs[0].GetBody())
	}
	// for multi segments
	if len(segs) == int(segs[0].GetTotal()) {
		return d.multi(segs)
	}
	return nil, nil
}

// The size of body of segments are the same, except the last segment
func (d *Decoder) multi(segs []*protocol.Segment) (*protocol.Message, error) {
	buf := make([]byte, MAX_SEGMENT_SIZE* (len(segs) - 1))
	for i := 0; i < len(segs) - 1; i++ {
		if len(segs[i].GetBody()) > MAX_SEGMENT_SIZE {
			return nil, errors.New("unexpected segment size")
		}
		begin := int(segs[i].GetIndex()) * MAX_SEGMENT_SIZE
		end := int(segs[i].GetIndex() + 1) * MAX_SEGMENT_SIZE
		copy(buf[begin:end], segs[i].GetBody())
	}
	// last one is rarely full.
	buf = append(buf, segs[len(segs) - 1].GetBody()...)
	return d.single(buf)
}

// A message with single one segment
func (d *Decoder) single(buf []byte) (*protocol.Message, error) {
	msg := &protocol.Message{}
	if err := proto.Unmarshal(buf, msg); err != nil {
		return nil, errors.New("unmarshal error: " + err.Error())
	}
	return msg, nil
}

