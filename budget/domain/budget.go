package domain

import (
	"api-your-accounts/shared/domain/persistent"
	"context"
	"time"
)

type Budget struct {
	ID                    uint
	Name                  string
	Year                  uint16
	Month                 uint8
	FixedIncome           float64
	AdditionalIncome      float64
	TotalPendingPayment   float64
	TotalAvailableBalance float64
	PendingBills          uint8
	TotalBalance          float64
	Total                 float64
	EstimatedBalance      float64
	TotalPayment          float64
	ProjectId             uint
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

//go:generate mockery --name BudgetRepository --filename budget-repository.go
type BudgetRepository interface {
	persistent.TransactionRepository[BudgetRepository]
	persistent.CreateRepository[Budget]
	persistent.ReadRepository[Budget, uint]
	FindByProjectIds(ctx context.Context, projectIds []uint) ([]*Budget, error)
	DeleteByProjectId(ctx context.Context, projectId uint) error
}
