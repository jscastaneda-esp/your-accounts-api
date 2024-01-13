package user_token

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

func (r *gormRepository) WithTransaction(tx persistent.Transaction) domain.UserTokenRepository {
	return db.DefaultWithTransaction[domain.UserTokenRepository](tx, NewRepository, r)
}

func (r *gormRepository) Save(ctx context.Context, userToken domain.UserToken) (uint, error) {
	model := &entity.UserToken{
		Token:     userToken.Token,
		UserId:    userToken.UserId,
		ExpiresAt: userToken.ExpiresAt,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

func (r *gormRepository) SearchByExample(ctx context.Context, example domain.UserToken) (*domain.UserToken, error) {
	where := &entity.UserToken{
		Token:  example.Token,
		UserId: example.UserId,
	}
	model := new(entity.UserToken)
	if err := r.db.WithContext(ctx).Where(where).First(model).Error; err != nil {
		return nil, err
	}

	return &domain.UserToken{
		ID:        model.ID,
		Token:     model.Token,
		UserId:    model.UserId,
		ExpiresAt: model.ExpiresAt,
	}, nil
}

func (r *gormRepository) DeleteByExpiresAtGreaterThanNow(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Where("expires_at < NOW()").Delete(&entity.UserToken{}).Error; err != nil {
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB) domain.UserTokenRepository {
	return &gormRepository{db}
}
