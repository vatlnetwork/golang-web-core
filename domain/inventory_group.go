package domain

import "errors"

type InventoryGroup struct {
	Id          string          `json:"id"`
	UserId      string          `json:"userId"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Items       []InventoryItem `json:"items"`
}

func NewInventoryGroup(userId, name, description string) (InventoryGroup, error) {
	if userId == "" {
		return InventoryGroup{}, errors.New("userId is required")
	}

	if name == "" {
		return InventoryGroup{}, errors.New("name is required")
	}

	return InventoryGroup{
		Id:          "",
		UserId:      userId,
		Name:        name,
		Description: description,
	}, nil
}
