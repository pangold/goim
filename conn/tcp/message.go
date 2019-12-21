package tcp

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
	magic    uint32
	kind     uint8 // Heartbeat, Goodbye, or Token message
	size     uint8
	body     []byte
	checksum uint32
}

func NewInternalMessage() *InternalMessage {
	return &InternalMessage{}
}

func NewHeartbeatMessage() *InternalMessage {
	return &InternalMessage{
		magic: MAGIC,
		kind:  HEARTBEAT,
		size:  0,
	}
}

func NewGoodbyeMessage() *InternalMessage {
	return &InternalMessage{
		magic: MAGIC,
		kind:  GOODBYE,
		size:  0,
	}
}

func NewTokenMessage(token []byte) *InternalMessage {
	// maximum size of token: 256 bytes
	if len(token) > 256 {
		panic("invalid token size")
	}
	return &InternalMessage{
		magic: MAGIC,
		kind:  TOKEN,
		size:  uint8(len(token)),
		body:  token,
	}
}

func (m *InternalMessage) Deserialize(data []byte) (*InternalMessage, int) {
	if len(data) < 6 {
		return nil, 0
	}
	m.magic = FromBytes(data[: 4]).(uint32)
	if m.magic != MAGIC {
		return nil, 0
	}
	m.kind = data[4]
	m.size = data[5]
	next := 6 + m.size
	if len(data) < int(next + 4) {
		return nil, 0
	}
	if m.size > 0 {
		m.body = data[6 : next]
	}
	m.checksum = FromBytes(data[next : next + 4]).(uint32)
	if crc32.Checksum(data[: next], table32) != m.checksum {
		return nil, 0
	}
 	return m, int(next) + 4
}

func (m *InternalMessage) Serialize() []byte {
	buf := make([]byte, 0)
	buf = append(buf, ToBytes(m.magic)...)
	buf = append(buf, m.kind)
	buf = append(buf, m.size)
	if m.size > 0 {
		buf = append(buf, m.body...)
	}
	m.checksum = crc32.Checksum(buf, table32)
	buf = append(buf, ToBytes(m.checksum)...)
	return buf
}