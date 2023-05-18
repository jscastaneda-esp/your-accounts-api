package persistent

//go:generate mockery --name Transaction --filename transaction.go
type Transaction interface {
	Set(tx any) error
	Get() any
}

//go:generate mockery --name TransactionManager --filename transaction-manager.go
type TransactionManager interface {
	Transaction(fc func(tx Transaction) error) error
}
