package impl

import (
	"github.com/golang/protobuf/proto"
	message "gitlab.com/pangold/goim/msg/protobuf"
	"testing"
	"time"
)



func TestSplitter_Send(t *testing.T) {
	//
	body := func(len int) []byte {
		buf := make([]byte, len)
		for i := 0; i < len; i++ {
			buf[i] = 'a'
		}
		return buf
	}

	msg := &message.Message{
		UserId:               proto.String("10001"),
		TargetId:             proto.String("10002"),
		GroupId:              nil,
		Type:                 (*message.Message_MessageType)(proto.Int32(int32(message.Message_TEXT))),
		Ack:                  (*message.Message_AckType)(proto.Int32(int32(message.Message_NONE))),
		Body:                 body(1200),
		Time:                 proto.Int64(time.Now().Unix() >> 1),
	}

	handle := func(seg *message.Segment) {
		// fmt.Printf("handle segment/resend callback, %d/%d, body: %s\n", seg.GetIndex(), seg.GetTotal(), seg.GetBody())
	}

	s := &Splitter {
		segmentHandler: handle,
		resendHandler:  handle,
		acks:           make(map[int64]*acknowledge),
	}

	if err := s.Send(msg); err != nil {
		t.Error(err)
	}

	if len(s.acks) != 1 {
		t.Error("unexpected size of ack size")
	}
	time.Sleep(time.Millisecond * 1000)
	//msg.
	if err := s.Send(msg); err != nil {
		t.Error(err)
	}

	if len(s.acks) != 2 {
		t.Error("unexpected size of ack size")
	}

	// s.SetAckSegment()
}
