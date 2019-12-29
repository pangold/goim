package codec

import "time"

const (
	MSG_TEXT     = 0
	MSG_IMAGE    = 1
	MSG_VIDEO    = 2
	MSG_AUDIO    = 3
	MSG_SYSTEM   = 4
	MSG_FILE     = 5

	ACK_NONE     = 0
	ACK_RECEIVED = 1
	ACK_READ     = 2
)

type MessageT struct {
	Id       int64
	UserId   string
	TargetId string
	GroupId  string
	Action   int32
	Ack      int32
	Type     int32
	Body     []byte
}

func NewMessageT(t int, uid, tid string, body []byte) *MessageT {
	return &MessageT {
		Id:       time.Now().UnixNano(),
		UserId:   uid,
		TargetId: tid,
		Ack:      ACK_NONE,
		Type:     int32(t),
		Body:     body,
	}
}

func NewGroupMessageT(t int, uid, gid string, body []byte) *MessageT {
	return &MessageT {
		Id:       time.Now().UnixNano(),
		UserId:   uid,
		GroupId:  gid,
		Ack:      ACK_NONE,
		Type:     int32(t),
		Body:     body,
	}
}

func NewAckMessageT(id int64) *MessageT {
	return &MessageT {
		Id:       id,
		Ack:      ACK_READ,
	}
}
