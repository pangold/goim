package impl

import (
	"errors"
	"github.com/golang/protobuf/proto"
	message "gitlab.com/pangold/goim/msg/protobuf"
)

type Merger struct {
	msgHandler     func(*message.Message)
	ackHandler     func(*message.Message)
	segments       map[int64][]*message.Segment
}

func NewMerger() *Merger {
	return &Merger{
		msgHandler: nil,
		ackHandler: nil,
		segments:   make(map[int64][]*message.Segment),
	}
}

func (m *Merger) SetMessageHandler(h func(*message.Message)) {
	m.msgHandler = h
}

func (m *Merger) SetAckHandler(h func(*message.Message)) {
	m.ackHandler = h
}

//
// func (m *Merger) Push(data []byte) {
func (m *Merger) Push(seg *message.Segment) error {
	// optimize for single segment
	if seg.GetTotal() == 1 {
		return m.single(seg.GetBody())
	}
	// for multi segments
	return m.append(seg)
}

// A message with single one segment
func (m *Merger) single(buf []byte) error {
	msg := &message.Message{}
	if err := proto.Unmarshal(buf, msg); err != nil {
		return errors.New("unmarshal message error: " + err.Error())
	}
	if msg.GetAck() == message.Message_NONE {
		m.msgHandler(msg)
	} else {
		m.ackHandler(msg)
	}
	return nil
}

func (m *Merger) append(seg *message.Segment) error {
	m.segments[seg.GetId()] = append(m.segments[seg.GetId()], seg)
	ss := m.segments[seg.GetId()]
	if len(ss) == int(ss[0].GetTotal()) {
		return m.combine(ss)
	}
	return nil
}

// []*message.Segment is not in order,
// But the size of every body of segment is the same
// FIXME: what if []*message.Segment's size is huge...
//        such as video file
func (m *Merger) combine(sl []*message.Segment) error {
	buf := make([]byte, MaxSegmentSize * (len(sl) - 1))
	for _, seg := range sl {
		// the last segment of a message
		if int(seg.GetIndex()) == len(sl) - 1 {
			buf = append(buf, seg.GetBody()...)
			continue
		}
		if len(seg.GetBody()) != MaxSegmentSize {
			return errors.New("combine error: unexpected segment size")
		}
		begin := int(seg.GetIndex()) * MaxSegmentSize
		end := int(seg.GetIndex() + 1) * MaxSegmentSize
		copy(buf[begin:end], seg.GetBody())
	}
	return m.single(buf)
}
