package system

import (
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/front"
	"gitlab.com/pangold/goim/protocol"
	"log"
)

type Chat struct {
	sessions *session.Sessions
	front *front.Server
}

func NewChatService(f *front.Server, s *session.Sessions) *Chat {
	return &Chat {
		sessions: s,
		front: f,
	}
}

// point to point chat, very simple.
func (this *Chat) Dispatch(msg *protocol.Message) {
	target := this.sessions.GetTokenByUserId(msg.TargetId)
	if target == "" {
		log.Println("target is not online.")
		return
	}
	if err := this.front.Send(target, msg); err != nil {
		log.Printf(err.Error())
	}
}

// I have already had session pool that be passed in by constructor NewChatService,
// So, I don't need SessionIn and SessionOut actually.
// func (this *Chat) SessionIn(session *protocol.Session) {
//
// }
//
// func (this *Chat) SessionOut(session *protocol.Session) {
//
// }


