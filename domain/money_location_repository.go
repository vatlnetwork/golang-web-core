package domain

const ErrorMoneyLocationNotFound = "money location not found"

type MoneyLocationRepository interface {
	CreateMoneyLocation(moneyLocation MoneyLocation) (MoneyLocation, error)
	GetMoneyLocation(id string) (MoneyLocation, error)
	GetMoneyLocationsForUser(userId string) ([]MoneyLocation, error)
	UpdateMoneyLocation(moneyLocation MoneyLocation) (MoneyLocation, error)
	DeleteMoneyLocation(id string) error
	DeleteAllMoneyLocationsForUser(userId string) error
}
