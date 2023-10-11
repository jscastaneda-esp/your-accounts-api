package budget_bill

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

func (r *gormRepository) WithTransaction(tx persistent.Transaction) domain.BudgetBillRepository {
	return db.DefaultWithTransaction[domain.BudgetBillRepository](tx, NewRepository, r)
}

func (r *gormRepository) Save(ctx context.Context, bill domain.BudgetBill) (uint, error) {
	model := new(entity.BudgetBill)
	if bill.ID != nil {
		if err := r.db.WithContext(ctx).First(model, *bill.ID).Error; err != nil {
			return 0, err
		}
	} else {
		model.BudgetId = *bill.BudgetId
	}

	if bill.Description != nil {
		model.Description = *bill.Description
	}

	if bill.Amount != nil {
		model.Amount = *bill.Amount
	}

	if bill.Payment != nil {
		model.Payment = *bill.Payment
	}

	if bill.DueDate != nil {
		model.DueDate = *bill.DueDate
	}

	if bill.Complete != nil {
		model.Complete = *bill.Complete
	}

	if bill.Category != nil {
		model.Category = *bill.Category
	}

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

func (r *gormRepository) SaveAll(ctx context.Context, bills []domain.BudgetBill) error {
	if len(bills) == 0 {
		return nil
	}

	models := []*entity.BudgetBill{}
	for _, bill := range bills {
		model := new(entity.BudgetBill)
		if bill.ID != nil {
			if err := r.db.WithContext(ctx).First(model, *bill.ID).Error; err != nil {
				return err
			}
		} else {
			model.BudgetId = *bill.BudgetId
		}

		if bill.Description != nil {
			model.Description = *bill.Description
		}

		if bill.Amount != nil {
			model.Amount = *bill.Amount
		}

		if bill.Payment != nil {
			model.Payment = *bill.Payment
		}

		if bill.DueDate != nil {
			model.DueDate = *bill.DueDate
		}

		if bill.Complete != nil {
			model.Complete = *bill.Complete
		}

		if bill.Category != nil {
			model.Category = *bill.Category
		}

		models = append(models, model)
	}

	if err := r.db.WithContext(ctx).Save(models).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.BudgetBill{
		BaseModel: shared_ent.BaseModel{
			ID: id,
		},
	}).Error; err != nil {
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB) domain.BudgetBillRepository {
	return &gormRepository{db}
}
