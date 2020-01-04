package codec

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/codec/protobuf"
)

const (
	MAX_SEGMENT_SIZE = 1024
)

type Encoder struct {

}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(msg *protobuf.Message) (res []*protobuf.Segment, err error) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return nil, errors.New("encode error: " + err.Error())
	}
	pages := int32(len(buf) / (MAX_SEGMENT_SIZE + 1)) + 1
	for index := int32(0); index < pages; index++ {
		seg := e.single(msg.GetId(), index, pages, buf)
		res = append(res, seg)
	}
	return res, nil
}

func (e *Encoder) single(id int64, index, pages int32, buf []byte) *protobuf.Segment {
	end := int((index + 1) * MAX_SEGMENT_SIZE)
	if index == pages - 1 {
		end = len(buf)
	}
	return &protobuf.Segment{
		Id:      &id,
		Index:   &index,
		Total:   &pages,
		Ack:     proto.Int32(0),
		Body:    buf[index * MAX_SEGMENT_SIZE : end],
	}
}