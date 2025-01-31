package types

import "time"

type RegisterUserPayload struct {
	Username string `json:"username" validate:"require"`
	Email    string `json:"email" validate:"require,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user User) error
}

type mockUserStore struct{}

func GetUserByEmail(email string) (*User, error) {
	return nil, nil
}

func GetUserByID(id int) (*User, error) {
	return nil, nil
}

func CreateUser(user User) error {
	return nil
}
