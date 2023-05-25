package project_log

import (
	"api-your-accounts/project/domain"
	"api-your-accounts/project/infrastructure/entity"
	"api-your-accounts/shared/domain/persistent"
	persistentInfra "api-your-accounts/shared/infrastructure/db/persistent"
	"context"

	"gorm.io/gorm"
)

type gormProjectLogRepository struct {
	db *gorm.DB
}

func (r *gormProjectLogRepository) WithTransaction(tx persistent.Transaction) domain.ProjectLogRepository {
	return persistentInfra.DefaultWithTransaction[domain.ProjectLogRepository](tx, NewRepository, r)
}

func (r *gormProjectLogRepository) Create(ctx context.Context, projectLog domain.ProjectLog) (*domain.ProjectLog, error) {
	model := &entity.ProjectLog{
		Description: projectLog.Description,
		Detail:      projectLog.Detail,
		ProjectId:   projectLog.ProjectId,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	return &domain.ProjectLog{
		ID:          model.ID,
		Description: model.Description,
		Detail:      model.Detail,
		ProjectId:   model.ProjectId,
		CreatedAt:   model.CreatedAt,
	}, nil
}

func (r *gormProjectLogRepository) FindByProjectId(ctx context.Context, projectId uint) ([]*domain.ProjectLog, error) {
	where := &entity.ProjectLog{
		ProjectId: projectId,
	}
	var models []entity.ProjectLog
	if err := r.db.WithContext(ctx).Where(where).Order("created_at desc").Limit(20).Find(&models).Error; err != nil {
		return nil, err
	}

	var projects []*domain.ProjectLog
	for _, model := range models {
		projects = append(projects, &domain.ProjectLog{
			ID:          model.ID,
			Description: model.Description,
			Detail:      model.Detail,
			ProjectId:   model.ProjectId,
			CreatedAt:   model.CreatedAt,
		})
	}

	return projects, nil
}

func NewRepository(db *gorm.DB) domain.ProjectLogRepository {
	return &gormProjectLogRepository{db}
}
