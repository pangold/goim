package v1

import (
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/codec/protobuf"
	"testing"
	"time"
)



func TestSplitter_1(t *testing.T) {
	//
	body := func(len int) []byte {
		buf := make([]byte, len)
		for i := 0; i < len; i++ {
			buf[i] = 'a'
		}
		return buf
	}

	msg := &protobuf.Message{
		Id:                   proto.Int64(time.Now().UnixNano()),
		UserId:               proto.String("10001"),
		TargetId:             proto.String("10002"),
		GroupId:              nil,
		Action:               proto.Int32(0),
		Ack:                  proto.Int32(0),
		Type:                 proto.Int32(0),
		Body:                 body(1200),
	}

	handle := func(conn interface{}, seg *protobuf.Segment) {
		// fmt.Printf("handle segment/resend callback, %d/%d, body: %s\n", seg.GetIndex(), seg.GetTotal(), seg.GetBody())
	}

	s := &Encoder{
		segmentHandler: handle,
		resendHandler:  &handle,
		acks:           make(map[int64]*acknowledge),
	}

	if err := s.Send(nil, msg); err != nil {
		t.Error(err)
	}

	if len(s.acks) != 1 {
		t.Error("unexpected size of ack size")
	}
	time.Sleep(time.Millisecond * 1000)
	//msg.
	msg.Id = proto.Int64(time.Now().UnixNano())
	if err := s.Send(nil, msg); err != nil {
		t.Error(err)
	}

	if len(s.acks) != 2 {
		t.Error("unexpected size of ack size")
	}

	// s.SetAckSegment()
}
