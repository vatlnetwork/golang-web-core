package domain

import (
	"fmt"
	"golang-web-core/srv/srverr"
	"net/http"
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

func NewUser(email, firstName, lastName, password string) (User, error) {
	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return User{}, srverr.New(fmt.Sprintf("password must be between %d and %d characters", minPasswordLength, maxPasswordLength), http.StatusBadRequest)
	}

	if !regexp.MustCompile(emailRegex).MatchString(email) {
		return User{}, srverr.New("invalid email address", http.StatusBadRequest)
	}

	user := User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, srverr.Wrap(err)
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
