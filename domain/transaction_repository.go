package domain

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
