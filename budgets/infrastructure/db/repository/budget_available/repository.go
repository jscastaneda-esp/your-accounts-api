package budget_available

import (
	"context"
	"your-accounts-api/budgets/domain"
	"your-accounts-api/budgets/infrastructure/db/entity"
	"your-accounts-api/shared/domain/persistent"
	"your-accounts-api/shared/infrastructure/db"
	shared_ent "your-accounts-api/shared/infrastructure/db/entity"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func (r *gormRepository) WithTransaction(tx persistent.Transaction) domain.BudgetAvailableRepository {
	return db.DefaultWithTransaction[domain.BudgetAvailableRepository](tx, NewRepository, r)
}

func (r *gormRepository) Save(ctx context.Context, available domain.BudgetAvailable) (uint, error) {
	model := new(entity.BudgetAvailable)
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

func (r *gormRepository) SaveAll(ctx context.Context, availables []domain.BudgetAvailable) error {
	if len(availables) == 0 {
		return nil
	}

	models := []*entity.BudgetAvailable{}
	for _, available := range availables {
		model := new(entity.BudgetAvailable)
		if available.ID != nil {
			if err := r.db.WithContext(ctx).First(model, *available.ID).Error; err != nil {
				return err
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

		models = append(models, model)
	}

	if err := r.db.WithContext(ctx).Save(models).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.BudgetAvailable{
		BaseModel: shared_ent.BaseModel{
			ID: id,
		},
	}).Error; err != nil {
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB) domain.BudgetAvailableRepository {
	return &gormRepository{db}
}
