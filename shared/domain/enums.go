package domain

type CodeLog string

const (
	Budget     CodeLog = "budget"
	BudgetBill CodeLog = "budget_bill"
)

type Action string

const (
	Update Action = "update"
	Delete Action = "delete"
)
