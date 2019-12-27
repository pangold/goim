package impl

import (
	"errors"
	"github.com/golang/protobuf/proto"
	message "gitlab.com/pangold/goim/msg/protobuf"
)

type Buffer struct {
	err     error
	handler func(*message.Segment)
	buf     []byte
}

func (b *Buffer) Error() error {
	return b.err
}

func (b *Buffer) SetSegmentHandler(handler func(*message.Segment)) {
	b.handler = handler
}

func (b *Buffer) Push(data []byte) {
	b.err = nil
	offset := 0
	b.buf = append(b.buf, data...)
	for {
		size := b.parse(b.buf[offset:])
		if size < 1 {
			break
		}
		offset += size
	}
	b.buf = b.buf[offset:]
}

// Too much copy...
// Does the middle Package neccessary ???????????????????
// Test it.
func (b *Buffer) parse(buf []byte) int {
	p := Package{}
	m, size := p.Deserialize(buf)
	if size == -1 {
		b.err = errors.New("unexpected data")
		return -1
	}
	if m == nil || size == 0 {
		return 0
	}
	seg := &message.Segment{}
	if err := proto.Unmarshal(m.body, seg); err != nil {
		b.err = errors.New("unmarshal error: " + err.Error())
		return -1
	}
	b.handler(seg)
	return size
}