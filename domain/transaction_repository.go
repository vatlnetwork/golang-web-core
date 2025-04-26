package domain

type TransactionRepository interface {
	CreateTransaction(transaction Transaction) (Transaction, error)
	GetTransactionsForUser(userId string) ([]Transaction, error)
	GetTransaction(transactionId string) (Transaction, error)
	UpdateTransaction(transaction Transaction) error
	DeleteTransaction(transactionId string) error
}
