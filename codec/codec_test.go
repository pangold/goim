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

	c.SetEncodeHandler(func(conn interface{}, data []byte) {
		// c.Merge(data)
		length := len(data)
		c.Decode(conn, data[:length / 2])
		c.Decode(conn, data[length / 2 :])
	})

	c.EnableResend(false)

	var mt2 *MessageT
	c.SetDecodeHandler(func(conn interface{}, msg *MessageT) {
		mt2 = msg
	})

	if err := c.Encode(nil, mt); err != nil {
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
