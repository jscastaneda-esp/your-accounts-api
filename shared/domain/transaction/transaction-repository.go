package transaction

type TransactionRepository interface {
	WithTransaction(tx Transaction) any
}
