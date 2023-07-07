package application

import (
	"context"
	"your-accounts-api/budget/domain"
	"your-accounts-api/shared/domain/persistent"
)

//go:generate mockery --name IBudgetApp --filename budget-app.go
type IBudgetApp interface {
	FindById(ctx context.Context, id uint) (*domain.Budget, error)
}

type budgetApp struct {
	tm         persistent.TransactionManager
	budgetRepo domain.BudgetRepository
}

func (app *budgetApp) FindById(ctx context.Context, id uint) (*domain.Budget, error) {
	budget, err := app.budgetRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return budget, nil
}

func NewBudgetApp(tm persistent.TransactionManager, budgetRepo domain.BudgetRepository) IBudgetApp {
	return &budgetApp{tm, budgetRepo}
}
