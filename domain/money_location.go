package domain

import "errors"

const ErrorMoneyLocationNotFound = "money location not found"

type MoneyLocationRepository interface {
	CreateMoneyLocation(moneyLocation MoneyLocation) (MoneyLocation, error)
	GetMoneyLocation(id string) (MoneyLocation, error)
	GetMoneyLocationsForUser(userId string) ([]MoneyLocation, error)
	UpdateMoneyLocation(moneyLocation MoneyLocation) (MoneyLocation, error)
	DeleteMoneyLocation(id string) error
	DeleteAllMoneyLocationsForUser(userId string) error
}

type MoneyLocation struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
	Name   string `json:"name"`
}

func NewMoneyLocation(name string, userId string) (MoneyLocation, error) {
	if name == "" {
		name = "Unnamed Location"
	}

	if userId == "" {
		return MoneyLocation{}, errors.New("user id is required")
	}

	return MoneyLocation{
		Name:   name,
		UserId: userId,
	}, nil
}
