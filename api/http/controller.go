package http

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"gitlab.com/pangold/goim/api/session"
	"gitlab.com/pangold/goim/front"
	"gitlab.com/pangold/goim/protocol"
	"net/http"
	"strings"
)

type Controller struct {
	front    *front.Server
	sessions *session.Sessions
}

func NewController(front *front.Server, ss *session.Sessions) *Controller {
	return &Controller{
		front:    front,
		sessions: ss,
	}
}

func (c *Controller) List(ctx *gin.Context) {
	// simple encode,
	// to decode: strings.Split(s, ",")
	ctx.String(http.StatusOK, strings.Join(c.sessions.GetUserIds(), ","))
}

// body: could be any word, such '&', '='
// it will influence others
// so, all the others params are query params
func (c *Controller) Send(ctx *gin.Context) {
	// user id to token
	uid := ctx.Query("uid")
	if uid == "" {
		ctx.String(http.StatusBadRequest, "uid could not be null")
		return
	}
	token := c.sessions.GetTokenByUserId(uid)
	if token == "" {
		ctx.String(http.StatusBadRequest, "uid(%s) is not online", uid)
		return
	}
	// get data that needs to be sent in http request body
	data := make([]byte, 0)
	if _, err := ctx.Request.Body.Read(data); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	// decode
	msg := &protocol.Message{}
	if err := proto.Unmarshal(data, msg); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	// send
	if err := c.front.Send(token, msg); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *Controller) Broadcast(ctx *gin.Context) {
	// get data that needs to be sent in http request body
	data := make([]byte, 0)
	_, err := ctx.Request.Body.Read(data)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	// decode
	msg := &protocol.Message{}
	if err := proto.Unmarshal(data, msg); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	c.front.Broadcast(msg)
	ctx.Status(http.StatusOK)
}

func (c *Controller) Online(ctx *gin.Context) {
	// user id to token
	uid := ctx.Query("uid")
	if uid == "" {
		ctx.String(http.StatusBadRequest, "uid could not be null")
		return
	}
	token := c.sessions.GetTokenByUserId(uid)
	if token == "" {
		ctx.String(http.StatusOK, "0")
	}
	ctx.String(http.StatusOK, "1")
}

func (c *Controller) Kick(ctx *gin.Context) {
	// user id to token
	uid := ctx.Query("uid")
	if uid == "" {
		ctx.String(http.StatusBadRequest, "uid could not be null")
		return
	}
	token := c.sessions.GetTokenByUserId(uid)
	if token == "" {
		ctx.String(http.StatusBadRequest, "uid(%s) is not online", uid)
	}
	c.front.Remove(token)
	ctx.Status(http.StatusOK)
}


