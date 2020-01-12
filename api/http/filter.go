package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/pangold/goim/api/middleware"
	"net/http"
)

type Filter struct {
	token middleware.Token
}

func NewFilter(token middleware.Token) *Filter {
	return &Filter {
		token: token,
	}
}

func (this *Filter) Do(ctx *gin.Context) {
	token := ctx.GetHeader("token")
	if this.token.ExplainToken(token) == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Next()
}