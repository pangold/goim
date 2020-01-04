package codec

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/codec/protobuf"
	"gitlab.com/pangold/goim/conn/interfaces"
)

type Decoder struct {
	MsgHandler     func(interfaces.Conn, *protobuf.Message)
	segments       map[int64]*segments
}

type segments struct {
	segs  []*protobuf.Segment
	count int
}

func NewDecoder() *Decoder {
	return &Decoder{
		MsgHandler: nil,
		segments:   make(map[int64]*segments),
	}
}

// because of reset mechanism, seg may be exist
func (d *Decoder) Decode(conn interfaces.Conn, seg *protobuf.Segment) error {
	// optimize for single segment
	if seg.GetTotal() == 1 {
		return d.single(conn, seg.GetBody())
	}
	// for multi segments
	if _, ok := d.segments[seg.GetId()]; !ok {
		d.segments[seg.GetId()] = &segments{count: 0, segs: make([]*protobuf.Segment, seg.GetTotal())}
	}
	// check if this segment is resent,
	// but another segment that with the same id/index/total
	// had already been received
	ss := d.segments[seg.GetId()]
	if ss.segs[seg.GetIndex()] != nil {
		return fmt.Errorf("%d(%d/%d) had already confirmed", seg.GetId(), seg.GetIndex(), seg.GetTotal())
	}
	ss.count++
	ss.segs[seg.GetIndex()] = seg
	if ss.count == int(seg.GetTotal()) {
		res := d.multi(conn, ss.segs)
		delete(d.segments, seg.GetId())
		return res
	}
	return nil
}

// []*message.Segment is not in order,
// The size of body of segments are the same, except the last segment
func (d *Decoder) multi(conn interfaces.Conn, sl []*protobuf.Segment) error {
	buf := make([]byte, MAX_SEGMENT_SIZE* (len(sl) - 1))
	for i := 0; i < len(sl) - 1; i++ {
		if len(sl[i].GetBody()) > MAX_SEGMENT_SIZE {
			return errors.New("unexpected segment size")
		}
		begin := int(sl[i].GetIndex()) * MAX_SEGMENT_SIZE
		end := int(sl[i].GetIndex() + 1) * MAX_SEGMENT_SIZE
		copy(buf[begin:end], sl[i].GetBody())
	}
	// last one is rarely full.
	buf = append(buf, sl[len(sl) - 1].GetBody()...)
	return d.single(conn, buf)
}

// A message with single one segment
func (m *Decoder) single(conn interfaces.Conn, buf []byte) error {
	msg := &protobuf.Message{}
	if err := proto.Unmarshal(buf, msg); err != nil {
		return errors.New("unmarshal error: " + err.Error())
	}
	m.MsgHandler(conn, msg)
	return nil
}
