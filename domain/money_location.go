package domain

import "errors"

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
