package application

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"your-accounts-api/budgets/domain"
	"your-accounts-api/shared/application"
	shared "your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/persistent"
	"your-accounts-api/shared/domain/utils/convert"
)

var ErrIncompleteData = errors.New("incomplete data")

type Change struct {
	ID      uint
	Section domain.BudgetSection
	Action  shared.Action
	Detail  map[string]any
}

type ChangeResult struct {
	Change Change
	Err    error
}

//go:generate mockery --name IBudgetApp --filename budget-app.go
type IBudgetApp interface {
	Create(ctx context.Context, userId uint, name string) (uint, error)
	Clone(ctx context.Context, userId uint, baseId uint) (uint, error)
	FindById(ctx context.Context, id uint) (*domain.Budget, error)
	FindByUserId(ctx context.Context, userId uint) ([]domain.Budget, error)
	Changes(ctx context.Context, id uint, changes []Change) []ChangeResult
	Delete(ctx context.Context, id uint) error
}

type budgetApp struct {
	tm                  persistent.TransactionManager
	budgetRepo          domain.BudgetRepository
	budgetAvailableRepo domain.BudgetAvailableRepository
	budgetBillRepo      domain.BudgetBillRepository
	logApp              application.ILogApp
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

		budgetAvailableRepo := app.budgetAvailableRepo.WithTransaction(tx)
		newAvailables := []domain.BudgetAvailable{}
		for _, available := range baseBudget.BudgetAvailables {
			newAvailable := domain.BudgetAvailable{
				Name:     available.Name,
				Amount:   available.Amount,
				BudgetId: &id,
			}
			newAvailables = append(newAvailables, newAvailable)
		}
		err = budgetAvailableRepo.SaveAll(ctx, newAvailables)
		if err != nil {
			return err
		}

		budgetBillRepo := app.budgetBillRepo.WithTransaction(tx)
		newBills := []domain.BudgetBill{}
		for _, bill := range baseBudget.BudgetBills {
			newBill := domain.BudgetBill{
				Description: bill.Description,
				Amount:      bill.Amount,
				DueDate:     bill.DueDate,
				BudgetId:    &id,
				Category:    bill.Category,
			}
			newBills = append(newBills, newBill)
		}
		err = budgetBillRepo.SaveAll(ctx, newBills)
		if err != nil {
			return err
		}

		description := fmt.Sprintf("Creación a partir del presupuesto %s(%d)", *baseBudget.Name, baseId)
		detail := map[string]any{
			"cloneId":   baseId,
			"cloneName": *baseBudget.Name,
		}
		return app.logApp.CreateLog(ctx, description, shared.Budget, id, detail, tx)
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

func (app *budgetApp) FindByUserId(ctx context.Context, userId uint) ([]domain.Budget, error) {
	example := domain.Budget{
		UserId: &userId,
	}
	budgets, err := app.budgetRepo.SearchAllByExample(ctx, example)
	if err != nil {
		return nil, err
	}

	return budgets, nil
}

func (app *budgetApp) Changes(ctx context.Context, id uint, changes []Change) []ChangeResult {
	changesChan := make(chan Change, len(changes))
	resultsChan := make(chan ChangeResult)
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go app.changeWorker(ctx, id, changesChan, resultsChan, &wg)
	}

	for _, change := range changes {
		changesChan <- change
	}
	close(changesChan)

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	changeResults := []ChangeResult{}
	resultsWithError := 0
	for result := range resultsChan {
		changeResults = append(changeResults, result)

		if result.Err != nil {
			resultsWithError += 1
		}
	}

	return app.updateTotals(ctx, id, changeResults, resultsWithError)
}

func (app *budgetApp) Delete(ctx context.Context, id uint) error {
	budget, err := app.FindById(ctx, id)
	if err != nil {
		return err
	}

	return app.budgetRepo.Delete(ctx, *budget.ID)
}

