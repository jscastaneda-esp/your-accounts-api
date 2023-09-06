package application

import (
	"context"
	"fmt"
	"time"
	"your-accounts-api/budgets/domain"
	"your-accounts-api/shared/application"
	shared "your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/persistent"
)

//go:generate mockery --name IBudgetApp --filename budget-app.go
type IBudgetApp interface {
	Create(ctx context.Context, userId uint, name string) (uint, error)
	Clone(ctx context.Context, userId uint, baseId uint) (uint, error)
	FindById(ctx context.Context, id uint) (*domain.Budget, error)
	FindByUserId(ctx context.Context, userId uint) ([]*domain.Budget, error)
	Delete(ctx context.Context, id uint) error
}

type budgetApp struct {
	tm         persistent.TransactionManager
	budgetRepo domain.BudgetRepository
	logApp     application.ILogApp
}

func (app *budgetApp) Create(ctx context.Context, userId uint, name string) (uint, error) {
	var id uint
	err := app.tm.Transaction(func(tx persistent.Transaction) error {
		budgetRepo := app.budgetRepo.WithTransaction(tx)
		now := time.Now()
		year := uint16(now.Year())
		month := uint8(now.Month())
		newBudget := domain.Budget{
			Name:   &name,
			Year:   &year,
			Month:  &month,
			UserId: &userId,
		}

		var err error
		id, err = budgetRepo.Save(ctx, newBudget)
		if err != nil {
			return err
		}

		return app.logApp.CreateLog(ctx, "Creación", shared.Budget, id, nil, tx)
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (app *budgetApp) Clone(ctx context.Context, userId uint, baseId uint) (uint, error) {
	baseBudget, err := app.FindById(ctx, baseId)
	if err != nil {
		return 0, err
	}

	var id uint
	err = app.tm.Transaction(func(tx persistent.Transaction) error {
		budgetRepo := app.budgetRepo.WithTransaction(tx)
		name := fmt.Sprintf("%s Copia", *baseBudget.Name)
		newBudget := domain.Budget{
			Name:             &name,
			Year:             baseBudget.Year,
			Month:            baseBudget.Month,
			FixedIncome:      baseBudget.FixedIncome,
			AdditionalIncome: baseBudget.AdditionalIncome,
			UserId:           &userId,
		}

		var err error
		id, err = budgetRepo.Save(ctx, newBudget)
		if err != nil {
			return err
		}

		description := fmt.Sprintf("Creación a partir del presupuesto %s(%d)", *baseBudget.Name, baseId)
		detail := fmt.Sprintf(`{"cloneId": %d, "cloneName": "%s"}`, baseId, *baseBudget.Name)
		err = app.logApp.CreateLog(ctx, description, shared.Budget, id, &detail, tx)
		if err != nil {
			return err
		}

		// TODO Pendiente la creación de AvailableBalances, Bills y BillShared

		return nil
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (app *budgetApp) FindById(ctx context.Context, id uint) (*domain.Budget, error) {
	budget, err := app.budgetRepo.Search(ctx, id)
	if err != nil {
		return nil, err
	}

	return budget, nil
}

func (app *budgetApp) FindByUserId(ctx context.Context, userId uint) ([]*domain.Budget, error) {
	example := domain.Budget{
		UserId: &userId,
	}
	budgets, err := app.budgetRepo.SearchAllByExample(ctx, example)
	if err != nil {
		return nil, err
	}

	return budgets, nil
}

func (app *budgetApp) Delete(ctx context.Context, id uint) error {
	budget, err := app.FindById(ctx, id)
	if err != nil {
		return err
	}

	return app.tm.Transaction(func(tx persistent.Transaction) error {
		return app.budgetRepo.Delete(ctx, *budget.ID) // TODO Validar si se puede eliminar toda la información
	})
}

func NewBudgetApp(tm persistent.TransactionManager, budgetRepo domain.BudgetRepository, logApp application.ILogApp) IBudgetApp {
	return &budgetApp{tm, budgetRepo, logApp}
}
