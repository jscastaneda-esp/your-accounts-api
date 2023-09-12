package domain

import (
	"your-accounts-api/shared/domain/persistent"
)

type BudgetAvailable struct {
	ID       *uint
	Name     *string
	Amount   *float64
	BudgetId *uint
}

//go:generate mockery --name BudgetAvailableRepository --filename budget-available-repository.go
type BudgetAvailableRepository interface {
	persistent.TransactionRepository[BudgetAvailableRepository]
	persistent.SaveRepository[BudgetAvailable]
	// persistent.SaveAllRepository[BudgetAvailable]
	// persistent.SearchAllByExample[BudgetAvailable]
	// persistent.DeleteRepository
}
