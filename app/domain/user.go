package domain

import (
	"fmt"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLength = 6
	maxPasswordLength = 128
	emailRegex        = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

type User struct {
	Id                bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Email             string        `json:"email" bson:"email"`
	FirstName         string        `json:"firstName" bson:"firstName"`
	LastName          string        `json:"lastName" bson:"lastName"`
	EncryptedPassword []byte        `json:"encryptedPassword" bson:"encryptedPassword"`
	CreatedAt         time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt" bson:"updatedAt"`
}

func NewUser(email, password string) (User, error) {
	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return User{}, fmt.Errorf("password must be between %d and %d characters", minPasswordLength, maxPasswordLength)
	}

	if !regexp.MustCompile(emailRegex).MatchString(email) {
		return User{}, fmt.Errorf("invalid email address")
	}

	user := User{
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	user.EncryptedPassword = encryptedPassword
	return user, nil
}

func (u User) CheckPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.EncryptedPassword, []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (u *User) UpdatePassword(password string) error {
	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return fmt.Errorf("password must be between %d and %d characters", minPasswordLength, maxPasswordLength)
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.EncryptedPassword = encryptedPassword
	return nil
}

func (u *User) UpdateEmail(email string) error {
	if !regexp.MustCompile(emailRegex).MatchString(email) {
		return fmt.Errorf("invalid email address")
	}

	u.Email = email
	return nil
}
