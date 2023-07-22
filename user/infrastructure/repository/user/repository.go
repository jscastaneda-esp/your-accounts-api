package user

import (
	"context"
	"your-accounts-api/shared/domain/persistent"
	persistent_infra "your-accounts-api/shared/infrastructure/db/persistent"
	"your-accounts-api/user/domain"
	"your-accounts-api/user/infrastructure/entity"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func (r *gormRepository) WithTransaction(tx persistent.Transaction) domain.UserRepository {
	return persistent_infra.DefaultWithTransaction[domain.UserRepository](tx, NewRepository, r)
}

func (r *gormRepository) FindByUIDAndEmail(ctx context.Context, uid string, email string) (*domain.User, error) {
	where := &entity.User{
		UID:   uid,
		Email: email,
	}
	model := new(entity.User)
	if err := r.db.WithContext(ctx).Where(where).First(model).Error; err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        model.ID,
		UID:       model.UID,
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

func (r *gormRepository) ExistsByUID(ctx context.Context, uid string) (bool, error) {
	var count int
	err := r.db.Raw("SELECT COUNT(1) FROM users WHERE uid = ?", uid).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *gormRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int
	err := r.db.Raw("SELECT COUNT(1) FROM users WHERE email = ?", email).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *gormRepository) Create(ctx context.Context, user domain.User) (uint, error) {
	model := &entity.User{
		UID:   user.UID,
		Email: user.Email,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

func NewRepository(db *gorm.DB) domain.UserRepository {
	return &gormRepository{db}
}