func (app *budgetApp) changeWorker(ctx context.Context, budgetId uint, changes <-chan Change, results chan<- ChangeResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for change := range changes {
		var changeResult ChangeResult
		switch change.Section {
		case domain.Main:
			changeResult = app.changeMain(ctx, budgetId, change)
		case domain.Available:
			changeResult = app.changeAvailable(ctx, budgetId, change)
		case domain.Bill:
			changeResult = app.changeBill(ctx, budgetId, change)
		}

		results <- changeResult
	}
}

func (app *budgetApp) changeMain(ctx context.Context, budgetId uint, change Change) ChangeResult {
	var err error

	switch change.Action {
	case shared.Update:
		if change.Detail == nil || len(change.Detail) == 0 {
			err = ErrIncompleteData
		} else {
			err = app.tm.Transaction(func(tx persistent.Transaction) error {
				budget := domain.Budget{
					ID: &change.ID,
				}

				if change.Detail["name"] != nil {
					name := fmt.Sprint(change.Detail["name"])
					budget.Name = &name
				}

				if change.Detail["year"] != nil {
					value, err := convert.AsUint64(change.Detail["year"])
					if err != nil {
						return err
					}

					year := uint16(value)
					budget.Year = &year
				}

				if change.Detail["month"] != nil {
					value, err := convert.AsUint64(change.Detail["month"])
					if err != nil {
						return err
					}

					month := uint8(value)
					budget.Month = &month
				}

				if change.Detail["fixedIncome"] != nil {
					fixedIncome, err := convert.AsFloat64(change.Detail["fixedIncome"])
					if err != nil {
						return err
					}

					budget.FixedIncome = &fixedIncome
				}

				if change.Detail["additionalIncome"] != nil {
					additionalIncome, err := convert.AsFloat64(change.Detail["additionalIncome"])
					if err != nil {
						return err
					}

					budget.AdditionalIncome = &additionalIncome
				}

				budgetRepo := app.budgetRepo.WithTransaction(tx)
				_, err := budgetRepo.Save(ctx, budget)
				if err != nil {
					return err
				}

				description := "Se actualiza la información principal"
				return app.logApp.CreateLog(ctx, description, shared.Budget, budgetId, change.Detail, tx)
			})
		}
	}

	return ChangeResult{change, err}
}

func (app *budgetApp) changeAvailable(ctx context.Context, budgetId uint, change Change) ChangeResult {
	var err error

	switch change.Action {
	case shared.Update:
		if change.Detail == nil || len(change.Detail) == 0 {
			err = ErrIncompleteData
		} else {
			err = app.tm.Transaction(func(tx persistent.Transaction) error {
				available := domain.BudgetAvailable{
					ID: &change.ID,
				}

				if change.Detail["name"] != nil {
					name := fmt.Sprint(change.Detail["name"])
					available.Name = &name
				}

				if change.Detail["amount"] != nil {
					amount, err := convert.AsFloat64(change.Detail["amount"])
					if err != nil {
						return err
					}

					available.Amount = &amount
				}

				budgetAvailableRepo := app.budgetAvailableRepo.WithTransaction(tx)
				_, err := budgetAvailableRepo.Save(ctx, available)
				if err != nil {
					return err
				}

				description := "Se actualizaron los disponibles"
				change.Detail["availableId"] = change.ID
				return app.logApp.CreateLog(ctx, description, shared.Budget, budgetId, change.Detail, tx)
			})
		}
	case shared.Delete:
		if change.Detail == nil || change.Detail["name"] == nil {
			err = ErrIncompleteData
		} else {
			err = app.tm.Transaction(func(tx persistent.Transaction) error {
				budgetAvailableRepo := app.budgetAvailableRepo.WithTransaction(tx)
				if err := budgetAvailableRepo.Delete(ctx, change.ID); err != nil {
					return err
				}

				description := fmt.Sprintf("Se elimino el disponible %s", change.Detail["name"])
				change.Detail["availableId"] = change.ID
				return app.logApp.CreateLog(ctx, description, shared.Budget, budgetId, change.Detail, tx)
			})
		}
	}

	return ChangeResult{change, err}
}

