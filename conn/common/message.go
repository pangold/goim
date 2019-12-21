package common

// Extra work for TCP

import (
	"hash/crc32"
)

const (
	MAGIC     = 0x20171007 //
	HEARTBEAT = 1          // Heartbeat(No body)
	GOODBYE   = 2          // Goodbye(No Body)
	TOKEN     = 3          // Token(With body, but only first message is available)
)

var (
	table32 = crc32.MakeTable(0xD5828281)
)

type InternalMessage struct {
	Magic    uint32
	Kind     uint8 // Heartbeat, Goodbye, or Token message
	Size     uint8
	Body     []byte
	Checksum uint32
}

func NewInternalMessage() *InternalMessage {
	return &InternalMessage{}
}

func NewHeartbeatMessage() *InternalMessage {
	return &InternalMessage{
		Magic: MAGIC,
		Kind:  HEARTBEAT,
		Size:  0,
	}
}

func NewGoodbyeMessage() *InternalMessage {
	return &InternalMessage{
		Magic: MAGIC,
		Kind:  GOODBYE,
		Size:  0,
	}
}

func NewTokenMessage(token []byte) *InternalMessage {
	// maximum size of token: 256 bytes
	if len(token) > 256 {
		panic("invalid token size")
	}
	return &InternalMessage{
		Magic: MAGIC,
		Kind:  TOKEN,
		Size:  uint8(len(token)),
		Body:  token,
	}
}

func (m *InternalMessage) Deserialize(data []byte) (*InternalMessage, int) {
	if len(data) < 6 {
		return nil, 0
	}
	m.Magic = FromBytes(data[: 4]).(uint32)
	if m.Magic != MAGIC {
		return nil, 0
	}
	m.Kind = data[4]
	m.Size = data[5]
	next := 6 + m.Size
	if len(data) < int(next + 4) {
		return nil, 0
	}
	if m.Size > 0 {
		m.Body = data[6 : next]
	}
	m.Checksum = FromBytes(data[next : next + 4]).(uint32)
	if crc32.Checksum(data[: next], table32) != m.Checksum {
		return nil, 0
	}
 	return m, int(next) + 4
}

func (m *InternalMessage) Serialize() []byte {
	buf := make([]byte, 0)
	buf = append(buf, ToBytes(m.Magic)...)
	buf = append(buf, m.Kind)
	buf = append(buf, m.Size)
	if m.Size > 0 {
		buf = append(buf, m.Body...)
	}
	m.Checksum = crc32.Checksum(buf, table32)
	buf = append(buf, ToBytes(m.Checksum)...)
	return buf
}