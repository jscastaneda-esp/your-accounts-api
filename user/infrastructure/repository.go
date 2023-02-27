package infrastructure

import (
	"api-your-accounts/user/domain"
	"context"

	"gorm.io/gorm"
)

type GORMUserRepository struct {
	db *gorm.DB
}

func (r *GORMUserRepository) FindByUUIDAndEmail(ctx context.Context, uuid string, email string) (*domain.User, error) {
	where := &User{
		UUID:  uuid,
		Email: email,
	}
	model := new(User)
	if err := r.db.WithContext(ctx).Where(where).First(model).Error; err != nil {
		return nil, err
	}

	return &domain.User{
		Id:        model.ID,
		UUID:      model.UUID,
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

func (r *GORMUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	model := &User{
		UUID:  user.UUID,
		Email: user.Email,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	return &domain.User{
		Id:        model.ID,
		UUID:      model.UUID,
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

func (r *GORMUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	return nil, nil
}

var instance *GORMUserRepository

func NewRepository(db *gorm.DB) *GORMUserRepository {
	if instance == nil {
		instance = &GORMUserRepository{db}
	}

	return instance
}
