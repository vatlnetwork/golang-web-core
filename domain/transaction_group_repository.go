package domain

type TransactionGroupRepository interface {
	CreateTransactionGroup(transactionGroup TransactionGroup) (TransactionGroup, error)
	GetTransactionGroupsForUser(userId string) ([]TransactionGroup, error)
	GetTransactionGroup(transactionGroupId string) (TransactionGroup, error)
	UpdateTransactionGroup(transactionGroup TransactionGroup) error
	DeleteTransactionGroup(transactionGroupId string) error
}
