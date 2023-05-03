package domain

import (
	"api-your-accounts/shared/domain/transaction"
	"context"
	"time"
)

type UserToken struct {
	ID          uint
	Token       string
	UserId      uint
	RefreshedBy *uint
	CreatedAt   time.Time
	ExpiresAt   time.Time
	RefreshedAt *time.Time
}

//go:generate mockery --name UserTokenRepository --filename user-token-repository.go
type UserTokenRepository interface {
	WithTransaction(tx transaction.Transaction) UserTokenRepository
	Create(ctx context.Context, userToken *UserToken) (*UserToken, error)
	FindByTokenAndUserId(ctx context.Context, token string, userId uint) (*UserToken, error)
	Update(ctx context.Context, userToken *UserToken) (*UserToken, error)
}
