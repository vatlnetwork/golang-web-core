package usersdb

import (
	"golang-web-core/src/domain"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserRecord struct {
	Id                string         `bson:"id"`
	Email             string         `bson:"email"`
	Username          string         `bson:"username"`
	EncryptedPassword string         `bson:"encryptedPassword"`
	IsAdmin           bool           `bson:"isAdmin"`
	SignUpDate        int64          `bson:"signUpDate"`
	LastSignIn        int64          `bson:"lastSignIn"`
	Address           domain.Address `bson:"address"`
	ImageId           string         `bson:"imageId"`
	Theme             domain.Theme   `bson:"theme"`
	PhoneNo           int64          `bson:"phoneNo"`
	ThemeColor        domain.Color   `bson:"themeColor"`
}

func UserRecordFromDomain(user domain.User) (UserRecord, error) {
	// encrypt given password
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserRecord{}, err
	}
	// make sure email is valid
	_, err = mail.ParseAddress(user.Email)
	if err != nil {
		return UserRecord{}, err
	}
	// check to see if passed is password is already decrypted, if so, set the encrypted password to the already existing value
	_, err = bcrypt.Cost([]byte(user.Password))
	if err == nil {
		encryptedPassword = []byte(user.Password)
	}
	return UserRecord{
		Id:                user.Id,
		Email:             strings.ToLower(user.Email),
		Username:          user.Username,
		EncryptedPassword: string(encryptedPassword),
		IsAdmin:           user.IsAdmin,
		SignUpDate:        user.SignUpDate,
		LastSignIn:        user.LastSignIn,
		Address:           user.Address,
		Theme:             user.Theme,
		ImageId:           user.Image.Id,
		PhoneNo:           user.PhoneNo,
		ThemeColor:        user.ThemeColor,
	}, nil
}
