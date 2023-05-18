package domain

import (
	"api-your-accounts/shared/domain/persistent"
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
	persistent.TransactionRepository[UserTokenRepository]
	persistent.CreateRepository[UserToken]
	FindByTokenAndUserId(ctx context.Context, token string, userId uint) (*UserToken, error)
	persistent.UpdateRepository[UserToken]
}
