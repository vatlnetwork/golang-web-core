package domain

type TransactionGroup struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	Description string `json:"description"`
}

func NewTransactionGroup(userId, description string) TransactionGroup {
	return TransactionGroup{
		UserId:      userId,
		Description: description,
	}
}
