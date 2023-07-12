package application

import (
	"context"
	"fmt"
	"time"
	budgetDom "your-accounts-api/budget/domain"
	"your-accounts-api/project/domain"
	"your-accounts-api/shared/domain/persistent"
	"your-accounts-api/shared/domain/utils/slices"
)

//go:generate mockery --name IProjectApp --filename project-app.go
type IProjectApp interface {
	Create(ctx context.Context, createData CreateData) (uint, error)
	Clone(ctx context.Context, baseId uint) (uint, error)
	FindByUser(ctx context.Context, userId uint) ([]*FindByUserRecord, error)
	FindLogsByProject(ctx context.Context, projectId uint) ([]*domain.ProjectLog, error)
	Delete(ctx context.Context, id uint) error
}

type projectApp struct {
	tm             persistent.TransactionManager
	projectRepo    domain.ProjectRepository
	projectLogRepo domain.ProjectLogRepository
	budgetRepo     budgetDom.BudgetRepository
}

func (app *projectApp) Create(ctx context.Context, createData CreateData) (uint, error) {
	var id uint
	err := app.tm.Transaction(func(tx persistent.Transaction) error {
		var err error
		id, err = app.createProject(ctx, createData.UserId, createData.Type, tx)
		if err != nil {
			return err
		}

		if createData.Type == domain.TypeBudget {
			budgetRepo := app.budgetRepo.WithTransaction(tx)
			now := time.Now()
			newBudget := budgetDom.Budget{
				Name:      createData.Name,
				Year:      uint16(now.Year()),
				Month:     uint8(now.Month()),
				ProjectId: id,
			}
			_, err = budgetRepo.Create(ctx, newBudget)
			if err != nil {
				return err
			}
		}

		err = app.createLog(ctx, "Creación", id, nil, tx)
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

func (app *projectApp) Clone(ctx context.Context, baseId uint) (uint, error) {
	baseProject, err := app.projectRepo.FindById(ctx, baseId)
	if err != nil {
		return 0, err
	}

	var id uint
	err = app.tm.Transaction(func(tx persistent.Transaction) error {
		var err error
		id, err = app.createProject(ctx, baseProject.UserId, baseProject.Type, tx)
		if err != nil {
			return err
		}

		var name string
		if baseProject.Type == domain.TypeBudget {
			baseBudgets, err := app.budgetRepo.FindByProjectIds(ctx, []uint{baseProject.ID})
			if err != nil {
				return err
			}

			if len(baseBudgets) > 0 {
				name = baseBudgets[0].Name
				budgetRepo := app.budgetRepo.WithTransaction(tx)
				newBudget := budgetDom.Budget{
					Name:             fmt.Sprintf("%s Copia", name),
					Year:             baseBudgets[0].Year,
					Month:            baseBudgets[0].Month,
					FixedIncome:      baseBudgets[0].FixedIncome,
					AdditionalIncome: baseBudgets[0].AdditionalIncome,
					Total:            baseBudgets[0].Total,
					EstimatedBalance: baseBudgets[0].EstimatedBalance,
					ProjectId:        id,
				}
				_, err = budgetRepo.Create(ctx, newBudget)
				if err != nil {
					return err
				}
			}
		}

		detail := fmt.Sprintf(`{"cloneId": %d, "cloneName": "%s"}`, baseId, name)
		err = app.createLog(ctx, fmt.Sprintf("Creación a partir de %s(%d)", name, baseId), id, &detail, tx)
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

func (app *projectApp) FindByUser(ctx context.Context, userId uint) ([]*FindByUserRecord, error) {
	projects, err := app.projectRepo.FindByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	response := []*FindByUserRecord{}
	projectBudgetIds := []uint{}
	for _, project := range projects {
		response = append(response, &FindByUserRecord{
			ID:   project.ID,
			Type: project.Type,
			Data: make(map[string]any),
		})

		if domain.TypeBudget == project.Type {
			projectBudgetIds = append(projectBudgetIds, project.ID)
		}
	}

	if len(projectBudgetIds) > 0 {
		budgets, err := app.budgetRepo.FindByProjectIds(ctx, projectBudgetIds)
		if err != nil {
			return nil, err
		}

		for _, budget := range budgets {
			if record := slices.Find(response, func(record *FindByUserRecord) bool {
				return record.ID == budget.ProjectId
			}); record != nil {
				record.Name = budget.Name
				record.Data["year"] = budget.Year
				record.Data["month"] = budget.Month
				record.Data["totalAvailableBalance"] = budget.TotalAvailableBalance
				record.Data["totalPendingPayment"] = budget.TotalPendingPayment
				record.Data["totalBalance"] = budget.TotalBalance
				record.Data["pendingBills"] = budget.PendingBills
			}
		}
	}

	return response, nil
}

func (app *projectApp) FindLogsByProject(ctx context.Context, projectId uint) ([]*domain.ProjectLog, error) {
	logs, err := app.projectLogRepo.FindByProjectId(ctx, projectId)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (app *projectApp) Delete(ctx context.Context, id uint) error {
	existsProject, err := app.projectRepo.FindById(ctx, id)
	if err != nil {
		return err
	}

	return app.tm.Transaction(func(tx persistent.Transaction) error {
		if domain.TypeBudget == existsProject.Type {
			budgetRepo := app.budgetRepo.WithTransaction(tx)
			if err := budgetRepo.DeleteByProjectId(ctx, existsProject.ID); err != nil {
				return err
			}
		}

		projectRepo := app.projectRepo.WithTransaction(tx)
		return projectRepo.Delete(ctx, id)
	})
}

func (app *projectApp) createProject(ctx context.Context, userId uint, typeProject domain.ProjectType, tx persistent.Transaction) (uint, error) {
	projectRepo := app.projectRepo.WithTransaction(tx)
	newProject := domain.Project{
		UserId: userId,
		Type:   typeProject,
	}
	id, err := projectRepo.Create(ctx, newProject)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (app *projectApp) createLog(ctx context.Context, description string, projectId uint, detail *string, tx persistent.Transaction) error {
	projectLogRepo := app.projectLogRepo.WithTransaction(tx)
	newLog := domain.ProjectLog{
		Description: description,
		ProjectId:   projectId,
		Detail:      detail,
	}
	_, err := projectLogRepo.Create(ctx, newLog)
	if err != nil {
		return err
	}

	return nil
}

func NewProjectApp(
	tm persistent.TransactionManager, projectRepo domain.ProjectRepository,
	projectLogRepo domain.ProjectLogRepository, budgetRepo budgetDom.BudgetRepository,
) IProjectApp {
	return &projectApp{tm, projectRepo, projectLogRepo, budgetRepo}
}
