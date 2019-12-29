package protobuf

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"time"
)

const (
	MAX_SEGMENT_SIZE = 1024
	WAIT_RESEND_TIME = time.Second * 3
)

type Encoder struct {
	segmentHandler  func(*Segment)
	resendHandler  *func(*Segment)
	acks            map[int64]*acknowledge //
}

type acknowledge struct {
	id       int64
	segments []*Segment
	timer    *time.Timer
	retry    int
}

func NewEncoder() *Encoder {
	return &Encoder{
		segmentHandler: nil,
		resendHandler:  nil,
		acks:           make(map[int64]*acknowledge),
	}
}

func (e *Encoder) SetSegmentHandler(h func(*Segment)) {
	e.segmentHandler = h
}

func (e *Encoder) SetResendHandler(h func(*Segment)) {
	e.resendHandler = &h
}

// without ack, segment will be held in acknowledge list
func (e *Encoder) SetAckSegment(seg *Segment) error {
	if ack, ok := e.acks[seg.GetId()]; ok {
		// being ack, clear
		// optimize: release by gc, or release immediately yourself
		ack.segments[seg.GetIndex()] = nil
		return nil
	}
	return errors.New(fmt.Sprintf("invalid ack id: %d, (%d/%d)", seg.GetId(), seg.GetIndex(), seg.GetTotal()))
}

func (e *Encoder) Send(msg *Message) error {
	ack := &acknowledge{}
	//
	ack.id = msg.GetId()
	e.acks[ack.id] = ack
	// split will trigger segment callback
	if err := e.split(ack, msg); err != nil {
		return err
	}
	// ack.segments = segments
	// check timeout, trigger resend callback
	// resend works unless s.resendHandler is being specific.
	if e.resendHandler != nil {
		ack.timer = time.NewTimer(WAIT_RESEND_TIME)
		// FIXME: unstable retry times
		ack.retry = 3 * len(ack.segments)
		go e.timeout(ack)
	}
	return nil
}

func (e *Encoder) timeout(ack *acknowledge) {
	select {
	case <-ack.timer.C:
		e.resend(ack)
	}
}

func (e *Encoder) resend(ack *acknowledge) {
	var count = 0
	for _, seg := range ack.segments {
		if seg == nil {
			continue
		}
		if e.resendHandler != nil {
			(*e.resendHandler)(seg)
		}
		if ack.retry == 0 {
			log.Printf("segment retry failed: message %d, body: %s", ack.id, ack.segments)
		} else {
			go e.timeout(ack)
		}
		ack.retry--
		count++
	}
	// nothing left, clear it
	if count == 0 {
		delete(e.acks, ack.id)
	}
}

func (e *Encoder) split(ack *acknowledge, msg *Message) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return errors.New("split error: " + err.Error())
	}
	count := int32(len(buf) / (MAX_SEGMENT_SIZE + 1)) + 1
	ack.segments = make([]*Segment, count)
	for index := int32(0); index < count; index++ {
		seg := e.single(ack.id, index, count, buf)
		ack.segments[index] = seg
		e.segmentHandler(seg)
	}
	return nil
}

func (e *Encoder) single(id int64, index, count int32, buf []byte) *Segment {
	end := int((index + 1) * MAX_SEGMENT_SIZE)
	if index == count - 1 {
		end = len(buf)
	}
	return &Segment {
		Id:      proto.Int64(id),
		Index:   proto.Int32(index),
		Total:   proto.Int32(count),
		Ack:     proto.Int32(0),
		Body:    buf[index * MAX_SEGMENT_SIZE : end],
	}
}