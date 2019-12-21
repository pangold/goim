package common

import (
	"encoding/binary"
	"reflect"
)

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
		return binary.BigEndian.Uint16(buf)
	case 4:
		return binary.BigEndian.Uint32(buf)
	case 8:
		return binary.BigEndian.Uint64(buf)
	}
	panic("invalid size")
}
