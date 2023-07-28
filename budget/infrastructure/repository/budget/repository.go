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

func (r *gormRepository) Save(ctx context.Context, budget domain.Budget) (uint, error) {
	model := new(entity.Budget)
	if budget.ID != nil {
		if err := r.db.WithContext(ctx).First(model, *budget.ID).Error; err != nil {
			return 0, err
		}
	} else {
		model.ProjectId = *budget.ProjectId
	}

	if budget.Name != nil {
		model.Name = *budget.Name
	}

	if budget.Year != nil {
		model.Year = *budget.Year
	}

	if budget.Month != nil {
		model.Month = *budget.Month
	}

	if budget.FixedIncome != nil {
		model.FixedIncome = *budget.FixedIncome
	}

	if budget.AdditionalIncome != nil {
		model.AdditionalIncome = *budget.AdditionalIncome
	}

	if budget.TotalPendingPayment != nil {
		model.TotalPendingPayment = *budget.TotalPendingPayment
	}

	if budget.TotalAvailableBalance != nil {
		model.TotalAvailableBalance = *budget.TotalAvailableBalance
	}

	if budget.PendingBills != nil {
		model.PendingBills = *budget.PendingBills
	}

	if budget.TotalBalance != nil {
		model.TotalBalance = *budget.TotalBalance
	}

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

func (r *gormRepository) Search(ctx context.Context, id uint) (*domain.Budget, error) {
	model := new(entity.Budget)
	if err := r.db.WithContext(ctx).First(model, id).Error; err != nil {
		return nil, err
	}

	return &domain.Budget{
		ID:                    &model.ID,
		Name:                  &model.Name,
		Year:                  &model.Year,
		Month:                 &model.Month,
		FixedIncome:           &model.FixedIncome,
		AdditionalIncome:      &model.AdditionalIncome,
		TotalPendingPayment:   &model.TotalPendingPayment,
		TotalAvailableBalance: &model.TotalAvailableBalance,
		PendingBills:          &model.PendingBills,
		TotalBalance:          &model.TotalBalance,
		ProjectId:             &model.ProjectId,
	}, nil
}

func (r *gormRepository) SearchByUserId(ctx context.Context, userId uint) ([]*domain.Budget, error) {
	var models []entity.Budget
	if err := r.db.WithContext(ctx).Joins("inner join projects projects on projects.id = budgets.project_id").
		Where("projects.user_id = ?", userId).Find(&models).Error; err != nil {
		return nil, err
	}

	var budgets []*domain.Budget
	for _, model := range models {
		modelC := model
		budgets = append(budgets, &domain.Budget{
			ID:                    &modelC.ID,
			Name:                  &modelC.Name,
			Year:                  &modelC.Year,
			Month:                 &modelC.Month,
			FixedIncome:           &modelC.FixedIncome,
			AdditionalIncome:      &modelC.AdditionalIncome,
			TotalPendingPayment:   &modelC.TotalPendingPayment,
			TotalAvailableBalance: &modelC.TotalAvailableBalance,
			PendingBills:          &modelC.PendingBills,
			TotalBalance:          &modelC.TotalBalance,
			ProjectId:             &modelC.ProjectId,
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
