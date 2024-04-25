package persistent

type Transaction interface {
	Set(tx any) error
	Get() any
}

type TransactionManager interface {
	Transaction(fc func(tx Transaction) error) error
}
