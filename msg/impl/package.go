package impl

// Extra work for TCP

import (
	"gitlab.com/pangold/goim/utils"
	"hash/crc32"
)

var (
	table32 = crc32.MakeTable(0xD5828281)
)

type Package struct {
	size     uint16
	body     []byte
	checksum uint32
}

func NewEmptyPackage() *Package {
	return &Package{}
}

func NewPackage(body []byte) *Package {
	return &Package {
		size:     uint16(len(body)),
		body:     body,
		checksum: crc32.Checksum(body, table32),
	}
}

func (m *Package) Deserialize(data []byte) (*Package, int) {
	m.size = utils.FromBytes(data[0:2]).(uint16)
	// unexpected length
	if m.size < 0 {
		return nil, -1
	}
	// not complete yet
	if int(m.size) > len(data) {
		return nil, 0
	}
	m.body = data[2 : 2 + m.size]
	m.checksum = utils.FromBytes(data[2 + m.size : ]).(uint32)
	// unexpected checksum
	if crc32.Checksum(m.body, table32) != m.checksum {
		return nil, -1
	}
 	return m, int(m.size) + 2 + 4
}

func (m *Package) Serialize() []byte {
	buf := make([]byte, 0)
	buf = append(buf, utils.ToBytes(m.size)...)
	buf = append(buf, m.body...)
	m.checksum = crc32.Checksum(m.body, table32)
	buf = append(buf, utils.ToBytes(m.checksum)...)
	return buf
}