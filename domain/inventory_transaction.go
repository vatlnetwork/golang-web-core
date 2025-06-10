package domain

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
