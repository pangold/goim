package msg

import (
	"testing"
)

func TestMessage_Merge(t *testing.T) {

	body := func(len int) []byte {
		buf := make([]byte, len)
		for i := 0; i < len; i++ {
			buf[i] = 'a'
		}
		return buf
	}

	msg := NewMessage()
	mt := NewMessageT(MSG_TEXT, "10001", "10002", body(200000000))

	msg.SetSplitHandler(func(data []byte) {
		// msg.Merge(data)
		length := len(data)
		msg.Merge(data[:length / 2])
		msg.Merge(data[length / 2 :])
	})

	var mt2 *MessageT
	msg.SetMessageHandler(func(msg *MessageT) {
		mt2 = msg
	})

	msg.SetAckHandler(func(msg *MessageT){

	})

	if err := msg.Split(mt); err != nil {
		t.Error(err)
	}

	if mt.UserId != mt2.UserId {
		t.Errorf("unexpected user, %s : %s", mt.UserId, mt2.UserId)
	}
	if mt.TargetId != mt2.TargetId {
		t.Errorf("unexpected target, %s : %s", mt.TargetId, mt2.TargetId)
	}
	if string(mt.Body) != string(mt2.Body) {
		t.Errorf("unexpected body")
	}
}
