package budget

import (
	"context"
	"your-accounts-api/budgets/domain"
	"your-accounts-api/budgets/infrastructure/db/entity"
	"your-accounts-api/shared/domain/persistent"
	"your-accounts-api/shared/infrastructure/db"
	shared_ent "your-accounts-api/shared/infrastructure/db/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormRepository struct {
	db *gorm.DB
}

func (r *gormRepository) WithTransaction(tx persistent.Transaction) domain.BudgetRepository {
	return db.DefaultWithTransaction[domain.BudgetRepository](tx, NewRepository, r)
}

func (r *gormRepository) Save(ctx context.Context, budget domain.Budget) (uint, error) {
	model := new(entity.Budget)
	if budget.ID != nil {
		if err := r.db.WithContext(ctx).First(model, *budget.ID).Error; err != nil {
			return 0, err
		}
	} else {
		model.UserId = *budget.UserId
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

	if budget.TotalPending != nil {
		model.TotalPending = *budget.TotalPending
	}

	if budget.TotalAvailable != nil {
		model.TotalAvailable = *budget.TotalAvailable
	}

	if budget.TotalSaving != nil {
		model.TotalSaving = *budget.TotalSaving
	}

	if budget.PendingBills != nil {
		model.PendingBills = *budget.PendingBills
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
		ID:               &model.ID,
		Name:             &model.Name,
		Year:             &model.Year,
		Month:            &model.Month,
		FixedIncome:      &model.FixedIncome,
		AdditionalIncome: &model.AdditionalIncome,
		TotalPending:     &model.TotalPending,
		TotalAvailable:   &model.TotalAvailable,
		TotalSaving:      &model.TotalSaving,
		PendingBills:     &model.PendingBills,
		UserId:           &model.UserId,
	}, nil
}

func (r *gormRepository) SearchAllByExample(ctx context.Context, example domain.Budget) ([]*domain.Budget, error) {
	where := &entity.Budget{
		UserId: *example.UserId,
	}
	var models []entity.Budget
	if err := r.db.WithContext(ctx).Where(where).Find(&models).Error; err != nil {
		return nil, err
	}

	var budgets []*domain.Budget
	for _, model := range models {
		modelC := model
		budgets = append(budgets, &domain.Budget{
			ID:               &modelC.ID,
			Name:             &modelC.Name,
			Year:             &modelC.Year,
			Month:            &modelC.Month,
			FixedIncome:      &modelC.FixedIncome,
			AdditionalIncome: &modelC.AdditionalIncome,
			TotalPending:     &modelC.TotalPending,
			TotalAvailable:   &modelC.TotalAvailable,
			TotalSaving:      &modelC.TotalSaving,
			PendingBills:     &modelC.PendingBills,
			UserId:           &modelC.UserId,
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
