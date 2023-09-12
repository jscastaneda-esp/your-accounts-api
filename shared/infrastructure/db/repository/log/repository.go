package log

import (
	"context"
	"your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/persistent"
	"your-accounts-api/shared/infrastructure/db"
	"your-accounts-api/shared/infrastructure/db/entity"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func (r *gormRepository) WithTransaction(tx persistent.Transaction) domain.LogRepository {
	return db.DefaultWithTransaction[domain.LogRepository](tx, NewRepository, r)
}

func (r *gormRepository) Save(ctx context.Context, log domain.Log) (uint, error) {
	if log.Detail == nil {
		log.Detail = map[string]any{}
	}

	model := &entity.Log{
		Description: log.Description,
		Detail:      log.Detail,
		Code:        log.Code,
		ResourceId:  log.ResourceId,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

func (r *gormRepository) SearchAllByExample(ctx context.Context, example domain.Log) ([]*domain.Log, error) {
	where := &entity.Log{
		Code:       example.Code,
		ResourceId: example.ResourceId,
	}
	var models []entity.Log
	if err := r.db.WithContext(ctx).Where(where).Order("created_at desc").Limit(20).Find(&models).Error; err != nil {
		return nil, err
	}

	var logs []*domain.Log
	for _, model := range models {
		logs = append(logs, &domain.Log{
			ID:          model.ID,
			Description: model.Description,
			Detail:      model.Detail,
			Code:        model.Code,
			ResourceId:  model.ResourceId,
			CreatedAt:   model.CreatedAt,
		})
	}

	return logs, nil
}

func NewRepository(db *gorm.DB) domain.LogRepository {
	return &gormRepository{db}
}
