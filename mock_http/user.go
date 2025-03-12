package mock_http

import "github.com/pdcgo/bertaut/mock_http/models"

type User struct {
	Name string
}

type ItemQuery struct {
	ItemID uint
}

type CCPayload struct {
	Dta string
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
