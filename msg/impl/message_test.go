package impl

import (
	"github.com/golang/protobuf/proto"
	message "gitlab.com/pangold/goim/msg/protobuf"
	"testing"
	"time"
)

func TestProtobuf_ReceivedData(t *testing.T) {

	//
	body := func(len int) []byte {
		buf := make([]byte, len)
		for i := 0; i < len; i++ {
			buf[i] = 'a'
		}
		return buf
	}

	msg := &message.Message{
		UserId:               proto.String("10001"),
		TargetId:             proto.String("10002"),
		GroupId:              nil,
		Type:                 (*message.Message_MessageType)(proto.Int32(int32(message.Message_TEXT))),
		Ack:                  (*message.Message_AckType)(proto.Int32(int32(message.Message_NONE))),
		// Body:                 body(3000000),
		Body:                 body(200000000),
		Time:                 proto.Int64(time.Now().Unix() >> 1),
	}

	p := NewMessage()

	p.SetSplitHandler(func(data []byte) {
		//p.SetReceived(data)
		length := len(data)
		p.Merge(data[:length / 2])
		p.Merge(data[length / 2 :])
	})

	var msg2 *message.Message
	p.SetMessageHandler(func(msg *message.Message) {
		msg2 = msg
	})

	if err := p.Split(msg); err != nil {
		t.Error(err)
	}

	//time.Sleep(time.Second * 3)

	if msg.GetUserId() != msg2.GetUserId() {
		t.Errorf("unexpected user, %s : %s", msg.GetUserId(), msg2.GetUserId())
	}
	if msg.GetTargetId() != msg2.GetTargetId() {
		t.Errorf("unexpected target, %s : %s", msg.GetTargetId(), msg2.GetTargetId())
	}
	b1 := string(msg.GetBody())
	b2 := string(msg2.GetBody())
	if b1 != b2 {
		t.Errorf("unexpected data, %s : %s", b1, b2)
	}
}
