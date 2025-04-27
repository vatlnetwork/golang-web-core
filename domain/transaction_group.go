package domain

import "errors"

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
