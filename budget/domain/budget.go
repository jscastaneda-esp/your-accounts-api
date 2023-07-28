package domain

import (
	"context"
	"your-accounts-api/shared/domain/persistent"
)

type Budget struct {
	ID                    *uint
	Name                  *string
	Year                  *uint16
	Month                 *uint8
	FixedIncome           *float64
	AdditionalIncome      *float64
	TotalPendingPayment   *float64
	TotalAvailableBalance *float64
	PendingBills          *uint8
	TotalBalance          *float64
	ProjectId             *uint
}

//go:generate mockery --name BudgetRepository --filename budget-repository.go
type BudgetRepository interface {
	persistent.TransactionRepository[BudgetRepository]
	persistent.SaveRepository[Budget]
	persistent.SearchRepository[Budget]
	SearchByUserId(ctx context.Context, userId uint) ([]*Budget, error)
	persistent.DeleteRepository
}
