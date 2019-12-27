package impl

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	message "gitlab.com/pangold/goim/msg/protobuf"
	"time"
)

const (
	MaxSegmentSize = 1024
	WaitResendTime = time.Second * 10
)

type Splitor struct {
	err            error
	timer          time.Timer
	segmentHandler func(*message.Segment)
	resendHandler  func(*message.Segment)
	acks           map[int64]*acknowledge //
}

type acknowledge struct {
	id int64
	segments []*message.Segment
	timer *time.Timer
}

func (s *Splitor) Error() error {
	return s.err
}

func (s *Splitor) SetSegmentHandler(handler func(*message.Segment)) {
	s.segmentHandler = handler
}

func (s *Splitor) SetResendHandler(handler func(*message.Segment)) {
	s.resendHandler = handler
}

func (s *Splitor) SetAckSegment(seg *message.Segment) {
	if ack, ok := s.acks[seg.GetId()]; ok {
		// being ack, clear
		ack.segments[seg.GetIndex()] = nil // release by gc, or release immediately yourself
	} else {
		s.err = errors.New(fmt.Sprintf("invalid ack id: %d, (%d/%d)", seg.GetId(), seg.GetIndex(), seg.GetTotal()))
	}
}

func (s *Splitor) Send(msg *message.Message) {
	ack := &acknowledge{}
	// generate msg id: simply using Unix Time as id
	// very moment, only one message for one client
	ack.id = time.Now().Unix()
	// split will trigger segment callback
	ack.segments = s.split(ack.id, msg)
	// check timeout, trigger resend callback
	// optimization: a common
	ack.timer = time.NewTimer(WaitResendTime)
	go s.timeout(ack)
}

func (s *Splitor) timeout(ack *acknowledge) {
	select {
	case <-ack.timer.C:
		s.resend(ack)
	}
}

func (s *Splitor) resend(ack *acknowledge) {
	for _, seg := range ack.segments {
		if seg != nil {
			s.resendHandler(seg)
		}
	}
}

func (s *Splitor) split(id int64, msg *message.Message) []*message.Segment {
	buf, err := proto.Marshal(msg)
	if err != nil {
		s.err = errors.New("split error: " + err.Error())
		return nil
	}
	count := len(buf) / (MaxSegmentSize + 1)
	result := make([]*message.Segment, count)
	for index := 0; index < count; index++ {
		var body []byte
		if index == count {
			body = buf[index * MaxSegmentSize : ]
		} else {
			body = buf[index * MaxSegmentSize : (index + 1) * MaxSegmentSize]
		}
		result[index] = &message.Segment {
			Id:      proto.Int64(id),
			Index:   proto.Int32(int32(index)),
			Total:   proto.Int32(int32(count)),
			Ack:     nil,
			Body:    body,
		}
		s.segmentHandler(result[index])
	}
	return result
}