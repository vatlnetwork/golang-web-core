package domain

import (
	"errors"
	"time"
)

type Transaction struct {
	Id          string  `json:"id"`
	UserId      string  `json:"userId"`
	Amount      float64 `json:"amount"`
	Timestamp   int64   `json:"timestamp"`
	Year        int     `json:"year"`
	Description string  `json:"description,omitempty"`
	GroupId     string  `json:"groupId,omitempty"`
}

func NewTransaction(userId string, amount float64, description, groupId string) (Transaction, error) {
	if userId == "" {
		return Transaction{}, errors.New("user id is required")
	}

	if amount == 0 {
		return Transaction{}, errors.New("amount is required")
	}

	return Transaction{
		UserId:      userId,
		Amount:      amount,
		Timestamp:   time.Now().UnixMilli(),
		Year:        time.Now().Year(),
		Description: description,
		GroupId:     groupId,
	}, nil
}
