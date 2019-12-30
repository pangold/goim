package protobuf

import (
	"github.com/golang/protobuf/proto"
	"testing"
	"time"
)

func TestCodec_1(t *testing.T) {

	//
	body := func(len int) []byte {
		buf := make([]byte, len)
		for i := 0; i < len; i++ {
			buf[i] = 'a'
		}
		return buf
	}

	msg := &Message{
		Id:                   proto.Int64(time.Now().UnixNano()),
		UserId:               proto.String("10001"),
		TargetId:             proto.String("10002"),
		GroupId:              nil,
		Action:               proto.Int32(0),
		Ack:                  proto.Int32(0),
		Type:                 proto.Int32(0),
		Body:                 body(200000000),
	}

	c := NewCodec()

	c.SetEncodeHandler(func(conn interface{}, data []byte) {
		//c.SetReceived(data)
		length := len(data)
		c.Decode(conn, data[:length / 2])
		c.Decode(conn, data[length / 2 :])
	})

	var msg2 *Message
	c.SetDecodeHandler(func(conn interface{}, msg *Message) {
		msg2 = msg
	})

	if err := c.Encode(nil, msg); err != nil {
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
