package impl

import (
	"github.com/golang/protobuf/proto"
	message "gitlab.com/pangold/goim/msg/protobuf"
	"log"
)

type Combiner struct {
	err            error
	messageHandler func(*message.Message)
	ackHandler func(*message.Message)
	segments       map[int64][]*message.Segment
}

func (c *Combiner) Error() error {
	return c.err
}

func (c *Combiner) SetMessageHandler(h func(*message.Message)) {
	c.messageHandler = h
}

func (c *Combiner) SetAckHandler(h func(*message.Message)) {
	c.ackHandler = h
}

func (c *Combiner) Push(seg *message.Segment) {
	// optimize for single segment
	if seg.GetTotal() == 1 {
		c.single(seg.GetBody())
		return
	}
	// for multi segments
	c.append(seg)
}

// A message with single one segment
func (c *Combiner) single(buf []byte) {
	msg := &message.Message{}
	if err := proto.Unmarshal(buf, msg); err != nil {
		log.Printf("unmarshal message error: %v", err)
		return
	}
	if msg.GetAck() == message.Message_NONE {
		c.messageHandler(msg)
		return
	}
	c.ackHandler(msg)
}

func (c *Combiner) append(seg *message.Segment) {
	c.segments[seg.GetId()] = append(c.segments[seg.GetId()], seg)
	ss := c.segments[seg.GetId()]
	total := ss[0].GetTotal()
	if len(ss) == int(total) {
		c.combine(ss)
	}
}

// []*message.Segment is not in order,
// But the size of every body of segment is the same
// FIXME: what if []*message.Segment's size is huge...
//        such as video file
func (c *Combiner) combine(sl []*message.Segment) {
	// TODO: cap of buf
	buf := make([]byte, MaxSegmentSize * (len(sl) - 1))
	for _, seg := range sl {
		// the last segment of a message
		if int(seg.GetIndex()) == len(sl) - 1 {
			buf = append(buf, seg.GetBody()...)
			continue
		}
		if len(seg.GetBody()) != MaxSegmentSize {
			log.Fatal("combine error: unexpected segment size")
			return
		}
		begin := int(seg.GetIndex()) * MaxSegmentSize
		end := int(seg.GetIndex()+1) * MaxSegmentSize
		copy(buf[begin:end], seg.GetBody())
	}
	c.single(buf)
}
