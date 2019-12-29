package codec

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

	c := NewCodec()
	mt := NewMessageT(MSG_TEXT, "10001", "10002", body(200000000))

	c.SetEncodeHandler(func(data []byte) {
		// c.Merge(data)
		length := len(data)
		c.Decode(data[:length / 2])
		c.Decode(data[length / 2 :])
	})

	var mt2 *MessageT
	c.SetDecodeHandler(func(msg *MessageT) {
		mt2 = msg
	})

	if err := c.Encode(mt); err != nil {
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
