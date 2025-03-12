package mock_http

import "github.com/pdcgo/bertaut/mock_http/models"

type User struct {
	Name string
}

// bertaut_api: /users
// test doc
type UserService interface {
	// method: get
	// Deprecated: asdasdasd
	Item(userID uint) (User, error)
	// method: get
	Info() (*models.UserInfo, error)
	// method: post
	CreateUser() error
	// method: delete
	DeleteUser(user User) error
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
