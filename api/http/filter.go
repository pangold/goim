package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/pangold/goim/utils"
	"net/http"
)

func filter(ctx *gin.Context) {
	// check token
	token := ctx.GetHeader("token")
	var cid, uid, uname string
	if err := utils.ExplainJwt(token, &cid, &uid, &uname); err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// TODO: check permission
	ctx.Next()
}