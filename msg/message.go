package msg

import (
	"gitlab.com/pangold/goim/msg/impl"
	message "gitlab.com/pangold/goim/msg/protobuf"
)

type Message struct {
	msg            *impl.Message
	ackHandler     func(*MessageT)
	messageHandler func(*MessageT)
}

func NewMessage() *Message {
	m := &Message {
		msg: impl.NewMessage(),
	}
	m.msg.SetAckHandler(m.handleAck)
	m.msg.SetMessageHandler(m.handleMessage)
	return m
}

func (m *Message) SetSplitHandler(h func([]byte)) {
	m.msg.SetSplitHandler(h)
}

func (m *Message) SetMessageHandler(h func(*MessageT)) {
	m.messageHandler = h
}

func (m *Message) SetAckHandler(h func(*MessageT)) {
	m.ackHandler = h
}

func (m *Message) Merge(data []byte) int {
	return m.msg.Merge(data)
}

func (m *Message) Split(msg *MessageT) error {
	mm := &message.Message {
		Id:       &msg.Id,
		UserId:   &msg.UserId,
		TargetId: &msg.TargetId,
		GroupId:  &msg.GroupId,
		Type:     (*message.Message_MessageType)(&msg.Type),
		Ack:      (*message.Message_AckType)(&msg.Ack),
		Body:     msg.Body,
	}
	return m.msg.Split(mm)
}

func (m *Message) handleAck(msg *message.Message) {
	mm := &MessageT {
		Id:       msg.GetId(),
		UserId:   msg.GetUserId(),
		TargetId: msg.GetTargetId(),
		GroupId:  msg.GetGroupId(),
		Type:     int32(msg.GetType()),
		Ack:      int32(msg.GetAck()),
		Body:     msg.Body,
	}
	m.ackHandler(mm)
}

func (m *Message) handleMessage(msg *message.Message) {
	mm := &MessageT {
		Id:       msg.GetId(),
		UserId:   msg.GetUserId(),
		TargetId: msg.GetTargetId(),
		GroupId:  msg.GetGroupId(),
		Type:     int32(msg.GetType()),
		Ack:      int32(msg.GetAck()),
		Body:     msg.Body,
	}
	m.messageHandler(mm)
}