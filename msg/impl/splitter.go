package impl

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	message "gitlab.com/pangold/goim/msg/protobuf"
	"log"
	"time"
)

const (
	MaxSegmentSize = 1024
	WaitResendTime = time.Second * 10
)

type Splitter struct {
	segmentHandler func(*message.Segment)
	resendHandler  func(*message.Segment)
	acks           map[int64]*acknowledge //
}

type acknowledge struct {
	id       int64
	segments []*message.Segment
	timer    *time.Timer
	retry    int
}

func NewSplitter() *Splitter {
	return &Splitter{
		segmentHandler: nil,
		resendHandler:  nil,
		acks:           make(map[int64]*acknowledge),
	}
}

func (s *Splitter) SetSegmentHandler(h func(*message.Segment)) {
	s.segmentHandler = h
}

func (s *Splitter) SetResendHandler(h func(*message.Segment)) {
	s.resendHandler = h
}

// without ack, segment will be held in acknowledge list
func (s *Splitter) SetAckSegment(seg *message.Segment) error {
	if ack, ok := s.acks[seg.GetId()]; ok {
		// being ack, clear
		// optimize: release by gc, or release immediately yourself
		ack.segments[seg.GetIndex()] = nil
		return nil
	}
	return errors.New(fmt.Sprintf("invalid ack id: %d, (%d/%d)", seg.GetId(), seg.GetIndex(), seg.GetTotal()))
}

func (s *Splitter) Send(msg *message.Message) error {
	ack := &acknowledge{}
	// generate msg id: simply using Unix Time as id
	// very moment, only one message for one client
	// FIXME: second level....
	ack.id = time.Now().Unix()
	s.acks[ack.id] = ack
	// split will trigger segment callback
	if err := s.split(ack, msg); err != nil {
		return err
	}
	// ack.segments = segments
	// check timeout, trigger resend callback
	// optimization: a common
	ack.timer = time.NewTimer(WaitResendTime)
	// FIXME: unstable retry times
	ack.retry = 3 * len(ack.segments)
	go s.timeout(ack)
	return nil
}

func (s *Splitter) timeout(ack *acknowledge) {
	select {
	case <-ack.timer.C:
		s.resend(ack)
	}
}

func (s *Splitter) resend(ack *acknowledge) {
	var count = 0
	for _, seg := range ack.segments {
		if seg == nil {
			continue
		}
		s.resendHandler(seg)
		if ack.retry == 0 {
			log.Printf("segment retry failed: message %d, body: %s", ack.id, ack.segments)
		} else {
			go s.timeout(ack)
		}
		ack.retry--
		count++
	}
	// nothing left, clear it
	if count == 0 {
		delete(s.acks, ack.id)
	}
}

func (s *Splitter) split(ack *acknowledge, msg *message.Message) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return errors.New("split error: " + err.Error())
	}
	count := int32(len(buf) / (MaxSegmentSize + 1)) + 1
	ack.segments = make([]*message.Segment, count)
	for index := int32(0); index < count; index++ {
		seg := s.single(ack.id, index, count, buf)
		ack.segments[index] = seg
		s.segmentHandler(seg)
	}
	return nil
}

func (s *Splitter) single(id int64, index, count int32, buf []byte) *message.Segment {
	end := int((index + 1) * MaxSegmentSize)
	if index == count - 1 {
		end = len(buf)
	}
	return &message.Segment {
		Id:      proto.Int64(id),
		Index:   proto.Int32(index),
		Total:   proto.Int32(count),
		Ack:     proto.Int32(0),
		Body:    buf[index * MaxSegmentSize : end],
	}
}