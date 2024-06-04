package usersmock

import (
	"fmt"
	"golang-web-core/src/domain"
	"golang-web-core/src/srv/srverr"
	"strings"

	"github.com/google/uuid"
)

type MockUserRepo struct{}

func NewMockUserRepo() MockUserRepo {
	return MockUserRepo{}
}

var users []domain.User = []domain.User{}

func (r MockUserRepo) CreateUser(user domain.User) (domain.User, error) {
	user.Id = uuid.NewString()
	for _, u := range users {
		if u.Email == user.Email {
			return domain.User{}, fmt.Errorf("unable to create user with email: %v, user with that email already exists", user.Email)
		}
		if u.Username == user.Username {
			return domain.User{}, fmt.Errorf("unable to create user with username: %v, user with that username already exists", user.Username)
		}
	}
	users = append(users, user)
	return user, nil
}

func (r MockUserRepo) Find(id string) (domain.User, error) {
	for _, user := range users {
		if user.Id == id {
			return user, nil
		}
	}
	return domain.User{}, fmt.Errorf("unable to find user with id: %v", id)
}

func (r MockUserRepo) FindByEmail(email string) (domain.User, error) {
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}
	return domain.User{}, fmt.Errorf("unable to find a user with email: %v", email)
}

func (r MockUserRepo) QueryByEmail(search string) ([]domain.User, error) {
	res := []domain.User{}
	for _, user := range users {
		if strings.Contains(user.Email, search) {
			res = append(res, user)
		}
	}
	return res, nil
}

func (r MockUserRepo) QueryByUsername(search string) ([]domain.User, error) {
	res := []domain.User{}
	for _, user := range users {
		if strings.Contains(user.Username, search) {
			res = append(res, user)
		}
	}
	return res, nil
}

func (r MockUserRepo) Update(user domain.User) error {
	userIndex := -1
	for i := 0; i < len(users); i++ {
		if users[i].Id == user.Id {
			userIndex = i
		}
	}
	if userIndex < 0 || userIndex > len(users)-1 {
		return fmt.Errorf("unable to find user with id: %v", user.Id)
	}
	users[userIndex] = user
	return nil
}

func (r MockUserRepo) Delete(id string) error {
	res := []domain.User{}
	for _, user := range users {
		if user.Id != id {
			res = append(res, user)
		}
	}
	users = res
	return nil
}

func (r MockUserRepo) MatchUser(emailOrUsername, password string) (domain.User, error) {
	emailMatches := []domain.User{}
	usernameMatches := []domain.User{}
	for _, user := range users {
		if user.Email == emailOrUsername {
			emailMatches = append(emailMatches, user)
		}
		if user.Username == emailOrUsername {
			usernameMatches = append(usernameMatches, user)
		}
	}
	for _, match := range emailMatches {
		if match.Password == password {
			return match, nil
		}
	}
	for _, match := range usernameMatches {
		if match.Password == password {
			return match, nil
		}
	}
	return domain.User{}, fmt.Errorf(srverr.InvalidCredentials)
}
