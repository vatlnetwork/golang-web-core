package domain

import (
	"errors"
	"time"
)

const ErrorTransactionNotFound string = "transaction not found"

type TransactionRepository interface {
	CreateTransaction(transaction Transaction) (Transaction, error)
	GetTransactionsForUser(userId string) ([]Transaction, error)
	GetTransactionsByLocation(locationId string) ([]Transaction, error)
	GetTransactionsByGroup(groupId string) ([]Transaction, error)
	GetTransactionsByYear(userId string, year int) ([]Transaction, error)
	GetTransaction(transactionId string) (Transaction, error)
	UpdateTransaction(transaction Transaction) error
	DeleteTransaction(transactionId string) error
	DeleteTransactionsInLocation(locationId string) error
	DeleteTransactionsInGroup(groupId string) error
	DeleteAllTransactionsForUser(userId string) error
}

type Transaction struct {
	Id              string  `json:"id"`
	UserId          string  `json:"userId"`
	Amount          float64 `json:"amount"`
	Timestamp       int64   `json:"timestamp"`
	Year            int     `json:"year"`
	Description     string  `json:"description,omitempty"`
	GroupId         string  `json:"groupId,omitempty"`
	MoneyLocationId string  `json:"moneyLocationId"`
}

func NewTransaction(userId string, amount float64, description, groupId, moneyLocationId string) (Transaction, error) {
	if userId == "" {
		return Transaction{}, errors.New("user id is required")
	}

	if amount == 0 {
		return Transaction{}, errors.New("amount is required")
	}

	if moneyLocationId == "" {
		return Transaction{}, errors.New("money location id is required")
	}

	return Transaction{
		UserId:          userId,
		Amount:          amount,
		Timestamp:       time.Now().UnixMilli(),
		Year:            time.Now().Year(),
		Description:     description,
		GroupId:         groupId,
		MoneyLocationId: moneyLocationId,
	}, nil
}
