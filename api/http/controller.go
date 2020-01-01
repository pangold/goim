package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/pangold/goim/api/front"
	"net/http"
	"strings"
)

type Controller struct {
	front *front.Server
}

func NewController(front *front.Server) *Controller {
	return &Controller{
		front: front,
	}
}

func (c *Controller) List(ctx *gin.Context) {
	// simple encode,
	// to decode: strings.Split(s, ",")
	ctx.String(http.StatusOK, strings.Join(c.front.GetOnlineUserIds(), ","))
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
	token := c.front.GetOnlineTokenByUserId(uid)
	if token == "" {
		ctx.String(http.StatusNotFound, "uid(%s) is not online", uid)
		return
	}
	// get data that needs to be sent in http request body
	data := make([]byte, 0)
	_, err := ctx.Request.Body.Read(data)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	c.front.Send(token, data)
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
	c.front.Broadcast(data)
	ctx.Status(http.StatusOK)
}

func (c *Controller) Online(ctx *gin.Context) {
	// user id to token
	uid := ctx.Query("uid")
	if uid == "" {
		ctx.String(http.StatusBadRequest, "uid could not be null")
		return
	}
	token := c.front.GetOnlineTokenByUserId(uid)
	if token == "" {
		ctx.Status(http.StatusNotFound)
	}
	ctx.Status(http.StatusOK)
}

func (c *Controller) Kick(ctx *gin.Context) {
	// user id to token
	uid := ctx.Query("uid")
	if uid == "" {
		ctx.String(http.StatusBadRequest, "uid could not be null")
		return
	}
	token := c.front.GetOnlineTokenByUserId(uid)
	if token == "" {
		ctx.String(http.StatusNotFound, "uid(%s) is not online", uid)
	}
	c.front.Remove(token)
	ctx.Status(http.StatusOK)
}


