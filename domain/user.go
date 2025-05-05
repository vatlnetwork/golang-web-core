package domain

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const emailRegex = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"

type User struct {
	Id                string `json:"id"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"encryptedPassword"`
}

func NewUser(email, password string) (User, error) {
	if !regexp.MustCompile(emailRegex).MatchString(email) {
		return User{}, errors.New("invalid email")
	}

	if len(password) < 6 || len(password) > 128 {
		return User{}, errors.New("password must be between 6 and 128 characters")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	return User{
		Email:             email,
		EncryptedPassword: string(encryptedPassword),
	}, nil
}

func (u User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))
}
