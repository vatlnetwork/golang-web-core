package domain

type UserRepository interface {
	CreateUser(user User) (User, error)
	GetUser(userId string) (User, error)
	GetUserByEmail(email string) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(userId string) error
}
