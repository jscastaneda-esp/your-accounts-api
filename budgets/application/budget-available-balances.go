package application

import (
	"context"
	"fmt"
	"your-accounts-api/budgets/domain"
	"your-accounts-api/shared/application"
	shared "your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/persistent"
)

//go:generate mockery --name IBudgetAvailableBalanceApp --filename budget-available-balance-app.go
type IBudgetAvailableBalanceApp interface {
	Create(ctx context.Context, name string, budgetId uint) (uint, error)
}

type budgetAvailableBalanceApp struct {
	tm                         persistent.TransactionManager
	budgetAvailableBalanceRepo domain.BudgetAvailableBalanceRepository
	budgetApp                  IBudgetApp
	logApp                     application.ILogApp
}

func (app *budgetAvailableBalanceApp) Create(ctx context.Context, name string, budgetId uint) (uint, error) {
	budget, err := app.budgetApp.FindById(ctx, budgetId)
	if err != nil {
		return 0, err
	}

	var id uint
	err = app.tm.Transaction(func(tx persistent.Transaction) error {
		var err error
		description := fmt.Sprintf("Se crea el disponible %s", name)
		err = app.logApp.CreateLog(ctx, description, shared.Budget, *budget.ID, nil, tx)
		if err != nil {
			return err
		}

		newAvailable := domain.BudgetAvailableBalance{
			Name:     &name,
			BudgetId: &budgetId,
		}
		budgetAvailableBalanceRepo := app.budgetAvailableBalanceRepo.WithTransaction(tx)
		id, err = budgetAvailableBalanceRepo.Save(ctx, newAvailable)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func NewBudgetAvailableBalanceApp(tm persistent.TransactionManager, budgetAvailableBalanceRepo domain.BudgetAvailableBalanceRepository, budgetApp IBudgetApp, logApp application.ILogApp) IBudgetAvailableBalanceApp {
	return &budgetAvailableBalanceApp{tm, budgetAvailableBalanceRepo, budgetApp, logApp}
}
