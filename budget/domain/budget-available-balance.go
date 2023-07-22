package domain

import (
	"time"
	"your-accounts-api/shared/domain/persistent"
)

type BudgetAvailableBalance struct {
	ID        uint
	Name      string
	Amount    float64
	BudgetId  uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

//go:generate mockery --name BudgetAvailableBalanceRepository --filename budget-available-balance-repository.go
type BudgetAvailableBalanceRepository interface {
	persistent.TransactionRepository[BudgetAvailableBalanceRepository]
	persistent.CreateRepository[BudgetAvailableBalance]
}
