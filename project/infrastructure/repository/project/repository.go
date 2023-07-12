package project

import (
	"context"
	"your-accounts-api/project/domain"
	"your-accounts-api/project/infrastructure/entity"
	"your-accounts-api/shared/domain/persistent"
	persistentInfra "your-accounts-api/shared/infrastructure/db/persistent"

	"gorm.io/gorm"
)

type gormProjectRepository struct {
	db *gorm.DB
}

func (r *gormProjectRepository) WithTransaction(tx persistent.Transaction) domain.ProjectRepository {
	return persistentInfra.DefaultWithTransaction[domain.ProjectRepository](tx, NewRepository, r)
}

func (r *gormProjectRepository) Create(ctx context.Context, project domain.Project) (uint, error) {
	model := &entity.Project{
		UserId: project.UserId,
		Type:   project.Type,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

func (r *gormProjectRepository) FindById(ctx context.Context, id uint) (*domain.Project, error) {
	model, err := r.findById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Project{
		ID:        model.ID,
		UserId:    model.UserId,
		Type:      model.Type,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

func (r *gormProjectRepository) FindByUserId(ctx context.Context, userId uint) ([]*domain.Project, error) {
	where := &entity.Project{
		UserId: userId,
	}
	var models []entity.Project
	if err := r.db.WithContext(ctx).Where(where).Order("created_at desc").Limit(10).Find(&models).Error; err != nil {
		return nil, err
	}

	var projects []*domain.Project
	for _, model := range models {
		projects = append(projects, &domain.Project{
			ID:        model.ID,
			UserId:    model.UserId,
			Type:      model.Type,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		})
	}

	return projects, nil
}

func (r *gormProjectRepository) Delete(ctx context.Context, id uint) error {
	model, err := r.findById(ctx, id)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Select("ProjectLogs").Delete(model).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormProjectRepository) findById(ctx context.Context, id uint) (*entity.Project, error) {
	model := new(entity.Project)
	if err := r.db.WithContext(ctx).First(model, id).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func NewRepository(db *gorm.DB) domain.ProjectRepository {
	return &gormProjectRepository{db}
}
