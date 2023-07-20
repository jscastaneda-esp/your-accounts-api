package user

import (
	"context"
	"your-accounts-api/shared/domain/persistent"
	"your-accounts-api/shared/infrastructure/db"
	persistentInfra "your-accounts-api/shared/infrastructure/db/persistent"
	"your-accounts-api/user/domain"
	"your-accounts-api/user/infrastructure/entity"

	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

func (r *gormUserRepository) WithTransaction(tx persistent.Transaction) domain.UserRepository {
	return persistentInfra.DefaultWithTransaction[domain.UserRepository](tx, newRepository, r)
}

func (r *gormUserRepository) FindByUIDAndEmail(ctx context.Context, uid string, email string) (*domain.User, error) {
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

func (r *gormUserRepository) ExistsByUID(ctx context.Context, uid string) (bool, error) {
	var count int
	err := r.db.Raw("SELECT COUNT(1) FROM users WHERE uid = ?", uid).Scan(&count).Error
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

func (r *gormUserRepository) Create(ctx context.Context, user domain.User) (uint, error) {
	model := &entity.User{
		UID:   user.UID,
		Email: user.Email,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

func newRepository(db *gorm.DB) domain.UserRepository {
	return &gormUserRepository{db}
}

var instance domain.UserRepository

func DefaultRepository() domain.UserRepository {
	if instance == nil {
		instance = newRepository(db.DB)
	}

	return instance
}
