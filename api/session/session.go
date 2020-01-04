package session

// A local session pool for current node.

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "gitlab.com/pangold/goim/api/grpc/proto"
	"gitlab.com/pangold/goim/utils"
	"log"
)

type Sessions struct {
	sessions map[string]*pb.Session
}

func NewSession(token string) *pb.Session {
	s := &pb.Session{
		Token:   &token,
		UserId:   proto.String(""),
		UserName: proto.String(""),
		ClientId: proto.String(""),
		NodeName: proto.String(""),
	}
	if err := utils.ExplainJwt(token, s.ClientId, s.UserId, s.UserName); err != nil {
		return nil
	}
	return s
}

func NewSessions() *Sessions {
	return &Sessions {
		sessions: make(map[string]*pb.Session),
	}
}

func (sp *Sessions) Add(token string, filter func(*pb.Session)error) error {
	s := NewSession(token)
	if s == nil {
		return fmt.Errorf("invalid token: %s", token)
	}
	if err := filter(s); err != nil {
		return err
	}
	log.Printf("new connection: cid = %s, uid = %s, name = %s", s.GetClientId(), s.GetUserId(), s.GetUserName())
	sp.sessions[s.GetUserId()] = s
	return nil
}

func (sp *Sessions) Remove(token string) *pb.Session {
	log.Printf("disconnection: token %s", token)
	tmp := NewSession(token)
	_, ok := sp.sessions[tmp.GetUserId()]
	if ok {
		delete(sp.sessions, token)
	}
	return tmp
}

func (sp *Sessions) Clear() {
	sp.sessions = nil
}

func (sp *Sessions) GetTokenByUserId(uid string) string {
	if s, ok := sp.sessions[uid]; ok {
		return s.GetToken()
	}
	return ""
}

func (sp *Sessions) GetUserIds() (res []string) {
	for _, s := range sp.sessions {
		res = append(res, s.GetUserId())
	}
	return res
}

func (sp *Sessions) GetTokens() (res []string) {
	for _, s := range sp.sessions {
		res = append(res, s.GetToken())
	}
	return res
}


