package tcp

// Extra work for TCP

import (
	"encoding/binary"
	"hash/crc32"
	"reflect"
)

const (
	MAGIC = 0x20171007 //
	HEARTBEAT = 1      // Heartbeat(No body)
	GOODBYE = 2        // Goodbye(No Body)
	TOKEN = 3          // Token(With body, but only first message is available)
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

func ToBytes(num interface{}) []byte {
	t := reflect.TypeOf(num)
	var buf = make([]byte, t.Bits() / 8)
	switch t.Kind() {
	case reflect.Uint8:
		buf[0] = num.(uint8)
	case reflect.Uint16:
		binary.BigEndian.PutUint16(buf, num.(uint16))
	case reflect.Uint32:
		binary.BigEndian.PutUint32(buf, num.(uint32))
	case reflect.Uint64:
		binary.BigEndian.PutUint64(buf, num.(uint64))
	}
	return buf
}

func FromBytes(buf []byte) interface{} {
	switch len(buf) {
	case 1:
		return buf[0]
	case 2:
		return binary.BigEndian.Uint64(buf)
	case 4:
		return binary.BigEndian.Uint32(buf)
	case 8:
		return binary.BigEndian.Uint64(buf)
	}
	panic("invalid size")
}

func DeserializeInternalMessage(data []byte) (*InternalMessage, int) {
	if len(data) < 6 {
		return nil, 0
	}
	m := &InternalMessage{}
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

func SerializeInternalMessage(m *InternalMessage) []byte {
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

func SerializeHeartbeatMessage() []byte {
	m := &InternalMessage{}
	m.magic = MAGIC
	m.kind = HEARTBEAT
	m.size = 0
	return SerializeInternalMessage(m)
}

func SerializeGoodbyeMessage() []byte {
	m := &InternalMessage{}
	m.magic = MAGIC
	m.kind = GOODBYE
	m.size = 0
	return SerializeInternalMessage(m)
}

// maximum size of token: 256 bytes
func SerializeTokenMessage(token []byte) []byte {
	m := &InternalMessage{}
	m.magic = MAGIC
	m.kind = TOKEN
	m.size = uint8(len(token))
	m.body = token
	return SerializeInternalMessage(m)
}