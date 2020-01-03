package v1

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/codec/protobuf"
	"log"
	"time"
)

const (
	MAX_SEGMENT_SIZE = 1024
	WAIT_RESEND_TIME = time.Second * 3
)

type Encoder struct {
	segmentHandler  func(interface{}, *protobuf.Segment)
	resendHandler  *func(interface{}, *protobuf.Segment)
	acks            map[int64]*acknowledge //
}

type acknowledge struct {
	id       int64
	segments []*protobuf.Segment
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

func (e *Encoder) SetSegmentHandler(handler func(interface{}, *protobuf.Segment)) {
	e.segmentHandler = handler
}

func (e *Encoder) SetResendHandler(handler func(interface{}, *protobuf.Segment)) {
	e.resendHandler = &handler
}

func (e *Encoder) ResendEnabled() bool {
	return e.resendHandler != nil
}

// without ack, segment will be held in acknowledge list
func (e *Encoder) SetAckSegment(seg *protobuf.Segment) {
	if ack, ok := e.acks[seg.GetId()]; ok && e.ResendEnabled() {
		// being ack, clear
		// optimize: release by gc, or release immediately yourself
		ack.segments[seg.GetIndex()] = nil
	}
}

func (e *Encoder) Send(conn interface{}, msg *protobuf.Message) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return errors.New("split error: " + err.Error())
	}
	pages := int32(len(buf) / (MAX_SEGMENT_SIZE + 1)) + 1
	var ack *acknowledge = nil
	if e.ResendEnabled() {
		ack = &acknowledge{id: msg.GetId()}
		e.acks[msg.GetId()] = ack
		ack.segments = make([]*protobuf.Segment, pages)
	}
	for index := int32(0); index < pages; index++ {
		seg := e.single(msg.GetId(), index, pages, buf)
		if ack != nil {
			ack.segments[index] = seg
		}
		e.segmentHandler(conn, seg)
	}
	// ack.segments = segments
	// check timeout, trigger resend callback
	// resend works unless s.resendHandler is being specific.
	if ack != nil && e.ResendEnabled() {
		ack.timer = time.NewTimer(WAIT_RESEND_TIME)
		// FIXME: unstable retry times
		ack.retry = 3 * len(ack.segments)
		go e.timeout(conn, ack)
	}
	return nil
}

func (e *Encoder) single(id int64, index, pages int32, buf []byte) *protobuf.Segment {
	end := int((index + 1) * MAX_SEGMENT_SIZE)
	if index == pages - 1 {
		end = len(buf)
	}
	return &protobuf.Segment{
		Id:      proto.Int64(id),
		Index:   proto.Int32(index),
		Total:   proto.Int32(pages),
		Ack:     proto.Int32(0),
		Body:    buf[index * MAX_SEGMENT_SIZE : end],
	}
}

func (e *Encoder) timeout(conn interface{}, ack *acknowledge) {
	select {
	case <-ack.timer.C:
		e.resend(conn, ack)
	}
}

func (e *Encoder) resend(conn interface{}, ack *acknowledge) {
	var count = 0
	for _, seg := range ack.segments {
		if seg == nil {
			continue
		}
		if e.resendHandler != nil {
			(*e.resendHandler)(conn, seg)
		}
		if ack.retry == 0 {
			log.Printf("segment retry failed: message %d, body: %s", ack.id, ack.segments)
		} else {
			go e.timeout(conn, ack)
		}
		ack.retry--
		count++
	}
	// nothing left, clear it
	if count == 0 {
		delete(e.acks, ack.id)
	}
}
