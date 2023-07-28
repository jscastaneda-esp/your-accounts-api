package project

import (
	"context"
	"your-accounts-api/project/domain"
	"your-accounts-api/project/infrastructure/entity"
	"your-accounts-api/shared/domain/persistent"
	shared_ent "your-accounts-api/shared/infrastructure/db/entity"
	persistent_infra "your-accounts-api/shared/infrastructure/db/persistent"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func (r *gormRepository) WithTransaction(tx persistent.Transaction) domain.ProjectRepository {
	return persistent_infra.DefaultWithTransaction[domain.ProjectRepository](tx, NewRepository, r)
}

func (r *gormRepository) Save(ctx context.Context, project domain.Project) (uint, error) {
	model := &entity.Project{
		UserId: project.UserId,
		Type:   project.Type,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Select("ProjectLogs").Delete(&entity.Project{
		BaseModel: shared_ent.BaseModel{
			ID: id,
		},
	}).Error; err != nil {
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB) domain.ProjectRepository {
	return &gormRepository{db}
}
