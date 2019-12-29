package codec

import "time"

const (
	// these only available when action = ACTION_CHAT & ACTION_UPLOAD & ACTION_SYSTEM
	MSG_TEXT     = 0
	MSG_IMAGE    = 1
	MSG_VIDEO    = 2
	MSG_AUDIO    = 3
	MSG_SYSTEM   = 4
	MSG_FILE     = 5

	ACTION_CHAT   = 0 // chatting by using TCP/Websocket
	ACTION_UPLOAD = 1 // p2p chatting by using UDP
	ACTION_SYSTEM = 2

	ACTION_FRIEND_REQUESTED = 10
	ACTION_FRIEND_REJECTED  = 11
	ACTION_FRIEND_ACCEPTED  = 12
	ACTION_FRIEND_RECOMMENT = 13 // recomment a friend to others
	ACTION_FRIEND_BROKEUP   = 14 //

	ACTION_GROUP_CREATED        = 20
	ACTION_GROUP_JOIN_REQUESTED = 21
	ACTION_GROUP_JOIN_REJECTED  = 22
	ACTION_GROUP_JOIN_ACCEPTED  = 23
	ACTION_GROUP_JOINED         = 24 // a new member joined
	ACTION_GROUP_LEFT           = 25 // a member's left
	ACTION_GROUP_KICKED         = 26 // a member's been kicked.
	ACTION_GROUP_DISMISSED      = 26 // group dismissed

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
		Action:   ACTION_CHAT,
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
		Action:   ACTION_CHAT,
		Ack:      ACK_NONE,
		Type:     int32(t),
		Body:     body,
	}
}

func NewSystemMessageT(t int, uid string, body []byte) *MessageT {
	return &MessageT {
		Id:       time.Now().UnixNano(),
		UserId:   uid,
		Action:   ACTION_SYSTEM,
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
