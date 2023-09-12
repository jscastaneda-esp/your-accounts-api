package domain

type BudgetBill struct {
	ID          *uint
	Description *string
	Amount      *float64
	Payment     *float64
	DueDate     *uint8
	Complete    *bool
	BudgetId    *uint
	Category    *BudgetBillCategory
}

//go:generate mockery --name BudgetBillRepository --filename budget-bill-repository.go
type BudgetBillRepository interface {
}
