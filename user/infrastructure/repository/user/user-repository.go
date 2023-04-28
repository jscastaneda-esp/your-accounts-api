package user

import (
	"api-your-accounts/user/domain"
	"api-your-accounts/user/infrastructure/entity"
	"context"

	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

func (r *gormUserRepository) FindByUUIDAndEmail(ctx context.Context, uuid string, email string) (*domain.User, error) {
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

func (r *gormUserRepository) ExistsByUUID(ctx context.Context, uuid string) (bool, error) {
	var count int
	err := r.db.Raw("SELECT COUNT(1) FROM users WHERE uuid = ?", uuid).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *gormUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int
	err := r.db.Raw("SELECT COUNT(1) FROM users WHERE email = ?", email).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *gormUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
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

func NewRepository(db *gorm.DB) domain.UserRepository {
	return &gormUserRepository{db}
}
