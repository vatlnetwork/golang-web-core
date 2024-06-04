package domain

import (
	"time"

	"github.com/google/uuid"
)

type Theme string

const (
	DarkTheme   Theme = "dark"
	LightTheme  Theme = "light"
	SystemTheme Theme = "system"
)

type User struct {
	Id         string    `json:"id" bson:"id"`
	Email      string    `json:"email" bson:"email"`
	Username   string    `json:"username" bson:"username"`
	Password   string    `json:"password" bson:"encryptedPassword"`
	IsAdmin    bool      `json:"isAdmin" bson:"isAdmin"`
	SignUpDate int64     `json:"signUpDate" bson:"signUpDate"`
	LastSignIn int64     `json:"lastSignIn" bson:"lastSignIn"`
	Address    Address   `json:"address" bson:"address"`
	Image      MediaFile `json:"image" bson:"image"`
	Theme      Theme     `json:"theme" bson:"theme"`
	PhoneNo    int64     `json:"phoneNo" bson:"phoneNo"`
	ThemeColor Color     `json:"themeColor" bson:"themeColor"`
}

func NewUser(email, username, password string, isAdmin bool) User {
	return User{
		Id:         uuid.NewString(),
		Email:      email,
		Username:   username,
		Password:   password,
		IsAdmin:    isAdmin,
		SignUpDate: time.Now().UnixMilli(),
		LastSignIn: time.Now().UnixMilli(),
		Theme:      SystemTheme,
		ThemeColor: Color{
			R: 68,
			B: 138,
			G: 255,
		},
	}
}
