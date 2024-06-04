package domain

type UserRepository interface {
	CreateUser(user User) (User, error)
	Find(id string) (User, error)
	FindByEmail(email string) (User, error)
	QueryByEmail(search string) ([]User, error)
	QueryByUsername(search string) ([]User, error)
	Update(user User) error
	Delete(id string) error
	MatchUser(emailOrUsername, password string) (User, error)
}
