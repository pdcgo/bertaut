package mock_http

import (
	"github.com/gin-gonic/gin"
	"github.com/pdcgo/bertaut/mock_http/models"
)

type User struct {
	Name string
}

func (idn *User) BuildFromCtx(ctx *gin.Context) error {
	return nil
}

type ItemQuery struct {
	ItemID uint
}

type CCPayload struct {
	Dta string
}

type Identity struct{}

func (idn *Identity) BuildFromCtx(ctx *gin.Context) error {
	return nil
}

type CustomPayload struct {
	DD string
}

func (idn *CustomPayload) BuildFromCtx(ctx *gin.Context) error {
	return nil
}

// bertaut_api: /users
// test doc
type UserService interface {
	// method: get
	// Deprecated: asdasdasd
	Item(query *ItemQuery, payload *CCPayload) (User, error)
	// method: get
	Info() (*models.UserInfo, error)
	// method: post
	CreateUser(payload *CCPayload) error
	// method: delete
	DeleteUser(user User) error
	// method: post
	GetInfo(payload *CCPayload) (*models.UserInfo, error)
	// method: get
	GetRole(identity *Identity, custom *CustomPayload) error
}

// bertaut_api: /wares
// test doc
type WareService interface {
	// method: get
	// Deprecated: asdasdasd
	Item(query *ItemQuery, payload *CCPayload) (User, error)
}

// func(ctx *gin.Context) {
// 	var err error
// 	err = srv.CreateUser()

// 	if err != nil {
// 		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"message": err.Error(),
// 		})
// 	}

// }

// type PermissionQuery interface{}

// type AuthService interface {
// 	ApiQueryCheckPermission(identity Identity, query PermissionQuery) (bool, error)
// }
