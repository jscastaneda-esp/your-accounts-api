package persistent

type TransactionRepository[T any] interface {
	WithTransaction(tx Transaction) T
}
