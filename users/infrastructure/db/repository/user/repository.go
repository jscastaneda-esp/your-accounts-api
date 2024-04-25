package user

import (
	"context"
	"your-accounts-api/shared/domain/persistent"
	"your-accounts-api/shared/infrastructure/db"
	"your-accounts-api/users/domain"
	"your-accounts-api/users/infrastructure/db/entity"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func (r *gormRepository) WithTransaction(tx persistent.Transaction) domain.UserRepository {
	return db.DefaultWithTransaction[domain.UserRepository](tx, NewRepository, r)
}

func (r *gormRepository) Save(ctx context.Context, user domain.User) (uint, error) {
	model := &entity.User{
		Email: user.Email,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

func (r *gormRepository) SearchByExample(ctx context.Context, example domain.User) (domain.User, error) {
	where := entity.User{
		Email: example.Email,
	}
	model := new(entity.User)
	if err := r.db.WithContext(ctx).Where(where).First(model).Error; err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:    model.ID,
		Email: model.Email,
	}, nil
}

func (r *gormRepository) ExistsByExample(ctx context.Context, example domain.User) (bool, error) {
	var count int64
	where := entity.User{}
	if example.Email != "" {
		where.Email = example.Email
	}

	err := r.db.WithContext(ctx).Model(where).Where(where).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func NewRepository(db *gorm.DB) domain.UserRepository {
	return &gormRepository{db}
}
