package session

import (
	"fmt"
	"gitlab.com/pangold/goim/api/middleware"
	"gitlab.com/pangold/goim/protocol"
	"log"
)

type Sessions struct {
	token middleware.Token
	sessions map[string]*protocol.Session
}

func NewSessions(token middleware.Token) *Sessions {
	return &Sessions {
		token: token,
		sessions: make(map[string]*protocol.Session),
	}
}

func (this *Sessions) GetToken() middleware.Token {
	return this.token
}

func (this *Sessions) ResetTokenExplainer(token middleware.Token) {
	this.token = token
}

func (this *Sessions) Add(token string, filter func(*protocol.Session)error) error {
	s := this.token.ExplainToken(token)
	if s == nil {
		return fmt.Errorf("invalid token: %s", token)
	}
	if err := filter(s); err != nil {
		return err
	}
	log.Printf("new connection: cid = %s, uid = %s, name = %s", s.ClientId, s.UserId, s.UserName)
	this.sessions[s.UserId] = s
	return nil
}

func (this *Sessions) Remove(token string) *protocol.Session {
	log.Printf("disconnection: token %s", token)
	tmp := this.token.ExplainToken(token)
	if _, ok := this.sessions[tmp.UserId]; ok {
		delete(this.sessions, token)
	}
	return tmp
}

func (this *Sessions) Clear() {
	// FIXME: manually
	this.sessions = nil
}

func (this *Sessions) GetTokenByUserId(uid string) string {
	if s, ok := this.sessions[uid]; ok {
		return s.Token
	}
	return ""
}

func (this *Sessions) GetUserIds() (res []string) {
	for _, s := range this.sessions {
		res = append(res, s.UserId)
	}
	return res
}

func (this *Sessions) GetTokens() (res []string) {
	for _, s := range this.sessions {
		res = append(res, s.Token)
	}
	return res
}


