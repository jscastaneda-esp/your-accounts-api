// TODO: Pendientes tests

package repository

import (
	"api-your-accounts/user/domain"
	"api-your-accounts/user/infrastructure/entity"
	"context"

	"gorm.io/gorm"
)

type GORMUserRepository struct {
	db *gorm.DB
}

func (r *GORMUserRepository) FindByUUIDAndEmail(ctx context.Context, uuid string, email string) (*domain.User, error) {
	where := &entity.User{
		UUID:  uuid,
		Email: email,
	}
	model := new(entity.User)
	if err := r.db.WithContext(ctx).Where(where).First(model).Error; err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        model.ID,
		UUID:      model.UUID,
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

func (r *GORMUserRepository) ExistsByUUID(ctx context.Context, uuid string) (bool, error) {
	var count int
	r.db.Raw("SELECT COUNT(1) FROM users WHERE uuid = ?", uuid).Scan(&count)
	if r.db.Error != nil {
		return false, r.db.Error
	}

	return count > 0, nil
}

func (r *GORMUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int
	r.db.Raw("SELECT COUNT(1) FROM users WHERE email = ?", email).Scan(&count)
	if r.db.Error != nil {
		return false, r.db.Error
	}

	return count > 0, nil
}

func (r *GORMUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	model := &entity.User{
		UUID:  user.UUID,
		Email: user.Email,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        model.ID,
		UUID:      model.UUID,
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

func (r *GORMUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	return nil, nil
}

var instance domain.UserRepository

func NewGORMRepository(db *gorm.DB) domain.UserRepository {
	if instance == nil {
		instance = &GORMUserRepository{db}
	}

	return instance
}
