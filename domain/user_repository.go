package domain

const ErrorUserNotFound string = "user not found"
const ErrorUserAlreadyExists string = "user already exists"

type UserRepository interface {
	CreateUser(user User) (User, error)
	GetUser(userId string) (User, error)
	GetUserByEmail(email string) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(userId string) error
}
