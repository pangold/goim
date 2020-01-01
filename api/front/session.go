package front

// A local session pool for current node.

import (
	"fmt"
	"gitlab.com/pangold/goim/utils"
	"log"
)

type Session struct {
	ClientId string
	UserId   string
	UserName string
	Token    string
}

type Sessions struct {
	sessions map[string]*Session
}

func NewSession(token string) *Session {
	s := &Session{Token: token}
	if err := utils.ExplainJwt(token, &s.ClientId, &s.UserId, &s.UserName); err != nil {
		return nil
	}
	return s
}

func NewSessions() *Sessions {
	return &Sessions {
		sessions: make(map[string]*Session),
	}
}

func (sp *Sessions) Add(token string, filter func(*Session)error) error {
	s := NewSession(token)
	if s == nil {
		return fmt.Errorf("invalid token: %s", token)
	}
	if err := filter(s); err != nil {
		return err
	}
	log.Printf("new connection: cid = %s, uid = %s, name = %s", s.ClientId, s.UserId, s.UserName)
	sp.sessions[s.UserId] = s
	return nil
}

func (sp *Sessions) Remove(token string) {
	log.Printf("disconnection: token %s", token)
	delete(sp.sessions, token)
}

func (sp *Sessions) Clear() {
	sp.sessions = nil
}

func (sp *Sessions) GetTokenByUserId(uid string) string {
	if s, ok := sp.sessions[uid]; ok {
		return s.Token
	}
	return ""
}

func (sp *Sessions) GetUserIds() (res []string) {
	for _, s := range sp.sessions {
		res = append(res, s.UserId)
	}
	return res
}

func (sp *Sessions) GetTokens() (res []string) {
	for _, s := range sp.sessions {
		res = append(res, s.Token)
	}
	return res
}


