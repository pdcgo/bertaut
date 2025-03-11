package mock_http

type User struct {
	Name string
}

// bertaut_api: /users
// test doc
type UserService interface {
	// method: get
	Item(userID uint) (User, error)
	// method: post
	CreateUser() error
}
