package domain

import "errors"

type InventoryLocation struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
	Name   string `json:"name"`
}

func NewInventoryLocation(userId, name string) (InventoryLocation, error) {
	if userId == "" {
		return InventoryLocation{}, errors.New("userId is required")
	}

	if name == "" {
		name = "Unnamed location"
	}

	return InventoryLocation{
		Id:     "",
		UserId: userId,
		Name:   name,
	}, nil
}
