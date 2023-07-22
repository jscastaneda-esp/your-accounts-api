package budget_available_balance

import (
	"context"
	"your-accounts-api/budget/domain"
	"your-accounts-api/budget/infrastructure/entity"
	"your-accounts-api/shared/domain/persistent"
	persistent_infra "your-accounts-api/shared/infrastructure/db/persistent"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func (r *gormRepository) WithTransaction(tx persistent.Transaction) domain.BudgetAvailableBalanceRepository {
	return persistent_infra.DefaultWithTransaction[domain.BudgetAvailableBalanceRepository](tx, NewRepository, r)
}

func (r *gormRepository) Create(ctx context.Context, available domain.BudgetAvailableBalance) (uint, error) {
	model := &entity.BudgetAvailableBalance{
		Name:     available.Name,
		Amount:   available.Amount,
		BudgetId: available.BudgetId,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

func NewRepository(db *gorm.DB) domain.BudgetAvailableBalanceRepository {
	return &gormRepository{db}
}
