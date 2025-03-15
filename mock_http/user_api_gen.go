package mock_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pdcgo/bertaut/mock_http/models"
)

func RegisterUserServiceApi(srv UserService, g *gin.RouterGroup, doc func(method string, path string, query any, payload any) error) {
	var err error
	g.Handle(http.MethodGet, "/users/item", func(ctx *gin.Context) {
		var result1 User
		var err error
		var query ItemQuery
		var payload CCPayload
		err = ctx.BindQuery(&query)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		err = ctx.BindJSON(&payload)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		result1, err = srv.Item(&query, &payload)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result1)
	})
	if doc != nil {
		err = doc(http.MethodGet, g.BasePath()+"/users", nil, nil)
		if err != nil {
			panic(err)
		}
	}
	g.Handle(http.MethodGet, "/users/info", func(ctx *gin.Context) {
		var result1 *models.UserInfo
		var err error
		result1, err = srv.Info()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result1)
	})
	if doc != nil {
		err = doc(http.MethodGet, g.BasePath()+"/users", nil, nil)
		if err != nil {
			panic(err)
		}
	}
	g.Handle(http.MethodPost, "/users/create_user", func(ctx *gin.Context) {
		var err error
		var payload CCPayload
		err = ctx.BindJSON(&payload)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		err = srv.CreateUser(&payload)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": ""})
	})
	if doc != nil {
		err = doc(http.MethodGet, g.BasePath()+"/users", nil, nil)
		if err != nil {
			panic(err)
		}
	}
	g.Handle(http.MethodDelete, "/users/delete_user", func(ctx *gin.Context) {
		var err error
		var param1 User
		err = param1.BuildFromCtx(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		err = srv.DeleteUser(param1)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": ""})
	})
	if doc != nil {
		err = doc(http.MethodGet, g.BasePath()+"/users", nil, nil)
		if err != nil {
			panic(err)
		}
	}
	g.Handle(http.MethodPost, "/users/get_info", func(ctx *gin.Context) {
		var result1 *models.UserInfo
		var err error
		var payload CCPayload
		err = ctx.BindJSON(&payload)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		result1, err = srv.GetInfo(&payload)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result1)
	})
	if doc != nil {
		err = doc(http.MethodGet, g.BasePath()+"/users", nil, nil)
		if err != nil {
			panic(err)
		}
	}
	g.Handle(http.MethodGet, "/users/get_role", func(ctx *gin.Context) {
		var err error
		var identity Identity
		var custom CustomPayload
		err = identity.BuildFromCtx(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		err = custom.BuildFromCtx(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		err = srv.GetRole(&identity, &custom)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": ""})
	})
	if doc != nil {
		err = doc(http.MethodGet, g.BasePath()+"/users", nil, nil)
		if err != nil {
			panic(err)
		}
	}
}
