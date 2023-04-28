package user_token

import (
	"api-your-accounts/user/domain"
	"api-your-accounts/user/infrastructure/entity"
	"context"

	"gorm.io/gorm"
)

type gormUserTokenRepository struct {
	db *gorm.DB
}

func (r *gormUserTokenRepository) Create(ctx context.Context, userToken *domain.UserToken) (*domain.UserToken, error) {
	model := &entity.UserToken{
		Token:       userToken.Token,
		UserId:      userToken.UserId,
		ExpiresAt:   userToken.ExpiresAt,
		RefreshedAt: userToken.RefreshedAt,
		RefreshedId: userToken.RefreshedId,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	return &domain.UserToken{
		ID:          model.ID,
		Token:       model.Token,
		UserId:      model.UserId,
		RefreshedId: model.RefreshedId,
		CreatedAt:   model.CreatedAt,
		ExpiresAt:   model.ExpiresAt,
		RefreshedAt: model.RefreshedAt,
	}, nil
}

func NewRepository(db *gorm.DB) domain.UserTokenRepository {
	return &gormUserTokenRepository{db}
}
