package domain

type Inventory struct {
	Id        string              `json:"id"`
	UserId    string              `json:"userId"`
	Name      string              `json:"name"`
	Groups    []InventoryGroup    `json:"groups"`
	Locations []InventoryLocation `json:"locations"`
}
