package domain

import "errors"

type InventoryItem struct {
	Id           string           `json:"id"`
	UserId       string           `json:"userId"`
	Description  string           `json:"description"`
	Locations    map[string]int64 `json:"locations"`
	InitialValue float64          `json:"initialValue"`
	Value        float64          `json:"value"`
}

func NewInventoryItem(userId, description string, initialValue float64) (InventoryItem, error) {
	if userId == "" {
		return InventoryItem{}, errors.New("userId is required")
	}

	if description == "" {
		description = "No description"
	}

	if initialValue < 0 {
		return InventoryItem{}, errors.New("initialValue must be greater than or equal to 0")
	}

	return InventoryItem{
		UserId:       userId,
		Description:  description,
		InitialValue: initialValue,
		Value:        initialValue,
	}, nil
}

func (i InventoryItem) AddQuantity(locationId string, quantity int64) (InventoryItem, error) {
	if locationId == "" {
		return InventoryItem{}, errors.New("locationId is required")
	}

	if quantity < 0 {
		return InventoryItem{}, errors.New("quantity must be greater than or equal to 0")
	}

	if _, ok := i.Locations[locationId]; !ok {
		i.Locations[locationId] = 0
	}

	i.Locations[locationId] += quantity
	return i, nil
}

func (i InventoryItem) RemoveQuantity(locationId string, quantity int64) (InventoryItem, error) {
	if locationId == "" {
		return InventoryItem{}, errors.New("locationId is required")
	}

	if quantity < 0 {
		return InventoryItem{}, errors.New("quantity must be greater than or equal to 0")
	}

	if _, ok := i.Locations[locationId]; !ok {
		return InventoryItem{}, errors.New("location not found")
	}

	i.Locations[locationId] -= quantity

	if i.Locations[locationId] < 0 {
		return InventoryItem{}, errors.New("new quantity cannot be negative")
	}

	return i, nil
}

func (i InventoryItem) UpdateValue(value float64) (InventoryItem, error) {
	if value < 0 {
		return InventoryItem{}, errors.New("value must be greater than or equal to 0")
	}

	i.Value = value
	return i, nil
}
