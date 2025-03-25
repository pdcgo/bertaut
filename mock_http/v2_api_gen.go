package mock_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUserServiceApi(srv UserService, g *gin.RouterGroup) {
	g.Handle(http.MethodGet, "/users/item", func(ctx *gin.Context) {
	})
	g.Handle(http.MethodGet, "/users/info", func(ctx *gin.Context) {
	})
	g.Handle(http.MethodPost, "/users/create_user", func(ctx *gin.Context) {
	})
	g.Handle(http.MethodDelete, "/users/delete_user", func(ctx *gin.Context) {
	})
	g.Handle(http.MethodPost, "/users/get_info", func(ctx *gin.Context) {
	})
	g.Handle(http.MethodGet, "/users/get_role", func(ctx *gin.Context) {
	})
}
func RegisterWareServiceApi(srv WareService, g *gin.RouterGroup) {
	g.Handle(http.MethodGet, "/wares/item", func(ctx *gin.Context) {
	})
}
