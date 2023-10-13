package domain

import (
	"your-accounts-api/shared/domain/persistent"
)

type Budget struct {
	ID               *uint
	Name             *string
	Year             *uint16
	Month            *uint8
	FixedIncome      *float64
	AdditionalIncome *float64
	TotalPending     *float64
	TotalAvailable   *float64
	TotalSaving      *float64
	PendingBills     *uint8
	UserId           *uint
	BudgetAvailables []BudgetAvailable
	BudgetBills      []BudgetBill
}

//go:generate mockery --name BudgetRepository --filename budget-repository.go
type BudgetRepository interface {
	persistent.TransactionRepository[BudgetRepository]
	persistent.SaveRepository[Budget]
	persistent.SearchRepository[Budget]
	persistent.SearchAllByExampleRepository[Budget]
	persistent.DeleteRepository
}
