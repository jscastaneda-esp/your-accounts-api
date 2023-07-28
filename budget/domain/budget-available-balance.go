package domain

import (
	"your-accounts-api/shared/domain/persistent"
)

type BudgetAvailableBalance struct {
	ID       *uint
	Name     *string
	Amount   *float64
	BudgetId *uint
}

//go:generate mockery --name BudgetAvailableBalanceRepository --filename budget-available-balance-repository.go
type BudgetAvailableBalanceRepository interface {
	persistent.TransactionRepository[BudgetAvailableBalanceRepository]
	persistent.SaveRepository[BudgetAvailableBalance]
	persistent.SaveAllRepository[BudgetAvailableBalance]
	persistent.SearchAllByExample[BudgetAvailableBalance]
	persistent.DeleteRepository
}
