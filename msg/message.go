package msg

import (
	"gitlab.com/pangold/goim/msg/impl"
	"gitlab.com/pangold/goim/msg/interfaces"
	message "gitlab.com/pangold/goim/msg/protobuf"
)

type Message struct {
	msg interfaces.Message
}

func NewMessage() *Message {
	return &Message {
		msg: impl.NewMessage(),
	}
}

func (p *Message) SetSplitHandler(h func([]byte)) {
	p.msg.SetSplitHandler(h)
}

func (p *Message) SetMessageHandler(h func(*message.Message)) {
	// TODO handle here, and convert to another type
	p.msg.SetMessageHandler(h)
}

func (p *Message) SetAckHandler(h func(*message.Message)) {
	p.msg.SetAckHandler(h)
}

func (p *Message) Split(msg *message.Message) error {
	return p.msg.Split(msg)
}

func (p *Message) Merge(data []byte) int {
	return p.msg.Merge(data)
}

