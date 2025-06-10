package domain

import "errors"

type InventoryTransaction interface {
	GetType() string
	GetId() string
	GetUserId() string
	GetInventoryItemId() string
}

type InventoryTransactionTypeDecoder struct {
	Type InventoryTransactionType `json:"type"`
}

type InventoryTransactionType string

const (
	InventoryTransactionTypeQuantityChange InventoryTransactionType = "quantity_change"
	InventoryTransactionTypeValueChange    InventoryTransactionType = "value_change"
)

type InventoryTransactionQuantityChange struct {
	Type            InventoryTransactionType `json:"type"`
	Id              string                   `json:"id"`
	UserId          string                   `json:"userId"`
	InventoryItemId string                   `json:"inventoryItemId"`
	LocationId      string                   `json:"locationId"`
	QuantityChange  int64                    `json:"quantityChange"`
}

func NewInventoryTransactionQuantityChange(userId, inventoryItemId, locationId string, quantityChange int64) (InventoryTransactionQuantityChange, error) {
	if userId == "" {
		return InventoryTransactionQuantityChange{}, errors.New("userId is required")
	}

	if inventoryItemId == "" {
		return InventoryTransactionQuantityChange{}, errors.New("inventoryItemId is required")
	}

	if locationId == "" {
		return InventoryTransactionQuantityChange{}, errors.New("locationId is required")
	}

	if quantityChange == 0 {
		return InventoryTransactionQuantityChange{}, errors.New("quantityChange must be non-zero")
	}

	return InventoryTransactionQuantityChange{
		Type:            InventoryTransactionTypeQuantityChange,
		UserId:          userId,
		InventoryItemId: inventoryItemId,
		LocationId:      locationId,
		QuantityChange:  quantityChange,
	}, nil
}

func (t InventoryTransactionQuantityChange) GetType() string {
	return string(t.Type)
}

func (t InventoryTransactionQuantityChange) GetId() string {
	return t.Id
}

func (t InventoryTransactionQuantityChange) GetUserId() string {
	return t.UserId
}

func (t InventoryTransactionQuantityChange) GetInventoryItemId() string {
	return t.InventoryItemId
}

type InventoryTransactionValueChange struct {
	Type            InventoryTransactionType `json:"type"`
	Id              string                   `json:"id"`
	UserId          string                   `json:"userId"`
	InventoryItemId string                   `json:"inventoryItemId"`
	ValueChange     float64                  `json:"valueChange"`
}

func NewInventoryTransactionValueChange(userId, inventoryItemId string, valueChange float64) (InventoryTransactionValueChange, error) {
	if userId == "" {
		return InventoryTransactionValueChange{}, errors.New("userId is required")
	}

	if inventoryItemId == "" {
		return InventoryTransactionValueChange{}, errors.New("inventoryItemId is required")
	}

	if valueChange == 0 {
		return InventoryTransactionValueChange{}, errors.New("valueChange must be non-zero")
	}

	return InventoryTransactionValueChange{
		Type:            InventoryTransactionTypeValueChange,
		UserId:          userId,
		InventoryItemId: inventoryItemId,
		ValueChange:     valueChange,
	}, nil
}

func (t InventoryTransactionValueChange) GetType() string {
	return string(t.Type)
}

func (t InventoryTransactionValueChange) GetId() string {
	return t.Id
}

func (t InventoryTransactionValueChange) GetUserId() string {
	return t.UserId
}

func (t InventoryTransactionValueChange) GetInventoryItemId() string {
	return t.InventoryItemId
}
