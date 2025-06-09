package domain

import "errors"

const ErrorTransactionGroupNotFound string = "transaction group not found"

type TransactionGroupRepository interface {
	CreateTransactionGroup(transactionGroup TransactionGroup) (TransactionGroup, error)
	GetTransactionGroupsForUser(userId string) ([]TransactionGroup, error)
	GetTransactionGroup(transactionGroupId string) (TransactionGroup, error)
	UpdateTransactionGroup(transactionGroup TransactionGroup) error
	DeleteTransactionGroup(transactionGroupId string) error
	DeleteAllTransactionGroupsForUser(userId string) error
}

type TransactionGroup struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	Description string `json:"description"`
}

func NewTransactionGroup(userId, description string) (TransactionGroup, error) {
	if userId == "" {
		return TransactionGroup{}, errors.New("user id is required")
	}

	if description == "" {
		description = "Unnamed Group"
	}

	return TransactionGroup{
		UserId:      userId,
		Description: description,
	}, nil
}
