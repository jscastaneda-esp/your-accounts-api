package user_token

import (
	"context"
	"your-accounts-api/shared/domain/persistent"
	persistentInfra "your-accounts-api/shared/infrastructure/db/persistent"
	"your-accounts-api/user/domain"
	"your-accounts-api/user/infrastructure/entity"

	"gorm.io/gorm"
)

type gormUserTokenRepository struct {
	db *gorm.DB
}

func (r *gormUserTokenRepository) WithTransaction(tx persistent.Transaction) domain.UserTokenRepository {
	return persistentInfra.DefaultWithTransaction[domain.UserTokenRepository](tx, NewRepository, r)
}

func (r *gormUserTokenRepository) Create(ctx context.Context, userToken domain.UserToken) (uint, error) {
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

func (r *gormUserTokenRepository) FindByTokenAndUserId(ctx context.Context, token string, userId uint) (*domain.UserToken, error) {
	where := &entity.UserToken{
		Token:  token,
		UserId: userId,
	}
	model := new(entity.UserToken)
	if err := r.db.WithContext(ctx).Where(where).First(model).Error; err != nil {
		return nil, err
	}

	return &domain.UserToken{
		ID:        model.ID,
		Token:     model.Token,
		UserId:    model.UserId,
		CreatedAt: model.CreatedAt,
		ExpiresAt: model.ExpiresAt,
	}, nil
}

func (r *gormUserTokenRepository) Update(ctx context.Context, userToken domain.UserToken) error {
	model := new(entity.UserToken)
	if err := r.db.WithContext(ctx).First(model, userToken.ID).Error; err != nil {
		return err
	}

	model.Token = userToken.Token
	model.UserId = userToken.UserId
	model.ExpiresAt = userToken.ExpiresAt
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB) domain.UserTokenRepository {
	return &gormUserTokenRepository{db}
}
