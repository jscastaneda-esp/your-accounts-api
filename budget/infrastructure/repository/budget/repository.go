package budget

import (
	"context"
	"your-accounts-api/budget/domain"
	"your-accounts-api/budget/infrastructure/entity"
	"your-accounts-api/shared/domain/persistent"
	shared_ent "your-accounts-api/shared/infrastructure/db/entity"
	persistent_infra "your-accounts-api/shared/infrastructure/db/persistent"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormRepository struct {
	db *gorm.DB
}

func (r *gormRepository) WithTransaction(tx persistent.Transaction) domain.BudgetRepository {
	return persistent_infra.DefaultWithTransaction[domain.BudgetRepository](tx, NewRepository, r)
}

func (r *gormRepository) Create(ctx context.Context, budget domain.Budget) (uint, error) {
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
		return 0, err
	}

	return model.ID, nil
}

func (r *gormRepository) FindById(ctx context.Context, id uint) (*domain.Budget, error) {
	model := new(entity.Budget)
	if err := r.db.WithContext(ctx).First(model, id).Error; err != nil {
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

func (r *gormRepository) FindByProjectIds(ctx context.Context, projectIds []uint) ([]*domain.Budget, error) {
	var models []entity.Budget
	if err := r.db.WithContext(ctx).Where("project_id IN ?", projectIds).Find(&models).Error; err != nil {
		return nil, err
	}

	var budgets []*domain.Budget
	for _, model := range models {
		budgets = append(budgets, &domain.Budget{
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
		})
	}

	return budgets, nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Select(clause.Associations).Delete(&entity.Budget{
		BaseModel: shared_ent.BaseModel{
			ID: id,
		},
	}).Error; err != nil {
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB) domain.BudgetRepository {
	return &gormRepository{db}
}
