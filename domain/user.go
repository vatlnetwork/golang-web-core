package domain

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const ErrorUserNotFound string = "user not found"
const ErrorUserAlreadyExists string = "user already exists"

type UserRepository interface {
	CreateUser(user User) (User, error)
	GetUser(userId string) (User, error)
	GetUserByEmail(email string) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(userId string) error
}

const emailRegex = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
const minPasswordLength = 6
const maxPasswordLength = 128

const ErrorInvalidEmail = "invalid email"
const ErrorInvalidPassword = "password must be between 6 and 128 characters"

type User struct {
	Id                string    `json:"id"`
	Email             string    `json:"email"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	EncryptedPassword string    `json:"-"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	LastSignIn        time.Time `json:"lastSignIn"`
}

func NewUser(email, firstName, lastName, password string) (User, error) {
	if !regexp.MustCompile(emailRegex).MatchString(email) {
		return User{}, errors.New(ErrorInvalidEmail)
	}

	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return User{}, errors.New(ErrorInvalidPassword)
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	if firstName == "" && lastName == "" {
		firstName = "Anonymous"
		lastName = "User"
	}

	return User{
		Email:             email,
		FirstName:         firstName,
		LastName:          lastName,
		EncryptedPassword: string(encryptedPassword),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}, nil
}

func (u User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))
}

func (u *User) UpdatePassword(password, passwordConfirmation string) error {
	if password != passwordConfirmation {
		return errors.New("passwords do not match")
	}

	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return errors.New(ErrorInvalidPassword)
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.EncryptedPassword = string(encryptedPassword)

	return nil
}
