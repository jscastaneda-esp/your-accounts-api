package log

import (
	"context"
	"database/sql"
	"your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/persistent"
	"your-accounts-api/shared/infrastructure/db"
	"your-accounts-api/shared/infrastructure/db/entity"

	"gorm.io/gorm"
)

const limit = 20

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

func (r *gormRepository) SearchAllByExample(ctx context.Context, example domain.Log) ([]domain.Log, error) {
	where := &entity.Log{
		Code:       example.Code,
		ResourceId: example.ResourceId,
	}
	var models []entity.Log
	if err := r.db.WithContext(ctx).Where(where).Order("created_at desc").Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}

	var logs []domain.Log
	for _, model := range models {
		logs = append(logs, domain.Log{
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

func (r *gormRepository) DeleteByResourceIdNotExists(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Where("resource_id NOT IN (SELECT id FROM budgets) AND resource_id NOT IN (SELECT id FROM budget_bills)").Delete(&entity.Log{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormRepository) SearchResourceIdsWithLimitExceeded(ctx context.Context) ([]uint, error) {
	var ids []uint
	if err := r.db.WithContext(ctx).Raw("SELECT resource_id FROM logs GROUP BY resource_id HAVING COUNT(resource_id) > ?", limit).Scan(&ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *gormRepository) DeleteByResourceIdAndIdLessThanLimit(ctx context.Context, resourceId uint) error {
	if err := r.db.WithContext(ctx).Where("resource_id = @resource_id AND id < (SELECT id FROM (SELECT id FROM logs WHERE resource_id = @resource_id ORDER BY id DESC LIMIT @limit) T ORDER BY id ASC LIMIT 1)", sql.Named("resource_id", resourceId), sql.Named("limit", limit)).Delete(&entity.Log{}).Error; err != nil {
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB) domain.LogRepository {
	return &gormRepository{db}
}
