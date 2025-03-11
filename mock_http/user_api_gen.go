package mock_http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterUserServiceApi(srv UserService, g *gin.RouterGroup) {
	g.Handle(http.MethodPost, "test", func(ctx *gin.Context) {
	})
}
