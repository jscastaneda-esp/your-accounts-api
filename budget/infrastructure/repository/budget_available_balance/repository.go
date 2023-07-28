package budget_available_balance

import (
	"context"
	"your-accounts-api/budget/domain"
	"your-accounts-api/budget/infrastructure/entity"
	"your-accounts-api/shared/domain/persistent"
	shared_ent "your-accounts-api/shared/infrastructure/db/entity"
	persistent_infra "your-accounts-api/shared/infrastructure/db/persistent"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func (r *gormRepository) WithTransaction(tx persistent.Transaction) domain.BudgetAvailableBalanceRepository {
	return persistent_infra.DefaultWithTransaction[domain.BudgetAvailableBalanceRepository](tx, NewRepository, r)
}

func (r *gormRepository) Save(ctx context.Context, available domain.BudgetAvailableBalance) (uint, error) {
	model := new(entity.BudgetAvailableBalance)
	if available.ID != nil {
		if err := r.db.WithContext(ctx).First(model, *available.ID).Error; err != nil {
			return 0, err
		}
	} else {
		model.BudgetId = *available.BudgetId
	}

	if available.Name != nil {
		model.Name = *available.Name
	}

	if available.Amount != nil {
		model.Amount = *available.Amount
	}

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

func (r *gormRepository) SaveAll(ctx context.Context, availables []domain.BudgetAvailableBalance) error {
	models := []*entity.BudgetAvailableBalance{}
	for _, available := range availables {
		model := new(entity.BudgetAvailableBalance)
		if available.ID != nil {
			model.ID = *available.ID
		}

		if available.Name != nil {
			model.Name = *available.Name
		}

		if available.Amount != nil {
			model.Amount = *available.Amount
		}

		if available.BudgetId != nil {
			model.BudgetId = *available.BudgetId
		}

		models = append(models, model)
	}

	if err := r.db.WithContext(ctx).Save(models).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormRepository) SearchAllByExample(ctx context.Context, example domain.BudgetAvailableBalance) ([]*domain.BudgetAvailableBalance, error) {
	where := &entity.BudgetAvailableBalance{
		BudgetId: *example.BudgetId,
	}
	var models []entity.BudgetAvailableBalance
	if err := r.db.WithContext(ctx).Where(where).Find(&models).Error; err != nil {
		return nil, err
	}

	var projects []*domain.BudgetAvailableBalance
	for _, model := range models {
		modelC := model
		projects = append(projects, &domain.BudgetAvailableBalance{
			ID:       &modelC.ID,
			Name:     &modelC.Name,
			Amount:   &modelC.Amount,
			BudgetId: &modelC.BudgetId,
		})
	}

	return projects, nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.BudgetAvailableBalance{
		BaseModel: shared_ent.BaseModel{
			ID: id,
		},
	}).Error; err != nil {
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB) domain.BudgetAvailableBalanceRepository {
	return &gormRepository{db}
}
