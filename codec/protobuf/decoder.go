package protobuf

import (
	"errors"
	"github.com/golang/protobuf/proto"
)

type Decoder struct {
	msgHandler     func(*Message)
	segments       map[int64][]*Segment
}

func NewDecoder() *Decoder {
	return &Decoder{
		msgHandler: nil,
		segments:   make(map[int64][]*Segment),
	}
}

func (d *Decoder) SetMessageHandler(h func(*Message)) {
	d.msgHandler = h
}

//
func (d *Decoder) Push(seg *Segment) error {
	// optimize for single segment
	if seg.GetTotal() == 1 {
		return d.single(seg.GetBody())
	}
	// for multi segments
	d.segments[seg.GetId()] = append(d.segments[seg.GetId()], seg)
	ss := d.segments[seg.GetId()]
	if len(ss) == int(ss[0].GetTotal()) {
		return d.multi(ss)
	}
	return nil
}

// []*message.Segment is not in order,
// But the size of every body of segment is the same
// FIXME: what if []*message.Segment's size is huge...
//        such as video file
func (d *Decoder) multi(sl []*Segment) error {
	buf := make([]byte, MAX_SEGMENT_SIZE* (len(sl) - 1))
	for _, seg := range sl {
		// the last segment of a message
		if int(seg.GetIndex()) == len(sl) - 1 {
			buf = append(buf, seg.GetBody()...)
			continue
		}
		if len(seg.GetBody()) != MAX_SEGMENT_SIZE {
			return errors.New("combine error: unexpected segment size")
		}
		begin := int(seg.GetIndex()) * MAX_SEGMENT_SIZE
		end := int(seg.GetIndex() + 1) * MAX_SEGMENT_SIZE
		copy(buf[begin:end], seg.GetBody())
	}
	return d.single(buf)
}

// A message with single one segment
func (m *Decoder) single(buf []byte) error {
	msg := &Message{}
	if err := proto.Unmarshal(buf, msg); err != nil {
		return errors.New("unmarshal message error: " + err.Error())
	}
	m.msgHandler(msg)
	return nil
}