func (app *budgetApp) changeBill(ctx context.Context, budgetId uint, change Change) ChangeResult {
	var err error

	switch change.Action {
	case shared.Update:
		if change.Detail == nil || len(change.Detail) == 0 {
			err = ErrIncompleteData
		} else {
			err = app.tm.Transaction(func(tx persistent.Transaction) error {
				budget := domain.BudgetBill{
					ID: &change.ID,
				}

				if change.Detail["description"] != nil {
					description := fmt.Sprint(change.Detail["description"])
					budget.Description = &description
				}

				if change.Detail["amount"] != nil {
					amount, err := convert.AsFloat64(change.Detail["amount"])
					if err != nil {
						return err
					}

					budget.Amount = &amount
				}

				if change.Detail["dueDate"] != nil {
					value, err := convert.AsUint64(change.Detail["dueDate"])
					if err != nil {
						return err
					}

					dueDate := uint8(value)
					budget.DueDate = &dueDate
				}

				if change.Detail["complete"] != nil {
					complete, err := convert.AsBool(change.Detail["complete"])
					if err != nil {
						return err
					}

					budget.Complete = &complete
				}

				if change.Detail["category"] != nil {
					category := domain.BudgetBillCategory(fmt.Sprint(change.Detail["category"]))
					budget.Category = &category
				}

				budgetBillRepo := app.budgetBillRepo.WithTransaction(tx)
				_, err := budgetBillRepo.Save(ctx, budget)
				if err != nil {
					return err
				}

				description := "Se actualizaron los pagos"
				change.Detail["billId"] = change.ID
				return app.logApp.CreateLog(ctx, description, shared.Budget, budgetId, change.Detail, tx)
			})
		}
	case shared.Delete:
		if change.Detail == nil || change.Detail["description"] == nil {
			err = ErrIncompleteData
		} else {
			err = app.tm.Transaction(func(tx persistent.Transaction) error {
				var err error
				budgetBillRepo := app.budgetBillRepo.WithTransaction(tx)
				err = budgetBillRepo.Delete(ctx, change.ID)
				if err != nil {
					return err
				}

				description := fmt.Sprintf("Se elimino el pago %s", change.Detail["description"])
				change.Detail["billId"] = change.ID
				return app.logApp.CreateLog(ctx, description, shared.Budget, budgetId, change.Detail, tx)
			})
		}
	}

	return ChangeResult{change, err}
}

func (app *budgetApp) updateTotals(ctx context.Context, id uint, changeResults []ChangeResult, resultsWithError int) []ChangeResult {
	if len(changeResults) > resultsWithError {
		budget, err := app.budgetRepo.Search(ctx, id)
		if err != nil {
			changeResults = append(changeResults, ChangeResult{
				Err: err,
			})
		} else {
			totalAvailable := 0.0
			for _, available := range budget.BudgetAvailables {
				totalAvailable += *available.Amount
			}
			budget.TotalAvailable = &totalAvailable

			totalPending := 0.0
			pendingBills := uint8(0)
			for _, bill := range budget.BudgetBills {
				pending := *bill.Amount - *bill.Payment
				if !*bill.Complete {
					if pending < 0 {
						pending = 0
					}
					totalPending += pending
				}

				if *bill.Amount > *bill.Payment && !*bill.Complete {
					pendingBills += 1
				}
			}
			budget.TotalPending = &totalPending
			budget.PendingBills = &pendingBills

			_, err := app.budgetRepo.Save(ctx, *budget)
			if err != nil {
				changeResults = append(changeResults, ChangeResult{
					Err: err,
				})
			}
		}
	}

	return changeResults
}

func NewBudgetApp(
	tm persistent.TransactionManager, budgetRepo domain.BudgetRepository, budgetAvailableRepo domain.BudgetAvailableRepository,
	budgetBillRepo domain.BudgetBillRepository, logApp application.ILogApp,
) IBudgetApp {
	return &budgetApp{tm, budgetRepo, budgetAvailableRepo, budgetBillRepo, logApp}
}
