package budget

import (
	"api-your-accounts/budget/domain"
	"api-your-accounts/budget/infrastructure/entity"
	"api-your-accounts/shared/domain/persistent"
	persistentInfra "api-your-accounts/shared/infrastructure/db/persistent"
	"context"

	"gorm.io/gorm"
)

type gormBudgetRepository struct {
	db *gorm.DB
}

func (r *gormBudgetRepository) WithTransaction(tx persistent.Transaction) domain.BudgetRepository {
	return persistentInfra.DefaultWithTransaction[domain.BudgetRepository](tx, NewRepository, r)
}

func (r *gormBudgetRepository) Create(ctx context.Context, budget *domain.Budget) (*domain.Budget, error) {
	model := &entity.Budget{
		Name:                  budget.Name,
		Year:                  budget.Year,
		Month:                 budget.Month,
		FixedIncome:           budget.FixedIncome,
		AdditionalIncome:      budget.AdditionalIncome,
		TotalPendingPayment:   budget.TotalPendingPayment,
		TotalAvailableBalance: budget.TotalAvailableBalance,
		PendingBills:          budget.PendingBills,
		TotalBalance:          budget.TotalBalance,
		Total:                 budget.Total,
		EstimatedBalance:      budget.EstimatedBalance,
		TotalPayment:          budget.TotalPayment,
		ProjectId:             budget.ProjectId,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	return &domain.Budget{
		ID:                    model.ID,
		Name:                  model.Name,
		Year:                  model.Year,
		Month:                 model.Month,
		FixedIncome:           model.FixedIncome,
		AdditionalIncome:      model.AdditionalIncome,
		TotalPendingPayment:   model.TotalPendingPayment,
		TotalAvailableBalance: model.TotalAvailableBalance,
		PendingBills:          model.PendingBills,
		TotalBalance:          model.TotalBalance,
		Total:                 model.Total,
		EstimatedBalance:      model.EstimatedBalance,
		TotalPayment:          model.TotalPayment,
		ProjectId:             model.ProjectId,
		CreatedAt:             model.CreatedAt,
		UpdatedAt:             model.UpdatedAt,
	}, nil
}

func (r *gormBudgetRepository) FindById(ctx context.Context, id uint) (*domain.Budget, error) {
	model, err := r.findById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Budget{
		ID:                    model.ID,
		Name:                  model.Name,
		Year:                  model.Year,
		Month:                 model.Month,
		FixedIncome:           model.FixedIncome,
		AdditionalIncome:      model.AdditionalIncome,
		TotalPendingPayment:   model.TotalPendingPayment,
		TotalAvailableBalance: model.TotalAvailableBalance,
		PendingBills:          model.PendingBills,
		TotalBalance:          model.TotalBalance,
		Total:                 model.Total,
		EstimatedBalance:      model.EstimatedBalance,
		TotalPayment:          model.TotalPayment,
		ProjectId:             model.ProjectId,
		CreatedAt:             model.CreatedAt,
		UpdatedAt:             model.UpdatedAt,
	}, nil
}

func (r *gormBudgetRepository) FindByProjectId(ctx context.Context, projectId uint) (*domain.Budget, error) {
	where := &entity.Budget{
		ProjectId: projectId,
	}
	model := new(entity.Budget)
	if err := r.db.WithContext(ctx).Where(where).First(model).Error; err != nil {
		return nil, err
	}

	budget := &domain.Budget{
		ID:                    model.ID,
		Name:                  model.Name,
		Year:                  model.Year,
		Month:                 model.Month,
		FixedIncome:           model.FixedIncome,
		AdditionalIncome:      model.AdditionalIncome,
		TotalPendingPayment:   model.TotalPendingPayment,
		TotalAvailableBalance: model.TotalAvailableBalance,
		PendingBills:          model.PendingBills,
		TotalBalance:          model.TotalBalance,
		Total:                 model.Total,
		EstimatedBalance:      model.EstimatedBalance,
		TotalPayment:          model.TotalPayment,
		ProjectId:             model.ProjectId,
		CreatedAt:             model.CreatedAt,
		UpdatedAt:             model.UpdatedAt,
	}

	return budget, nil
}

func (r *gormBudgetRepository) Delete(ctx context.Context, id uint) error {
	model, err := r.findById(ctx, id)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Delete(model).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormBudgetRepository) findById(ctx context.Context, id uint) (*entity.Budget, error) {
	model := new(entity.Budget)
	if err := r.db.WithContext(ctx).First(model, id).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func NewRepository(db *gorm.DB) domain.BudgetRepository {
	return &gormBudgetRepository{db}
}