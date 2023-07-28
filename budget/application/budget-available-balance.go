package application

import (
	"context"
	"fmt"
	"your-accounts-api/budget/domain"
	"your-accounts-api/project/application"
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
	projectApp                 application.IProjectApp
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
		err = app.projectApp.CreateLog(ctx, description, *budget.ProjectId, nil, tx)
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

func NewBudgetAvailableBalanceApp(tm persistent.TransactionManager, budgetAvailableBalanceRepo domain.BudgetAvailableBalanceRepository, budgetApp IBudgetApp, projectApp application.IProjectApp) IBudgetAvailableBalanceApp {
	return &budgetAvailableBalanceApp{tm, budgetAvailableBalanceRepo, budgetApp, projectApp}
}
