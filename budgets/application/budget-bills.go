package application

import (
	"context"
	"fmt"
	"your-accounts-api/budgets/domain"
	"your-accounts-api/shared/application"
	shared "your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/persistent"
)

//go:generate mockery --name IBudgetBillApp --filename budget-bill-app.go
type IBudgetBillApp interface {
	Create(ctx context.Context, description string, category domain.BudgetBillCategory, budgetId uint) (uint, error)
}

type budgetBillApp struct {
	tm             persistent.TransactionManager
	budgetBillRepo domain.BudgetBillRepository
	logApp         application.ILogApp
}

func (app *budgetBillApp) Create(ctx context.Context, description string, category domain.BudgetBillCategory, budgetId uint) (uint, error) {
	var id uint
	err := app.tm.Transaction(func(tx persistent.Transaction) error {
		var err error
		descriptionLog := fmt.Sprintf("Se crea el pago %s", description)
		err = app.logApp.CreateLog(ctx, descriptionLog, shared.Budget, budgetId, nil, tx)
		if err != nil {
			return err
		}

		newBill := domain.BudgetBill{
			Description: &description,
			BudgetId:    &budgetId,
			Category:    &category,
		}
		budgetBillRepo := app.budgetBillRepo.WithTransaction(tx)
		id, err = budgetBillRepo.Save(ctx, newBill)
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

func NewBudgetBillApp(tm persistent.TransactionManager, budgetBillRepo domain.BudgetBillRepository, logApp application.ILogApp) IBudgetBillApp {
	return &budgetBillApp{tm, budgetBillRepo, logApp}
}
