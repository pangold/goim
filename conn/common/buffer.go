package common

//import (
//	"gitlab.com/pangold/goim/msg"
//)
//
//type Buffer struct {
//	data []byte
//	handler func(*msg.Message)
//}
//
//func (b *Buffer) Add(data []byte) (*msg.Message, error) {
//	m := &msg.Message{}
//	b.data = append(b.data, data...)
//	if _, err := m.Deserialize(b.data); err != nil {
//		return nil, err
//	}
//	return m, nil
//}