package codec

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/codec/protobuf"
	"gitlab.com/pangold/goim/front/interfaces"
	"time"
)

const (
	MAX_SEGMENT_SIZE = 1024
	WAIT_RESEND_TIME = time.Second * 3
	RETRY_TIMES      = 3
)

type Encoder struct {
	acks            map[int64]*acknowledge //
	ResendEnabled   bool
}

type acknowledge struct {
	id       int64
	segments []*protobuf.Segment
	timer    *time.Timer
	retry    int
}

func NewEncoder() *Encoder {
	return &Encoder{
		acks:           make(map[int64]*acknowledge),
		ResendEnabled:  false,
	}
}

// without ack, segment will be held in acknowledge list
func (e *Encoder) SetAckSegment(seg *protobuf.Segment) {
	if ack, ok := e.acks[seg.GetId()]; ok && e.ResendEnabled {
		// being ack, clear
		// optimize: release by gc, or release immediately yourself
		ack.segments[seg.GetIndex()] = nil
	}
}

func (e *Encoder) Send(conn interfaces.Conn, msg *protobuf.Message) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return errors.New("split error: " + err.Error())
	}
	pages := int32(len(buf) / (MAX_SEGMENT_SIZE + 1)) + 1
	var ack *acknowledge = nil
	if e.ResendEnabled {
		ack = &acknowledge{id: msg.GetId(), retry: RETRY_TIMES}
		e.acks[msg.GetId()] = ack
		ack.segments = make([]*protobuf.Segment, pages)
	}
	for index := int32(0); index < pages; index++ {
		seg := e.single(msg.GetId(), index, pages, buf)
		if ack != nil {
			ack.segments[index] = seg
		}
		e.send(conn, seg)
	}
	// ack.segments = segments
	// check timeout, trigger resend callback
	// resend works unless s.resendHandler is being specific.
	if ack != nil && e.ResendEnabled {
		// FIXME: TTL
		ack.timer = time.NewTimer(WAIT_RESEND_TIME)
		go e.timeout(conn, ack)
	}
	return nil
}

func (e *Encoder) pos(index, pages int32, buf []byte) []byte {
	end := int((index + 1) * MAX_SEGMENT_SIZE)
	if index == pages - 1 {
		end = len(buf)
	}
	return buf[index * MAX_SEGMENT_SIZE : end]
}

func (e *Encoder) single(id int64, index, pages int32, buf []byte) *protobuf.Segment {
	return &protobuf.Segment{
		Id:      &id,
		Index:   &index,
		Total:   &pages,
		Ack:     proto.Int32(0),
		Body:    e.pos(index, pages, buf),
	}
}

func (e *Encoder) timeout(conn interfaces.Conn, ack *acknowledge) {
	select {
	case <-ack.timer.C:
		e.resend(conn, ack)
	}
}

func (e *Encoder) resend(conn interfaces.Conn, ack *acknowledge) {
	var resent = false
	for _, seg := range ack.segments {
		if seg != nil {
			e.send(conn, seg)
			resent = true
		}
		if ack.retry > 0 {
			ack.timer.Reset(WAIT_RESEND_TIME)
			go e.timeout(conn, ack)
		}
	}
	if resent || ack.retry == 0 {
		ack.timer.Stop()
		delete(e.acks, ack.id)
		return
	}
	ack.retry--
}

func (e *Encoder) send(conn interfaces.Conn, segment *protobuf.Segment) {
	if buf, err := proto.Marshal(segment); err == nil {
		conn.Send(buf)
	}
}
