package application

import (
	"context"
	"fmt"
	"your-accounts-api/budgets/domain"
	"your-accounts-api/shared/application"
	shared "your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/persistent"
)

//go:generate mockery --name IBudgetAvailableApp --filename budget-available-app.go
type IBudgetAvailableApp interface {
	Create(ctx context.Context, name string, budgetId uint) (uint, error)
}

type budgetAvailableApp struct {
	tm                  persistent.TransactionManager
	budgetAvailableRepo domain.BudgetAvailableRepository
	logApp              application.ILogApp
}

func (app *budgetAvailableApp) Create(ctx context.Context, name string, budgetId uint) (uint, error) {
	var id uint
	err := app.tm.Transaction(func(tx persistent.Transaction) error {
		var err error
		description := fmt.Sprintf("Se crea el disponible %s", name)
		err = app.logApp.CreateLog(ctx, description, shared.Budget, budgetId, nil, tx)
		if err != nil {
			return err
		}

		newAvailable := domain.BudgetAvailable{
			Name:     &name,
			BudgetId: &budgetId,
		}
		budgetAvailableRepo := app.budgetAvailableRepo.WithTransaction(tx)
		id, err = budgetAvailableRepo.Save(ctx, newAvailable)
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

func NewBudgetAvailableApp(tm persistent.TransactionManager, budgetAvailableRepo domain.BudgetAvailableRepository, logApp application.ILogApp) IBudgetAvailableApp {
	return &budgetAvailableApp{tm, budgetAvailableRepo, logApp}
}
